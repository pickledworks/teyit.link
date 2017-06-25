// teyit.link

'use strict';

// node modules

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _sequelize = require('sequelize');

var _config = require('../../config');

// variables
var options = {
  host: _config.MYSQL_HOST,
  port: _config.MYSQL_PORT,
  dialect: 'mysql',
  pool: {
    max: 16,
    min: 0,
    idle: 10000
  }
};

// local modules


var db = new _sequelize.Sequelize(_config.MYSQL_DATABASE, _config.MYSQL_USER, _config.MYSQL_PASSWORD, options);

// exports
exports.default = db;