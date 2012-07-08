/*globals init load populate up clickDir clickFile addToPlaylist addAll play next*/
var root = "/f/";
var path = [];
function init() {
	load(path);
	$('#player').bind('ended', next);
	$('#addall').click(addAll);
	$('#next').click(next);
}
function load(path)  {
	$.ajax({
		url:root+path.join('/'),
		dataType:"json",
		success: populate
	});
}
function populate(files) {
	var $b = $('#browser').empty();
	function add(i, f) {
		if (f.name[0] == '.') return;
		var dir = f.isDir;
		var cl = dir ? "dir" : "file";
		var evenOrOdd = (i %2 === 0) ? 'even' : 'odd';

		$('<a></a>').text(f.name).data('file', f).data('evenOrOdd', evenOrOdd)
			.addClass(cl).appendTo($b)
			.click(dir?clickDir:clickFile);
	}
	files.sort(function(a, b) {
		a = a.name.toLowerCase();
		b = b.name.toLowerCase();
		if (a > b) return 1;
		if (a < b) return -1;
		return 0;
	});
	$b.append(up());
	$.each(files, add);
}
function up() {
	return $('<a class="dir">..</a>').click(function() {
		path.pop();
		load(path);
	});
}
function clickDir(e) {
	path.push($(e.target).data('file').name);
	load(path);
}
function clickFile(e) {
  var className = $('#playlist > a:last').hasClass('even') ? 'odd' : 'even';
	addToPlaylist($(e.target).data('file'), className);
}
function addToPlaylist(f, eoo) {
	var $p = $('#playlist');
	var playnow = ($p.find('a').length === 0);
	var $d = $('<a></a>').text(f.name).data('file', f).data('path', path.map(function(i) { return i; })).addClass(eoo)
		.appendTo($p)
		.click(function(e) { play(e.target); });
	if (playnow) $d.click();
}
function addAll() {
  var hasEven, cls;
  hasEven = $('#playlist > a:last').hasClass('even');

	$('#browser a.file').each(function(i, e) {
	  cls = hasEven ? (i % 2 === 0) ? 'odd' : 'even' : (i % 2 === 0) ? 'even' : 'odd';
		addToPlaylist($(e).data('file'), cls);
	});
}
function play(el) {
	var name = $(el).data('file').name;
	var pth = $(el).data('path');
	var url = root+pth.join('/')+'/'+name;
	$('#playlist a').removeClass('playing');
	$(el).addClass('playing');
	$('#player').attr('src', url);
}
function next() {
	var $next = $('#playlist a.playing').next();
	if ($next.length) {
		setTimeout($next.click(), 2000);
	}
}

$(document).ready(function() {
	init();
	var audio,
	    loadingIndicator,
	    positionIndicator,
	    timeleft,
	    loaded = false,
	    manualSeek = false;

	audio = $('.player audio').get(0);
	loadingIndicator = $('.player #loading');
	positionIndicator = $('.player #handle');
	timeleft = $('.player #timeleft');

	if (audio && (audio.buffered !== undefined) && (audio.buffered.length !== 0)) {
	  $(audio).bind('progress', function() {
	    var loaded = parseInt(((audio.buffered.end(0) / audio.duration) * 100), 10);
	    loadingIndicator.css({width: loaded + '%'});
	  });
	}
	else {
	  loadingIndicator.remove();
	}

	$(audio).bind('timeupdate', function() {

	  var rem = parseInt(audio.duration - audio.currentTime, 10),
	  pos = (audio.currentTime / audio.duration) * 100,
	  mins = Math.floor(rem/60,10),
	  secs = rem - mins*60;

	  timeleft.text('-' + mins + ':' + (secs > 9 ? secs : '0' + secs));
	  if (!manualSeek) { positionIndicator.css({left: pos + '%'}); }
	  if (!loaded) {
	    loaded = true;

	    $('.player #gutter').slider({
	      value: 0,
	      step: 0.01,
	      orientation: "horizontal",
	      range: "min",
	      max: audio.duration,
	      animate: true,
	      slide: function() {
	        manualSeek = true;
	      },
	      stop:function(e,ui) {
	        manualSeek = false;
	        audio.currentTime = ui.value;
	      }
	    });
	  }

	});

	$(audio).bind('play',function() {
	  $("#playtoggle").addClass('playing');
	}).bind('pause', function() {
	  $("#playtoggle").removeClass('playing');
	});

	$("#playtoggle").click(function() {
	  if (audio.paused) { audio.play(); }
	  else { audio.pause(); }
	});

});