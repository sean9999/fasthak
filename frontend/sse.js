import hak from './hak.js';

/**
 * Filesystem events are exposed as an SSE service at ./.hak/fs/sse
 * This allows you to implement live-reloading logic
 */

const EVENT_NAMESPACE = "fs";

const registerSSE = async () => {
	try {
		const sse = new EventSource(`/${hak.PREFIX}/${EVENT_NAMESPACE}/sse`);
		hak.sse = sse;
		return sse;
	} finally {
		return null;
	}
};

export { hak, registerSSE };