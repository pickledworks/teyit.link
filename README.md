# teyitlink-archive

This project provides the [AWS Lambda](https://aws.amazon.com/lambda/) function used in the [teyitlink-web](https://github.com/teyit/teyitlink-web) repository for [teyit.link](https://teyit.link/?ref=teyitlink-archiver-github).

The function first loads a page using a headless browser, runs a few helper scripts on the page, takes a screenshot and extracts the HTML after a few seconds, loads and inlines most of the static assets (images, stylesheets, videos) into the HTML and uploads the resulted HTML and the screenshot to [AWS S3](https://aws.amazon.com/s3/).

We use a prebuilt [PhantomJS](http://phantomjs.org/) build for the inital load and screenshots, a modified version of the awesome [inliner](https://github.com/remy/inliner) library then loads the assets and does the actual inlining. We use the offical [aws-sdk](https://github.com/aws/aws-sdk-js) for uploading to S3.

You can build and deploy this project like any other Lambda function. If you are not on a Linux platform, run you need to set `PHANTOMJS_PLATFORM` environment variable to `linux` before running `npm install`.

## Copyright

Teyit, [teyit.org](https://teyit.org/?ref=teyitlink-archiver-github)
Platform which verifies rumors on the Internet
[info@teyit.org](mailto:info@teyit.org)