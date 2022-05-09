import { hak, registerSSE } from "./.hak/js/sse.js";

const fouc = () => {
	document.body.classList.add('fouc');
	hak.waitFor(250).then(() => {
		document.body.classList.remove('fouc');
	});
};

window.addEventListener("load", () => {
	registerSSE().then(sse => {
		sse.addEventListener('fs', ev => {
			const [fsEventName, filePath] = atob(ev.data).split("\n")

			//	maybe you could do something more intelligent 
			//	than a brute reload
			//	like hot-module reload, etc
			console.log({ fsEventName, filePath });

			//	allow for debounce time 
			//	since fsEvents seem to happen in clusters
			hak.waitFor(333).then(() => {
				window.location.reload();
			});
		});
	}).catch(console.error);
	fouc();
});