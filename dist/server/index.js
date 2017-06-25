// teyit.link

'use strict';

// core modules

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _http = require('http');

var _http2 = _interopRequireDefault(_http);

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _ratelimit = require('ratelimit.js');

var _bodyParser = require('body-parser');

var _bodyParser2 = _interopRequireDefault(_bodyParser);

var _compression = require('compression');

var _compression2 = _interopRequireDefault(_compression);

var _cookieParser = require('cookie-parser');

var _cookieParser2 = _interopRequireDefault(_cookieParser);

var _cors = require('cors');

var _cors2 = _interopRequireDefault(_cors);

var _express = require('express');

var _express2 = _interopRequireDefault(_express);

var _serveFavicon = require('serve-favicon');

var _serveFavicon2 = _interopRequireDefault(_serveFavicon);

var _expressHbs = require('express-hbs');

var _expressHbs2 = _interopRequireDefault(_expressHbs);

var _lusca = require('lusca');

var _lusca2 = _interopRequireDefault(_lusca);

var _morgan = require('morgan');

var _morgan2 = _interopRequireDefault(_morgan);

var _redis = require('redis');

var _redis2 = _interopRequireDefault(_redis);

var _expressSession = require('express-session');

var _expressSession2 = _interopRequireDefault(_expressSession);

var _config = require('../config');

var _routes = require('./routes');

var _routes2 = _interopRequireDefault(_routes);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

// variables


// local modules


// node modules
var rateLimiter = new _ratelimit.RateLimit(_redis2.default.createClient({
  host: _config.REDIS_HOST,
  port: _config.REDIS_PORT
}), [{
  interval: 3,
  limit: 10
}]);
var limitMiddleware = new _ratelimit.ExpressMiddleware(rateLimiter, {
  ignoreRedisErrors: true
});

// server instance
var Server = function Server() {
  // create epxress app
  var app = (0, _express2.default)();

  // setup express app
  app.set('trust proxy', 1);
  app.disable('x-powered-by');

  // :remote-addr :remote-user :method :url HTTP/:http-version :status :res[content-length] - :response-time ms
  app.use((0, _morgan2.default)('short'));

  // view engine
  app.engine('hbs', _expressHbs2.default.express4({
    defaultLayout: _path2.default.resolve((0, _config.root)('views/layouts/default.hbs')),
    layoutsDir: _path2.default.resolve((0, _config.root)('views/layouts')),
    partialsDir: _path2.default.resolve((0, _config.root)('views/partials'))
  }));
  app.set('view engine', 'hbs');
  app.set('views', _path2.default.resolve((0, _config.root)('views')));

  // prevent abuse with rate limiting
  app.use(limitMiddleware.middleware(function (req, res, next) {
    if (req.query.client) {
      next();
    }
    res.status(429).render('pages/error', {
      code: 429,
      message: 'Too Many Requests'
    });
  }));

  // compress response
  app.use((0, _compression2.default)());

  // session
  app.use((0, _expressSession2.default)({
    cookie: {
      secure: _config.production
    },
    name: _config.SESSION_NAME,
    resave: false,
    saveUninitialized: false,
    secret: _config.SESSION_SECRET
  }));

  // security
  app.use((0, _cors2.default)());
  // app.use(lusca.csrf(true));
  app.use(_lusca2.default.nosniff(true));
  app.use(_lusca2.default.xframe('SAMEORIGIN'));
  app.use(_lusca2.default.xssProtection(true));

  // static folders
  app.use('/public', _express2.default.static(_path2.default.resolve((0, _config.root)('/public')), {
    dotfiles: 'deny',
    index: false,
    maxage: '1d'
  }));

  // favicon
  app.use((0, _serveFavicon2.default)(_path2.default.resolve((0, _config.root)('/public/favicon.ico')), {
    maxage: '1y'
  }));

  // process requests
  app.use(_bodyParser2.default.urlencoded({
    extended: true,
    limit: '1mb'
  }));
  app.use(_bodyParser2.default.json({
    limit: '1mb'
  }));

  // parse cookies
  app.use((0, _cookieParser2.default)());

  // routes
  app.use('/', _routes2.default);

  // catch 404 and forward to error handler
  app.use(function (req, res, next) {
    var err = new Error('Not Found');
    err.status = 404;

    next(err);
  });

  // error handler
  app.use(function (err, req, res) {
    // set locals, only providing error in development
    res.locals.message = err.message;
    res.locals.error = _config.production ? {} : err;

    // print out error details
    // TODO: create an error page and print errors in a proper way
    res.status(err.status || 500).json({
      message: res.locals.message,
      error: res.locals.error,
      code: err.status
    });
  });

  // start servers
  var server = _http2.default.createServer(app).listen(parseInt(_config.PORT, 10), function (err) {
    if (err) {
      console.error(err);
    }

    console.log('Express: Running on port ' + _config.PORT + ' ' + (_config.production ? '(production)' : '(development)') + '.');
  });

  process.on('SIGINT', function () {
    console.log('Express: Server is closing.');

    server.close();
  });
};

// exports
exports.default = Server;