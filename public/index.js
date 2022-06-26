import { hak, registerSSE } from "/.hak/js/sse.js";

/**
 * A flash that serves as a visual indicator that the page just reloaded
 */
const fouc = () => {
	document.body.classList.add('fouc');
	hak.waitFor(250).then(() => {
		document.body.classList.remove('fouc');
	});
};

const debugLog = (stuff) => {
	localStorage.setItem(`fsEvents/${Date.now()}`, JSON.stringify(stuff));
};

const showDebugLog = () => {
	const events = [];
	for (let i = 0; i < localStorage.length; i++) {
		let k = localStorage.key(i);
		if (k.indexOf("fsEvents/") > -1) {
			events.push(JSON.parse(localStorage.getItem(k)));
		}
	}
	events.reverse();
	document.getElementById('debug').innerText = JSON.stringify(events, null, "\t");
};

hak.run(showDebugLog);
hak.run(() => {
	//	clear bug log button
	document.getElementById('clear-debug').addEventListener("click", ev => {
		ev.preventDefault();
		localStorage.clear();
		showDebugLog();
	});
});

hak.run(fouc);

window.addEventListener("load", () => {
	registerSSE().then(sse => {
		sse.addEventListener('fs', ev => {
			const [fsEventName, filePath] = atob(ev.data).split("\n")

			//	maybe you could do something more intelligent 
			//	than a brute reload
			//	like hot-module reload, etc
			if (hak.DEBUG) {
				debugLog([new Date().toISOString(), fsEventName, filePath]);
			}

			//	allow for debounce time 
			//	since fsEvents seem to happen in clusters
			hak.waitFor(333).then(() => {
				window.location.reload();
			});
		});
	}).catch(console.error);
});