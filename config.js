// teyit.link

'use strict';

// core modules
import path from 'path';

// node modules
import dotenv from 'dotenv';

// load .env content into process.env
dotenv.config({
  path: path.normalize('.env'),
  silent: true
});
dotenv.load();

// export configs
export const {
  PWD,
  NODE_ENV,
  HOST,
  PORT,
  SESSION_NAME,
  SESSION_SECRET,
  AWS_ACCESS_KEY_ID,
  AWS_SECRET_ACCESS_KEY,
  AWS_REGION,
  MYSQL_HOST,
  MYSQL_PORT,
  MYSQL_USER,
  MYSQL_PASSWORD,
  MYSQL_DATABASE,
  REDIS_HOST,
  REDIS_PORT
} = process.env;

// export custom variables or functions
export const production = process.env.NODE_ENV === 'production';
export const root = dir => path.join(process.env.PWD, dir);
