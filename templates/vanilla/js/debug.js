const debugElement = document.getElementById('debug');
const clearDebugButton = document.getElementById('clear-debug');
clearDebugButton.addEventListener('click', () => {
	clearDebugInfo('fsevents');
});

//	save to localStorage
const persist = (groupId, thing) => {
	const key = `debug/${groupId}/${Date.now()}`;
	localStorage.setItem(key, JSON.stringify(thing));
};

const getDateFromKey = (key) => {
	return Number(key.split('/').pop());
};

const compareByKey = (a, b) => {
	if (a.key > b.key) {
		return 1;
	}
	if (a.key < b.key) {
		return -1;
	}
	return 0;
};

const getSortedEventsFromStorage = (groupId) => {
	const events = [];
	const keys = Object.keys(localStorage);
	for (let key of keys) {
		if (key.startsWith(`debug/${groupId}/`)) {
			let event = localStorage.getItem(key); 
			events.push({"key": getDateFromKey(key), event});
		}
	}
	events.sort(compareByKey);
	return events;
};

//	render to the DOM
const render = (groupId) => {
	const groupNode = debugElement.querySelector(`#${groupId} ul`);
	const keys = Object.keys(localStorage);
	const events = getSortedEventsFromStorage(groupId);
	events.forEach(event => {
		const li = document.createElement('li');
		const friendlyDate = new Date(event.key).toLocaleTimeString('en-US');
		li.innerText = friendlyDate + "\t" + event.event;
		groupNode.appendChild(li);
	});
}

//	call this when there is an event to log
const debug = (groupId, stuff) => {
	persist(groupId, stuff);
	render(groupId);
};

const clearDebugInfo = (groupId) => {
	const groupNode = debugElement.querySelector(`#${groupId} ul`);
	const keys = Object.keys(localStorage);
	for (let key of keys) {
		if (key.startsWith(`debug/${groupId}/`)) {
			localStorage.removeItem(key);
		}
	}
	groupNode.innerHTML = "";
};

export { debug, render, persist, clearDebugInfo }
