

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
const linkRegex = /(http(s)?:\/\/.)?(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&\/\/=;]*)/g;

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
 * creates a new song info object
 * @param {String} circle
 * @param {String} title 
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
 * @description matches a regular expression and returns subgroup at index 'pos'
 * if the match fails, it will return 'or'
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

/**
 * @description creates a regular expression for searching for song information in a video description
 * @param {String} field field to search for
 * @returns {RegExp}
 */
function newParseReg(field) {
	let r = String.raw`^.*(\[?${field}\]?)([\s]|\.){0,30}:[\s]{0,5}([^ ].*)$`;
	return RegExp(r, 'mi');
}

/**
 * @description parses song information from a pages 'pre' elements and returns an array of song infos
 * @returns {SongInfo} song information
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

/**
 * @description creates a link to a youtube search for the given text
 * will return undefined if the supplied text is ""
 * @param {String} text text to return
 * @returns {HTMLElement} link
 */
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

/**
 * @description makes links leading to youtube searches for song information
 * parsed from video descriptions.
 */
function makeSearchButtons() {
	let songs = parseSongInformation();
	let container = document.querySelector(".search");
	let join = x => x.filter(x => x != "").join(" - ");

	for (let x of songs) {
		console.log(x);

		Array.from([
			document.querySelector(".webpage-url h3").innerHTML, // include video title as a search query
			join([x.title]),
			join([x.circle, x.title]),
			join([x.circle, x.album]),
			join([x.circle, x.album, x.title]),
		]
		.reduce((x, y) => { x.set(y, true); return x}, new Map()) // Remove duplicates using a Map
		.keys())               // Obtain the keys iterator and convert it to an array
		.map(createSearchLink) // create links for each search query
		.filter(x => !!x)      // Filter out any undefined elements
		.forEach(x => container.insertAdjacentElement("afterend", x)); // append buttons after search form
	}
}

window.addEventListener("load", function() {
	replaceURLs();
	makeSearchButtons();
});
