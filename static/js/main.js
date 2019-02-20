

/**
 * @description Sets the CSS cookie and reloads the window for it to take effect
 * @param {string} name the string to set the css cookie to
 */
function setCSS(name) {
	let d = new Date();
	document.cookie = `css=${name}; path=/; expires=${new Date(d.getFullYear() + 9, d.getMonth(), d.getDay()).toUTCString()}`;
	window.location.reload();
}

// linkRegex is a regular expression for matching links
const linkRegex = /(http(s)?:\/\/.)?(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)/g;

/**
 * @description Replaces all URLs in the video descriptions with hyperlinks
 */
function replaceURLs() {
	for (let x of document.getElementsByTagName("pre")) {
		x.innerHTML = x.innerHTML.replace(linkRegex, (match) => {
			return `<a href='${match}'>${match}</a>`;
		});
	}
}

/** 
 * @typedef {Object} SongInfo
 * @property {String} circle
 * @property {String} title
 * @property {String} album
*/

/**
 * @param {String} circle
 * @param {String} title 
 * @description creates a new song info object
 * @returns {SongInfo}
 */
function newSongInfo(circle="", title="", album="") {
	return {
		circle,
		title,
		album,
	}
}

/**
 * @param {RegExp} regex regular expression to use in matching
 * @param {String} text string to search
 * @param {Number} pos number of the capture group
 * @param {String} or text to return if the match fails
 * @returns {String}
 */
function matchGroupOr(regex, text, pos, or="") {
	console.log("match");
	console.log(regex);
	let match = text.match(regex);
	console.log(match);
	if (!match) {
		return or;
	}
	return match[pos];
}

function newParseReg(field) {
	let r = String.raw`^.*(\[?${field}\]?)([\s]|\.){0,30}:[\s]{0,5}([^ ].*)$`;
	return RegExp(r, 'mi');
}

/**
 * @description parses song information from the descriptions and returns an array of song infos
 */
function parseSongInformation() {
	let songs = [];
	for (let x of document.getElementsByTagName("pre")) {
		let songInfo = newSongInfo();

		function getField(f) {
			return matchGroupOr(newParseReg(f), x.innerHTML, 3);
		}
		
		songInfo.title = getField("Title|Track");
		songInfo.circle = getField("Circle|Performed by");
		songInfo.album = getField("Album");

		console.log(songInfo)
		songs.push(songInfo);
	}
	return songs;
}

function createSearchLink(text) {
	if (text == "") 
		return null;
	let link = document.createElement("a");
	link.href = `https://youtube.com/search?q=${encodeURIComponent(text)}`;
	// link.innerHTML = "[" + text + "]";
	link.appendChild(document.createTextNode("[" + text + "]"));
	link.style.marginRight = "10px";
	link.target = "_blank";
	return link;
}

window.addEventListener("load", function() {
	replaceURLs();
	let songs = parseSongInformation();
	for (let x of songs) {
		console.log(x);

		let join = (x) => x.filter(x => x != "").join(" - ");
		
		let elements = Array.from([
			document.querySelector(".webpage-url h3").innerHTML,
			join([x.title]),
			join([x.circle, x.title]),
			join([x.circle, x.album]),
			join([x.circle, x.album, x.title]),
		]
		.reduce((x, y) => { x.set(y, true); return x}, new Map())
		.keys())
		.map(createSearchLink)
		.filter(x => !!x);
		
		let container = document.querySelector(".search");
		for (let e of elements) {
			container.insertAdjacentElement("afterend", e);
		}
	}
});
