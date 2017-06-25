// teyit.link

'use strict';

// local modules

var _server = require('./server');

var _server2 = _interopRequireDefault(_server);

var _database = require('./server/utils/database');

var _database2 = _interopRequireDefault(_database);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

// bootstrap
_database2.default.authenticate().then(function () {
  console.log('Sequelize: Authenticated.');

  // start the server
  (0, _server2.default)();
}).catch(function (err) {
  console.error('Sequelize: Unable to authenticate -- ' + err);

  process.exit(1);
});