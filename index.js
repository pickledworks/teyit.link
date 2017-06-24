'use title';
require('dotenv').config();
const express = require('express');
const app = express();
const bodyParser = require('body-parser')
const mysql = require('mysql');
const AWS = require("aws-sdk");
const uuid = require("uuid");
const shortId = require('randomstring');
const dateFormat = require('dateformat');
const common_format = 'dd-mm-yyyy HH:MM:ss';

AWS.config.update({
  accessKeyId: process.env.AWS_ACCESS_KEY,
    secretAccessKey: process.env.AWS_SECRET_KEY,
    region: process.env.AWS_REGION,
});


app.use('/public', express.static('static'));
app.set('view engine', 'pug');

const lambda = new AWS.Lambda();

const connection = mysql.createConnection({
  host: process.env.DB_HOST,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME
});

connection.connect();
app.use(bodyParser.urlencoded({
  extended: true
}));
const createNew = (req, res) => {
  let archive_id = uuid();
  let slug = shortId.generate(7);
  let params = {
    FunctionName: "teyitlink-archive", 
    Payload: JSON.stringify({
      request_url : req.request_url,
      archive_id : archive_id,
    })
  };
  lambda.invoke(params, (err, data) => {
    let payload = JSON.parse(data.Payload);
    if(err){
      res.redirect("/?fail-archive");
    }else{
      let archiveObject = {
        'archive_id' : archive_id,
        'slug' : slug,
        'meta_title' : payload.title,
        'meta_description' : payload.description,
        'request_url' : req.request_url,
        'created_at' : new Date()
      };
      connection.query('insert into archives set  ?',archiveObject, (error, results) => {
        if(req.accepts('json')) {
            res.json({
              error : error,
              data : archiveObject
            });
        }
        else {
            if(error){
              res.redirect("/?fail-create");
            }else{
              res.redirect('/' + slug);
            }
        }

      });

    }
  });

}


app.post('/new', (req,res,next) => {
  req.request_url = req.body.request_url;
  next();
}, createNew);
app.get('/bookmark', (req,res,next) => {
  req.request_url = req.query.request_url;
  next();
}, createNew);

app.get('/', (req, res) => {
  res.render('index');
});

app.get(/^\/(\w{7})?$/, (req, res) => {
  let slug = req.params[0];
  connection.query('SELECT * from archives where slug=? limit 1',[slug], (error, results, fields) => {
    if(results.length > 0){
      results[0]['created_at'] = dateFormat(new Date(results[0]['created_at']), common_format);
      res.render('detail', results[0]);
    }else{
      res.redirect('/?notfound');
    }
  });
});


app.get('/search', (req, res) => {
  if(!req.query.q){
    res.redirect('/'); //Todo view
  }
  let query = req.query.q.replace(/\W/g, '');
  if(query.length < 3){
    res.redirect('/'); //Todo view
  }
  connection.query('SELECT * from archives where request_url like ? or meta_title like ?',['%' +query + '%','%' +query + '%'], (error, results, fields) =>{
    if(results.length < 1){
      res.redirect('/?empty');
    }else{
      results.map(res => {
        res.created_at = dateFormat(new Date(res.created_at), common_format);
      });
      res.render('search',{results : results});
    }

  }); 

});


app.listen(3000, () => {
  console.log('Example app listening on port 3000!')
})


