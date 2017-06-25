// teyit.link

'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
var homepage = {
  // permanently redirect all unauthorized requests and non-existing routes to hompage
  all: function all(req, res) {
    res.redirect(301, '/');
  },

  // /
  get: function get(req, res) {
    res.status(200).render('pages/homepage');
  }
};

// exports
exports.default = homepage;