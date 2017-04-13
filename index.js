require('dotenv').config();
var express = require('express');
var app = express();
var bodyParser = require('body-parser')
var mysql = require('mysql');
var AWS = require("aws-sdk");
var uuid = require("uuid");
var shortId = require('randomstring');
var dateFormat = require('dateformat');
var common_format = 'dd-mm-yyyy HH:MM:ss';

AWS.config.update({
  accessKeyId: process.env.AWS_ACCESS_KEY,
    secretAccessKey: process.env.AWS_SECRET_KEY,
    region: process.env.AWS_REGION,
});


app.use('/public', express.static('static'));
app.set('view engine', 'pug');

var lambda = new AWS.Lambda();

var connection = mysql.createConnection({
  host: process.env.DB_HOST,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME
});

connection.connect();
app.use(bodyParser.urlencoded({
  extended: true
}));
const createNew = function (req, res) {
  var archive_id = uuid();
  var slug = shortId.generate(7);
  var params = {
    FunctionName: "teyitlink-archive", 
    Payload: JSON.stringify({
      request_url : req.request_url,
      archive_id : archive_id,
    })
  };
  lambda.invoke(params, function(err, data) {
    var payload = JSON.parse(data.Payload);
    if(err){
      res.redirect("/?fail-archive");
    }else{
      var archiveObject = {
        'archive_id' : archive_id,
        'slug' : slug,
        'meta_title' : payload.title,
        'meta_description' : payload.description,
        'request_url' : req.request_url,
        'created_at' : new Date()
      };
      connection.query('insert into archives set  ?',archiveObject, function (error, results) {
        if(error){
          res.redirect("/?fail-create");
        }else{
          res.redirect('/' + slug);
        }

      });

    }
  });

}

app.post('/new', function(req,res,next){
  req.request_url = req.body.request_url;
  next();
}, createNew);
app.get('/bookmark', function(req,res,next){
  req.request_url = req.query.request_url;
  next();
}, createNew);

app.get('/', function (req, res) {
  res.render('index');
});

app.get(/^\/(\w{7})?$/,function (req, res) {
  var slug = req.params[0];
  connection.query('SELECT * from archives where slug=? limit 1',[slug], function (error, results, fields) {
    if(results.length > 0){
      results[0]['created_at'] = dateFormat(new Date(results[0]['created_at']), common_format);
      res.render('detail', results[0]);
    }else{
      res.redirect('/?notfound');
    }
  });
});


app.get('/search', function (req, res) {
  if(!req.query.q){
    res.redirect('/'); //Todo view
  }
  var query = req.query.q.replace(/\W/g, '');
  if(query.length < 3){
    res.redirect('/'); //Todo view
  }
  connection.query('SELECT * from archives where request_url like ? or meta_title like ?',['%' +query + '%','%' +query + '%'], function (error, results, fields) {
    if(results.length < 1){
      res.redirect('/?empty');
    }else{
      results.map(function(res){
        res.created_at = dateFormat(new Date(res.created_at), common_format);
      });
      res.render('search',{results : results});
    }

  }); 

});


app.listen(3000, function () {
  console.log('Example app listening on port 3000!')
})


