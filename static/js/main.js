/**
 * @description Sets the CSS cookie and reloads the window for it to take effect
 * @param {string} name the string to set the css cookie to
 */
function setCSS(name) {
	let d = new Date();
	document.cookie = `css=${name}; path=/; expires=${new Date(d.getFullYear() + 9, d.getMonth(), d.getDay()).toUTCString()}`;
	window.location.reload();
}




