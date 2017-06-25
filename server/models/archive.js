// teyit.link

'use strict';

// node modules
import {Sequelize} from 'sequelize';

// local modules
import db from '../utils/database';

// Archive model
const Archive = db.define('archive', {
  archive_id: {
    type: Sequelize.STRING,
    primaryKey: true,
    autoIncrement: true
  },
  slug: {
    type: Sequelize.STRING
  },
  meta_title: {
    type: Sequelize.STRING
  },
  meta_description: {
    type: Sequelize.TEXT
  },
  request_url: {
    type: Sequelize.STRING
  }
}, {
  timestamps: true,
  updatedAt: false,
  underscored: true
});

// exports
export default Archive;
