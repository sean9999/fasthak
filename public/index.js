import { hak, registerSSE } from "./.hak/js/sse.js";
import {debug, render, clearDebugInfo} from "./js/debug.js";

//	a function that provides a quick, easily identifiable clue that the window just reloaded
const fouc = () => {
	document.body.classList.add('fouc');
	hak.waitFor(250).then(() => {
		document.body.classList.remove('fouc');
	});
};

window.addEventListener("load", () => {

	registerSSE().then(sse => {

		//	these may be useful for debugging
		sse.addEventListener('open', (stuff) => {console.log('sse open',stuff)});
		sse.addEventListener('message', (stuff) => {console.log('sse message',stuff)});
		sse.addEventListener('error', (stuff) => {console.error('sse error',stuff)});

		sse.addEventListener('fs', ev => {
			const [macroEvent, microEvent, filePath] = ev.data.split("\n")

			console.log({macroEvent, microEvent, filePath});


			//	maybe you could do something more intelligent 
			//	than a brute reload
			//	like hot-module reload, etc
			//console.log({ fsEventName, filePath });
			//debugElement.innerHTML += JSON.stringify({macroEvent, microEvent, filePath}) + '<br />';
		
			debug('fsevents', {microEvent, filePath});

			//	allow for debounce time 
			//	since fsEvents seem to happen in clusters
			
			// hak.waitFor(1000).then(() => {
			// 	if (sse.readyState === 1) {
			// 		sse.close();
			// 	}
			// 	window.location.reload();
			// });
		

			window.addEventListener('beforeunload',(ev) => {
				//ev.preventDefault();
				if (sse.readyState === 1) {
					sse.close();
				}
			});
			
		});
	}).catch(console.error);
	fouc();
	render('fsevents');
});
