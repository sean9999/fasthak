import { registerSSE } from "./.hak/js/sse.js";

registerSSE().then(sse => {
    
    //  your business logic
    const biz = fsEvent => {

        const payload = atob(fsEvent.data);
        console.log(payload);
        
    };

    //  attach your business logic to the fs event
    sse.addEventListener('fs', biz);

    //  it's polite to close your connection to SSE before leaving
    window.addEventListener("beforeunload", sse.close);

});
