// teyit.link

'use strict';

// node modules
import moment from 'moment';

// local modules
import Archive from '../models/archive';

// configuration
moment.locale('tr');

const archive = {
  // /new
  new: {
    post: (req, res, next) => {
      req.request_url = req.body.request_url;

      next();
    }
  },

  // /search (/?q=)
  search: {
    get: (req, res) => {
      if (!req.query.q) {
        res.status(400).redirect('/');
      }

      const query = req.query.q.replace(/\W/g, '');

      if (query.length < 3) {
        res.status(400).redirect('/');
      }

      Archive.findAll({
        where: {
          $or: {
            request_url: {
              $like: `%${query}%`
            },
            meta_title: {
              $like: `%${query}%`
            }
          }
        }
      }).then(results => {
        if (results.length < 1) {
          res.status(204).redirect('/?empty');
        }
        else {
          res.status(200).render('pages/search', {
            results: results.map(result => result.dataValues)
          });
        }
      });
    }
  },

  // /slug (/^\/(\w{7})?$/)
  slug: {
    get: (req, res) => {
      const slug = req.params[0];

      Archive.findOne({
        where: {
          slug
        }
      }).then(result => {
        if (result && 'dataValues' in result) {
          result.dataValues.created_at = moment(result.dataValues.created_at).format('LL LTS'); // 01 Ocak 2017 12:34:56

          res.status(200).render('pages/detail', result.dataValues);
        }
        else {
          res.status(404).redirect('/?notfound');
        }
      });
    }
  }
};

// exports
export default archive;
