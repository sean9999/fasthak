import hak from './hak_globals.js';

hak.DEBUG = true;
hak.registerDebugNode(document.getElementById('debug'));

/**
 * Attach an extra event handler to SSE
 * Assume there is a hak.debugNode that is a valid DOM node. Dump events there
 * Potentially do other stuff that is appropriate for debug mode
 * @param eventSource
 * @returns {{element: HTMLElement}}
 */
const debug = () => {

    let debugText = hak.debugNode.innerText = sessionStorage.getItem('debug') || "";
    const handleSSEEvent = (ev) => {
        if (ev.data != null && ev.data !== "null") {
            debugText += "\n" + JSON.stringify(atob(ev.data).split("\n"));
            sessionStorage.setItem('debug', debugText);
            hak.debugNode.innerText = debugText;
        }
    };
    hak.SSE.addEventListener('fs', handleSSEEvent, { passive: true, once: true });
    return { element: hak.debugNode };
};

const clearDebugInfo = () => {
    //sessionStorage.setItem('debug', "");
    sessionStorage.clear();
    hak.debugNode.innerText = "";
};

/**
 * Reload the page whenever an SSE event of type "fs" is received
 * @returns {{eventSource: EventSource}}
 */
const sse = () => {
    const handleSSEEvent = (ev) => {
        if (ev.data) {
            hak.SSE.close();
            location.reload();
        }
    };
    hak.SSE.addEventListener('fs', handleSSEEvent, { passive: true, once: true });

    window.addEventListener('beforeunload', event => {
        hak.SSE.close();
    });

    hak.SSE.addEventListener("message", ev => {
        console.info("sse message", ev.data);
    });
    hak.SSE.addEventListener("error", err => {
        console.error("sse error", err);
        main();
    });
    hak.SSE.addEventListener("open", ev => {
        console.info("sse open");
    });

};

const sayHello = () => {
    document.body.classList.add('green');
    hak.waitFor(333).then(() => {
        document.body.classList.remove('green');
    });
};

const main = () => {

    if (hak.DEBUG) {
        sayHello();
    }

    const embark = () => {
        hak.registerSSE(
            new EventSource(`/${hak.PREFIX}/${hak.EVENT.NAMESPACE}/sse`)
        );
        if (hak.DEBUG) {
            debug();
        }
        sse();
    };

    //  unfortunately we need to force a delay before registration
    let iid = window.setTimeout(embark, 250);

};

hak.run(() => {
    document.getElementById('clear-debug').addEventListener("click", ev => {
        ev.preventDefault();
        clearDebugInfo();
    });
});

//hak.run(main);
window.addEventListener("load", main);
