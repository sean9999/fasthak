/**
 * hak is a singleton
 */
const hak = {
    DEBUG: true,
    PREFIX: ".hak",
    waitFor: (ms) => {
        return new Promise((resolve) => {
            setTimeout(resolve, ms);
        });
    },
    run: (func) => {
        /**
         * run the function when DOM is ready and any other necessary setup is done
         */
        if (document.readyState === 'interactive') {
            func();
        }
        else {
            document.addEventListener('DOMContentLoaded', func);
        }
    },
    sse: null
};
export default hak;
