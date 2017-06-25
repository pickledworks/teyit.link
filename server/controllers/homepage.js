// teyit.link

'use strict';

const homepage = {
  // permanently redirect all unauthorized requests and non-existing routes to hompage
  all: (req, res) => {
    res.redirect(301, '/');
  },

  // /
  get: (req, res) => {
    res.status(200).render('pages/homepage');
  }
};

// exports
export default homepage;
