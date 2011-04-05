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
		if (f.Name[0] == '.') return;
		var dir = (f.Mode & 040000);
		var cl = dir ? "dir" : "file";
		var evenOrOdd = (i %2 === 0) ? 'even' : 'odd';
		
		$('<a></a>').text(f.Name).data('file', f).data('evenOrOdd', evenOrOdd)
			.addClass(cl).appendTo($b)
			.click(dir?clickDir:clickFile);
	}
	files.sort(function(a, b) {
		a = a.Name.toLowerCase();
		b = b.Name.toLowerCase();
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
	path.push($(e.target).data('file').Name);
	load(path);
}
function clickFile(e) {
	addToPlaylist($(e.target).data('file'), $(e.target).data('evenOrOdd'));
}
function addToPlaylist(f, eoo) {
  console.log(eoo);
	var $p = $('#playlist');
	var playnow = ($p.find('a').length === 0);
	var $d = $('<a></a>').text(f.Name).data('file', f).data('path', path.map(function(i) { return i; })).addClass(eoo)
		.appendTo($p)
		.click(function(e) { play(e.target); });
	if (playnow) $d.click();
}
function addAll() {
	$('#browser a.file').each(function(i, e) {
		addToPlaylist($(e).data('file'), (i % 2) === 0 ? 'even' : 'odd');
	});
}
function play(el) {
	var name = $(el).data('file').Name;
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