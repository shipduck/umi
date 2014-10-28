//https://konklone.com/post/github-pages-now-supports-https-so-use-it
var host = "shipduck.github.io";
if ((host == window.location.host) && (window.location.protocol != "https:")) {
  window.location.protocol = "https";
}

var videoList= document.getElementsByTagName('video');
for(var i = 0 ; i < videoList.length ; ++i) {
  var video = videoList[i];
  video.addEventListener('click', function() {
	  if(video.paused) {
	    video.play();
	  } else {
	    video.pause();
	  }
  }, false);
}
