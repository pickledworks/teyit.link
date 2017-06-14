'use strict';
require('dotenv').config();
const fs = require('fs');
const AWS = require('aws-sdk');
const stripJs = require('strip-js');
const url = require('url');
const getter = require('./utils/getter');
const uploader = require('./utils/uploader');
const validator = require('validator');

const validateURL = (d) => validator.isURL(d);

AWS.config.update({
  accessKeyId: process.env.AWS_ACCESS_KEY,
  secretAccessKey: process.env.AWS_SECRET_KEY,
  region: process.env.AWS_REGION
});

const bucket = process.env.AWS_BUCKET;
const client = new AWS.S3();

exports.handler = (event, context, cb) => {
  const archive_id = event.archive_id;
  const request_url = event.request_url;

  if (!validateURL(request_url)) {
    cb(true, 'Url is not valid.');
    return;
  }

  return getter(request_url, archive_id) // Step one: Create rendered HTML and screenshot
    .then((data) => { // Step two: Upload generated files
      const htmlUpload = uploader(data.html, client, bucket, archive_id, 'index.html', 'text/html');
      const ss = fs.readFileSync('/tmp/' + archive_id + '.png');
      const ssUpload = uploader(ss, client, bucket, archive_id, 'screenshot.png', 'image/png');
      return Promise.all([htmlUpload, ssUpload]).then(() => data);
    })
    .then((data) => {
      cb(null, {
        description: data.description,
        title: data.title,
      });
    })
    .catch((err) => {
      console.log('Crawl or save error', err);
      cb(err, 'Crawl or save error!');
    });
};
