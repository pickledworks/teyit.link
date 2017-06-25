// teyit.link

'use strict';

// node modules

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _express = require('express');

var _archive = require('./controllers/archive');

var _archive2 = _interopRequireDefault(_archive);

var _homepage = require('./controllers/homepage');

var _homepage2 = _interopRequireDefault(_homepage);

var _lambda = require('./utils/lambda');

var _lambda2 = _interopRequireDefault(_lambda);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

// variables
var router = (0, _express.Router)();

// routes


// local modules
router.route('/').get(_homepage2.default.get);
router.route('/').all(_homepage2.default.all);
router.route(['/add', '/bookmark', '/new']).all(_archive2.default.new.post, _lambda2.default);
router.route('/search').get(_archive2.default.search.get);
router.route('/search').all(_homepage2.default.all);
router.route(/^\/(\w{7})?$/).get(_archive2.default.slug.get);
router.route(/^\/(\w{7})?$/).all(_homepage2.default.all);

// exports
exports.default = router;