/**
 * Some globals and helper functions. 9
 */

const hak = {
    DEBUG: true,
    PREFIX: ".hak",
    EVENT: {
        NAMESPACE: 'fs'
    },
    waitFor: ms => {
        return new Promise((resolve) => {
            setTimeout(resolve, ms);
        })
    },
    run: func => {
        /**
         * run the function when DOM is ready and any other necessary setup is done
         */
        if (document.readyState === 'interactive') {
            func();
        } else {
            document.addEventListener('DOMContentLoaded', func);
        }
    },
    registerSSE: (eventSource) => {
        hak.SSE = eventSource;
    },
    registerDebugNode: (el) => {
        hak.debugNode = el;
    }
};

export default hak;