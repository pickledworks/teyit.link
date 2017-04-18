// teyit.link

'use strict';

// node modules
import {Router} from 'express';

// local modules
import archive from './controllers/archive';
import homepage from './controllers/homepage';
import lambda from './utils/lambda';

// variables
const router = Router();

// routes
router.route('/').get(homepage.get);
router.route('/').all(homepage.all);
router.route(['/add', '/bookmark', '/new']).post(archive.new.post, lambda);
router.route(['/add', '/bookmark', '/new']).all(homepage.all);
router.route('/search').get(archive.search.get);
router.route('/search').all(homepage.all);
router.route(/^\/(\w{7})?$/).get(archive.slug.get);
router.route(/^\/(\w{7})?$/).all(homepage.all);

// exports
export default router;
