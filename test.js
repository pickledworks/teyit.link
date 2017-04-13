var func = require('./index').handler;

func({
    "request_url": "http://teyit.org",
    "archive_id" : 99911
},{},function(err,res){
    console.log(err,res);
});
