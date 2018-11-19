'use strict';
const Inliner = require('../vendor/inliner/lib/index');
const request = require('request');
const stripJs = require('strip-js');

const customRequestAdaptor = (source) => (url, settings, callback) => {
  const res = {};

  if (settings.initialRequest) {
    const html = source.get('_html');
    if (!html) {
      callback('no html found', res, null);
    } else {
      callback(null, res, stripJs(html));
    }

    return;
  }

  if (source.has(url)) {
    const body = source.get(url);
    callback(null, res, body);
    return;
  }

  request(url, settings, (error, res, body) => {
    callback(error, res, body);
  });
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
