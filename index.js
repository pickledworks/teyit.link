'use strict';
require('dotenv').config();
var fs = require('fs');
var AWS = require('aws-sdk');
var stripJs = require('strip-js');
var url = require('url');
var getter = require('./utils/getter');
var uploader = require('./utils/uploader');

const urlregex = /^(https?|ftp):\/\/([a-zA-Z0-9.-]+(:[a-zA-Z0-9.&%$-]+)*@)*((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])){3}|([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.(com|edu|gov|int|mil|net|org|biz|arpa|info|name|pro|aero|coop|museum|[a-zA-Z]{2}))(:[0-9]+)*(\/($|[a-zA-Z0-9.,?'\\+&%$#=~_-]+))*$/;

const validateURL = (d) => urlregex.test(d);

AWS.config.update({
    accessKeyId: process.env.AWS_ACCESS_KEY,
    secretAccessKey: process.env.AWS_SECRET_KEY,
    region: process.env.AWS_REGION
});

var bucket = process.env.AWS_BUCKET;
var client = new AWS.S3();

exports.handler = (event, context, cb) => {
    var archive_id = event.archive_id;
    var request_url = event.request_url;

    if(!validateURL(request_url)){
      cb(true,'Url is not valid.');
      return;
    }

    return getter(request_url, archive_id)
      .then((data) => { // Step one: Write the string into a temporary file to upload
          const htmlUpload = uploader(stripJs(data.html), client, bucket, archive_id, 'index.html', 'text/html');
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