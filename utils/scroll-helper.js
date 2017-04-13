const scrollHelper = "var body=document.body,html=document.documentElement,height=Math.max(body.scrollHeight,body.offsetHeight,html.clientHeight,html.scrollHeight,html.offsetHeight);window.scrollTo(0,height);";

module.exports = function(url) {
  if (!url.startsWith("https://www.facebook.com/")) {
    return scrollHelper;
  }
  return "";
};