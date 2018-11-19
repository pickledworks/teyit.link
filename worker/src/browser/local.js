const index = require('../index');
const config = require('./config');
const puppeteer = require('puppeteer');

function guid() {
    function s4() {
      return Math.floor((1 + Math.random()) * 0x10000)
        .toString(16)
        .substring(1);
    }
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4();
  }

(async () => {
    const browser = await puppeteer.launch({
        headless: false,
        slowMo: process.env.SLOWMO_MS,
        dumpio: !!config.DEBUG,
        // use chrome installed by puppeteer
    });

    await index.run(browser, guid(), 'https://github.com/noddigital/teyit.link/pulls')
    .then((result) => console.log(result))
    .catch((err) => console.error(err));
    await browser.close();
})();
