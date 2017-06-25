// teyit.link

'use strict';

// node modules

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _moment = require('moment');

var _moment2 = _interopRequireDefault(_moment);

var _archive = require('../models/archive');

var _archive2 = _interopRequireDefault(_archive);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

// configuration
_moment2.default.locale('tr');

// local modules


var archive = {
  // /new
  new: {
    post: function post(req, res, next) {
      req.request_url = req.body.request_url;
      next();
    },
    get: function get(req, res, next) {
      req.request_url = req.query.request_url;
      next();
    }
  },

  // /search (/?q=)
  search: {
    get: function get(req, res) {
      if (!req.query.q) {
        res.status(400).redirect('/');
      }

      var query = req.query.q.replace(/\W/g, '');

      if (query.length < 3) {
        res.status(400).redirect('/');
      }

      _archive2.default.findAll({
        where: {
          $or: {
            request_url: {
              $like: '%' + query + '%'
            },
            meta_title: {
              $like: '%' + query + '%'
            }
          }
        }
      }).then(function (results) {
        if (results.length < 1) {
          res.status(204).redirect('/?empty');
        } else {
          res.status(200).render('pages/search', {
            results: results.map(function (result) {
              return result.dataValues;
            })
          });
        }
      });
    }
  },

  // /slug (/^\/(\w{7})?$/)
  slug: {
    get: function get(req, res) {
      var slug = req.params[0];

      _archive2.default.findOne({
        where: {
          slug: slug
        }
      }).then(function (result) {
        if (result && 'dataValues' in result) {
          result.dataValues.created_at = (0, _moment2.default)(result.dataValues.created_at).format('LL LTS'); // 01 Ocak 2017 12:34:56

          res.status(200).render('pages/detail', result.dataValues);
        } else {
          res.status(404).redirect('/?notfound');
        }
      });
    }
  }
};

// exports
exports.default = archive;