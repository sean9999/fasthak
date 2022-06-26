import hak from './hak.js';

/**
 * Filesystem events are exposed as an SSE service at ./.hak/fs/sse
 * This allows you to implement live-reloading logic
 */

hak.DEBUG = true;

const registerSSE = async () => {
	try {
		const sse = new EventSource(`/${hak.PREFIX}/fs/sse`);
		hak.sse = sse;
		return sse;
	} catch (e) {
		console.error('caught failed registerSSE', e);
	}

};

export { hak, registerSSE };