// teyit.link

'use strict';

// node modules
import {Sequelize} from 'sequelize';

// local modules
import {
  MYSQL_HOST,
  MYSQL_PORT,
  MYSQL_USER,
  MYSQL_PASSWORD,
  MYSQL_DATABASE
} from '../../config';

// variables
const db = new Sequelize(
  MYSQL_DATABASE,
  MYSQL_USER,
  MYSQL_PASSWORD,
  {
    host: MYSQL_HOST,
    port: MYSQL_PORT,
    dialect: 'mysql',
    pool: {
      max: 16,
      min: 0,
      idle: 10000
    }
  }
);

// exports
export default db;
