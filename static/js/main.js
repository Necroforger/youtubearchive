

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
const linkRegex = /(http(s)?:\/\/.)?(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)/;

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

window.addEventListener("load", function() {
	replaceURLs();
});
