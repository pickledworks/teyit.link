'use strict';
const Inliner = require('../vendor/inliner/lib/index');
const request = require('request')
const phantomjs = require('phantomjs-prebuilt')
const webdriverio = require('webdriverio')
const wdOpts = { desiredCapabilities: { browserName: 'phantomjs' } }

const scrollHelper = require('./scroll-helper');

const customRequestAdaptor = (url, settings, callback) => {
  console.log("custom: " + url, settings);
  if (!settings.initialRequest) {
    if (url.includes('&amp;oe=')) { // quick and dirty fix for facebook (temporary)
      url = url.replace('&amp;oe=', '&oe=');
    }
    return request(url, settings, callback);
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
            browser.saveScreenshot('/tmp/' + settings.archiveID + '.png').then((screenshot) => {
              callback(error, res, body);
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
      resolve(data);
    }
  });
});

