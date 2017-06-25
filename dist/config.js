// teyit.link

'use strict';

// core modules

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.root = exports.production = exports.REDIS_PORT = exports.REDIS_HOST = exports.MYSQL_DATABASE = exports.MYSQL_PASSWORD = exports.MYSQL_USER = exports.MYSQL_PORT = exports.MYSQL_HOST = exports.AWS_REGION = exports.AWS_SECRET_ACCESS_KEY = exports.AWS_ACCESS_KEY_ID = exports.SESSION_SECRET = exports.SESSION_NAME = exports.PORT = exports.HOST = exports.NODE_ENV = exports.PWD = undefined;

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _dotenv = require('dotenv');

var _dotenv2 = _interopRequireDefault(_dotenv);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

// load .env content into process.env
_dotenv2.default.config({
  path: _path2.default.normalize('.env'),
  silent: true
});

// node modules

_dotenv2.default.load();

// export configs
var _process$env = process.env,
    PWD = _process$env.PWD,
    NODE_ENV = _process$env.NODE_ENV,
    HOST = _process$env.HOST,
    PORT = _process$env.PORT,
    SESSION_NAME = _process$env.SESSION_NAME,
    SESSION_SECRET = _process$env.SESSION_SECRET,
    AWS_ACCESS_KEY_ID = _process$env.AWS_ACCESS_KEY_ID,
    AWS_SECRET_ACCESS_KEY = _process$env.AWS_SECRET_ACCESS_KEY,
    AWS_REGION = _process$env.AWS_REGION,
    MYSQL_HOST = _process$env.MYSQL_HOST,
    MYSQL_PORT = _process$env.MYSQL_PORT,
    MYSQL_USER = _process$env.MYSQL_USER,
    MYSQL_PASSWORD = _process$env.MYSQL_PASSWORD,
    MYSQL_DATABASE = _process$env.MYSQL_DATABASE,
    REDIS_HOST = _process$env.REDIS_HOST,
    REDIS_PORT = _process$env.REDIS_PORT;

// export custom variables or functions

exports.PWD = PWD;
exports.NODE_ENV = NODE_ENV;
exports.HOST = HOST;
exports.PORT = PORT;
exports.SESSION_NAME = SESSION_NAME;
exports.SESSION_SECRET = SESSION_SECRET;
exports.AWS_ACCESS_KEY_ID = AWS_ACCESS_KEY_ID;
exports.AWS_SECRET_ACCESS_KEY = AWS_SECRET_ACCESS_KEY;
exports.AWS_REGION = AWS_REGION;
exports.MYSQL_HOST = MYSQL_HOST;
exports.MYSQL_PORT = MYSQL_PORT;
exports.MYSQL_USER = MYSQL_USER;
exports.MYSQL_PASSWORD = MYSQL_PASSWORD;
exports.MYSQL_DATABASE = MYSQL_DATABASE;
exports.REDIS_HOST = REDIS_HOST;
exports.REDIS_PORT = REDIS_PORT;
var production = exports.production = process.env.NODE_ENV === 'production';
var root = exports.root = function root(dir) {
  return _path2.default.join(process.env.PWD, dir);
};