var func = require("./index").handler;

func({
  "request_url": "http://serafettin.com/",
  "archive_id" : 999999991
}, {}, function(err, res) {
  console.log("[err]", err);
  console.log("[res]", res);
});
