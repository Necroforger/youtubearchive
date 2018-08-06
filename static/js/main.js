/**
 * @description Sets the CSS cookie and reloads the window
 * 				for it to take effect
 * @param {string} name the string to set the css cookie to
 */
function setCSS(name) {
	document.cookie = `css=${name}; path=/; expires=09/09/9999`;
	window.location.reload();
}




