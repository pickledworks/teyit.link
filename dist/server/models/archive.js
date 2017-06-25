// teyit.link

'use strict';

// node modules

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _sequelize = require('sequelize');

var _database = require('../utils/database');

var _database2 = _interopRequireDefault(_database);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

// Archive model
var Archive = _database2.default.define('archive', {
  archive_id: {
    type: _sequelize.Sequelize.STRING,
    primaryKey: true,
    autoIncrement: true
  },
  slug: {
    type: _sequelize.Sequelize.STRING
  },
  meta_title: {
    type: _sequelize.Sequelize.STRING
  },
  meta_description: {
    type: _sequelize.Sequelize.TEXT
  },
  request_url: {
    type: _sequelize.Sequelize.STRING
  }
}, {
  timestamps: true,
  updatedAt: false,
  underscored: true
});

// exports


// local modules
exports.default = Archive;