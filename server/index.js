// teyit.link

'use strict';

// core modules
import http from 'http';
import path from 'path';

// node modules
import {ExpressMiddleware, RateLimit} from 'ratelimit.js';
import bodyParser from 'body-parser';
import compression from 'compression';
import cookieParser from 'cookie-parser';
import cors from 'cors';
import express from 'express';
import favicon from 'serve-favicon';
import hbs from 'express-hbs';
import lusca from 'lusca';
import morgan from 'morgan';
import redis from 'redis';
import session from 'express-session';

// local modules
import {
  PORT,
  SESSION_NAME,
  SESSION_SECRET,
  REDIS_HOST,
  REDIS_PORT,
  production,
  root
} from '../config';
import routes from './routes';

// variables
const rateLimiter = new RateLimit(
  redis.createClient({
    host: REDIS_HOST,
    port: REDIS_PORT
  }),
  [{
    interval: 3,
    limit: 10
  }]
);
const limitMiddleware = new ExpressMiddleware(rateLimiter, {
  ignoreRedisErrors: true
});

// server instance
const Server = () => {
  // create epxress app
  const app = express();

  // setup express app
  app.set('trust proxy', 1);
  app.disable('x-powered-by');

  // :remote-addr :remote-user :method :url HTTP/:http-version :status :res[content-length] - :response-time ms
  app.use(morgan('short'));

  // view engine
  app.engine('hbs', hbs.express4({
    defaultLayout: path.resolve(root('views/layouts/default.hbs')),
    layoutsDir: path.resolve(root('views/layouts')),
    partialsDir: path.resolve(root('views/partials'))
  }));
  app.set('view engine', 'hbs');
  app.set('views', path.resolve(root('views')));

  // prevent abuse with rate limiting
  app.use(limitMiddleware.middleware((req, res, next) => {
    res.status(429).render('pages/error', {
      code: 429,
      message: 'Too Many Requests'
    });
  }));

  // compress response
  app.use(compression());

  // session
  app.use(session({
    cookie: {
      secure: production
    },
    name: SESSION_NAME,
    resave: false,
    saveUninitialized: false,
    secret: SESSION_SECRET
  }));

  // security
  app.use(cors());
  // app.use(lusca.csrf(true));
  app.use(lusca.nosniff(true));
  app.use(lusca.xframe('SAMEORIGIN'));
  app.use(lusca.xssProtection(true));

  // static folders
  app.use('/public', express.static(path.resolve(root('/public')), {
    dotfiles: 'deny',
    index: false,
    maxage: '1d'
  }));

  // favicon
  app.use(favicon(path.resolve(root('/public/favicon.ico')), {
    maxage: '1y'
  }));

  // process requests
  app.use(bodyParser.urlencoded({
    extended: true,
    limit: '1mb'
  }));
  app.use(bodyParser.json({
    limit: '1mb'
  }));

  // parse cookies
  app.use(cookieParser());

  // routes
  app.use('/', routes);

  // catch 404 and forward to error handler
  app.use((req, res, next) => {
    const err = new Error('Not Found');
    err.status = 404;

    next(err);
  });

  // error handler
  app.use((err, req, res) => {
    // set locals, only providing error in development
    res.locals.message = err.message;
    res.locals.error = production ? {} : err;

    // print out error details
    // TODO: create an error page and print errors in a proper way
    res.status(err.status || 500).json({
      message: res.locals.message,
      error: res.locals.error,
      code: err.status
    });
  });

  // start servers
  const server = http.createServer(app).listen(parseInt(PORT, 10), err => {
    if (err) {
      console.error(err);
    }

    console.log(`Express: Running on port ${PORT} ${production ? '(production)' : '(development)'}.`);
  });

  process.on('SIGINT', () => {
    console.log('Express: Server is closing.');

    server.close();
  });
};

// exports
export default Server;
