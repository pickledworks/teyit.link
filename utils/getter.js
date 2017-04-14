'use strict';
const fs = require('fs');
const Inliner = require('../vendor/inliner/lib/index');
const request = require('request');
const stripJs = require('strip-js');
const phantomjs = require('phantomjs-prebuilt');
const webshot = require('webshot');
const webdriverio = require('webdriverio');
const wdOpts = { desiredCapabilities: { browserName: 'phantomjs' } }

const scrollHelper = require('./scroll-helper');

const customRequestAdaptor = (url, settings, callback) => {
  console.log("custom: " + url, settings);
  const isFacebook = url.startsWith("https://www.facebook.com/");
  const ssFilePath = '/tmp/' + settings.archiveID + '.png';

  if (!settings.initialRequest || isFacebook) {
    if (url.includes('&amp;oe=')) { // quick and dirty fix for facebook (temporary)
      url = url.replace('&amp;oe=', '&oe=');
    }

    if (isFacebook) {
      url = url + '?_fb_noscript=1';
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
  let error = null;

  phantomjs.run('--webdriver=4444').then((program) => {
    const browser = webdriverio.remote(wdOpts).init();
    browser
      .url(url)
      .setViewportSize({
        width: 1280,
        height: 800
      })
      .execute(scrollHelper)
      .then(() => {
        setTimeout(() => {
           browser.getHTML('html', true).then((body) => {
            browser.saveScreenshot(ssFilePath).then((screenshot) => {
              callback(error, res, stripJs(body));
              program.kill();
            });
          });
        }, 5000);
      });
  });
};

module.exports = (url, archiveID) => new Promise((resolve, reject) => {
  console.log("url: " + url);
  console.log("archive id: " + archiveID);
  new Inliner(url, {
    adaptor: customRequestAdaptor,
    skipAbsoluteUrls: true,
    encoding: 'utf-8',
    archiveID: archiveID,
  }, (error, data) => {
    if (error) {
      reject(error);
    } else {

      const ssFilePath = '/tmp/' + archiveID + '.png';

      if (fs.existsSync(ssFilePath)) {
        resolve(data);
      } else {
        webshot(data.html, ssFilePath, {siteType:'html'}, function(err) {
          if (err) {
            reject(err);
          } else {
            resolve(data);
          }
        });
      }
    }
  });
});

