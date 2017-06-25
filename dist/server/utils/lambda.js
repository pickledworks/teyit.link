// teyit.link

'use strict';

// node modules

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _awsSdk = require('aws-sdk');

var _awsSdk2 = _interopRequireDefault(_awsSdk);

var _randomstring = require('randomstring');

var _uuid = require('uuid');

var _uuid2 = _interopRequireDefault(_uuid);

var _config = require('../../config');

var _archive = require('../models/archive');

var _archive2 = _interopRequireDefault(_archive);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

// configuration


// local modules
_awsSdk2.default.config.update({
  accessKeyId: _config.AWS_ACCESS_KEY_ID,
  secretAccessKey: _config.AWS_SECRET_ACCESS_KEY,
  region: _config.AWS_REGION
});

var lambdaClient = new _awsSdk2.default.Lambda({
  region: 'eu-central-1'
});

var createArchive = function createArchive(archive_id, slug) {
  var payload = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : {};
  var request_url = arguments[3];
  var res = arguments[4];

  var params = {
    archive_id: archive_id,
    slug: slug,
    meta_title: payload.title,
    meta_description: payload.description,
    request_url: request_url
  };
  _archive2.default.create(params).then(function () {
    if (req.accepts('json')) {
      res.status(200).json({
        data: params,
        status: true
      });
    } else {
      res.status(200).redirect('/' + slug);
    }
  }).catch(function () {
    if (req.accepts('json')) {
      res.status(503).json({
        status: false
      });
    } else {
      res.status(503).redirect('/?fail-create');
    }
  });
};

var lambda = function lambda(req, res) {

  var archive_id = (0, _uuid2.default)();
  var request_url = req.request_url;
  var slug = (0, _randomstring.generate)(7);

  var params = {
    FunctionName: 'teyitlink-archive',
    Payload: JSON.stringify({
      request_url: request_url,
      archive_id: archive_id
    })
  };
  lambdaClient.invoke(params, function (err, data) {
    var payload = JSON.parse(data.Payload);

    if (err) {
      if (req.accepts('json')) {
        res.status(503).json({
          status: false
        });
      } else {
        res.status(503).redirect('/?fail-create');
      }
    } else {
      createArchive(archive_id, slug, payload, request_url, res);
    }
  });
};

// exports
exports.default = lambda;