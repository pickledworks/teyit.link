require('dotenv').config();
const setup = require('./browser/setup');
const inline = require('./tasks/inline');
const upload = require('./tasks/upload');
const scrollHelper = require('./helpers/scroll-helper');
const fs = require('fs');
const AWS = require('aws-sdk');
const validator = require('validator');

AWS.config.update({
  accessKeyId: process.env.AWS_ACCESS_KEY,
  secretAccessKey: process.env.AWS_SECRET_KEY,
  region: process.env.AWS_REGION,
});

const bucket = process.env.AWS_BUCKET;
const client = new AWS.S3();

exports.handler = async (event, context, callback) => {
  // For keeping the browser launch
  context.callbackWaitsForEmptyEventLoop = false;
  const { archive_id, request_url } = event;

  if (!validator.isURL(request_url)) {
    callback('URL is not valid.');
    return;
  }

  const browser = await setup.getBrowser();

  try {
    const result = await exports.run(browser, archive_id, request_url);
    callback(null, result);
  } catch (e) {
    callback(e);
  }
};

const extract = async (archiveID, requestURL, source) => {
  const data = await inline(archiveID, requestURL, source);
  await upload(data.html, client, bucket, archiveID, 'index.html', 'text/html');
  return {
    title: data.title,
    description: data.description,
    image: data.image,
  };
};

const screenshot = async (archiveID, page) => {
  const ssFilePath = '/tmp/' + archiveID + '.png';
  await page.screenshot({path: ssFilePath});
  const ss = fs.readFileSync(ssFilePath);
  await upload(ss, client, bucket, archiveID, 'screenshot.png', 'image/png');
};

exports.run = async (browser, archiveID, requestURL) => {
  console.log("Starting to archive", archiveID, requestURL);

  const page = await browser.newPage();

  await page.setViewport({
    width: 1280,
    height: 800,
  });

  await page.goto(requestURL,
   {waitUntil: ['domcontentloaded', 'networkidle0']}
  );

  await page.evaluate(scrollHelper);

  const source = await page.evaluate('new XMLSerializer().serializeToString(document.doctype) + document.documentElement.outerHTML');

  const results = await Promise.all([
    extract(archiveID, requestURL, source),
    screenshot(archiveID, page),
  ]);

  await page.close();

  const data = results[0];
  return data;
};
