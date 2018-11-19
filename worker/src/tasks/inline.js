'use strict';
const Inliner = require('../vendor/inliner/lib/index');
const request = require('request');
const stripJs = require('strip-js');

const customRequestAdaptor = (source) => (url, settings, callback) => {
  const isFacebook = url.startsWith("https://www.facebook.com/");

  if (!settings.initialRequest || isFacebook) {
    if (url.includes('&amp;oe=')) { // quick and dirty fix for facebook (temporary)
      url = url.replace('&amp;oe=', '&oe=');
    }

    if (isFacebook) {

      if (url.indexOf("_fb_noscript=1") === -1) {
        const separator = (url.indexOf("?") === -1) ? "?" : "&";
        url = url + separator + '_fb_noscript=1';
      }

      settings.headers = {
        "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8"
      };

    }

    return request(url, settings, (error, res, body) => {
      if (isFacebook) {
        const renderedBody = stripJs(body.toString('utf8'));
        callback(error, res, renderedBody);
      } else {
        callback(error, res, body);
      }
    });

  }

  const res = {};

  callback(null, res, stripJs(source));
};

module.exports = (archiveID, requestURL, source) => new Promise((resolve, reject) => {
  new Inliner(requestURL, {
    adaptor: customRequestAdaptor(source),
    skipAbsoluteUrls: true,
    nosvg: true,
    encoding: 'utf-8',
    archiveID: archiveID,
  }, (error, data) => {
    if (error) {
      reject(error);
    } else {
      resolve(data);
    }
  });
});
