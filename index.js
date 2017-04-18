// teyit.link

'use strict';

// local modules
import Server from './server';
import db from './server/utils/database';

// bootstrap
db.authenticate().then(() => {
  console.log('Sequelize: Authenticated.');

  // start the server
  Server();
}).catch(err => {
  console.error(`Sequelize: Unable to authenticate -- ${err}`);

  process.exit(1);
});
