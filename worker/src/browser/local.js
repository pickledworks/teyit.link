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

    await index.run(browser, guid(), 'https://www.hurriyetemlak.com/projeler/juma-insaat/juma-plus-juma-insaat?utm_source=hurriyet_emlak&utm_medium=he_newsearch&utm_content=flat-plan-id-88&utm_campaign=juma-plus-juma-insaat')
    .then((result) => console.log(result))
    .catch((err) => console.error(err));
    await browser.close();
})();
