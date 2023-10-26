<img style="opacity: 0.75" src="fast_hak.png" alt="FastHak" title="FastHak" />


[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/sean9999/fasthak/graphs/commit-activity)

[![Go Reference](https://pkg.go.dev/badge/github.com/sean9999/fasthak.svg)](https://pkg.go.dev/github.com/sean9999/fasthak)

[![Go Report Card](https://goreportcard.com/badge/github.com/sean9999/fasthak)](https://goreportcard.com/report/github.com/sean9999/fasthak)

[![Go version](https://img.shields.io/github/go-mod/go-version/sean9999/fasthak.svg)](https://github.com/sean9999/fasthak)

FastHak is a web server written in Go designed for rapid front-end development. It uses Server Sent Events for live-reload, and automatically injects the necessary javascript files, allowing you to get straight to developing your awesome web app.

It is designed to serve on localhost using HTTPS, because modern web-apps need HTTPS to [do](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#sending_events_from_the_server) [awesome](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API#specifications) [stuff](https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API). It uses the awesome [rec.la](https://www.rec.la/) service for local HTTPS.


## Getting Started

Install fasthak:

```shell
go install github.com/sean9999/fasthak@latest
```

assuming you want to start your server against the ./public subdirectory:

```
$ fasthak -dir=./public # defaults to current dir
```

Use a different port with:

```
$ fasthak -port=12345 # port defaults to 9443
```

It will serve your app on `https://fasthak.rec.la:9443`. You'll want to at least have an `index.html` there.

The client-side code is available at `https://localhost:9443/.hak/js/`. You do not need to include these files in your project. They are embedded in the server itself. You will want to point to them _from_ your project. Example:

```js
import { hak, registerSSE } from "./.hak/js/sse.js";
```

## Bare Minimum front-end Code

The _bare_ bare minimum to see it in action would be to have one `index.html` file that looked like this:

> index.html
```html
<!DOCTYPE html>
<html>
    <head>
        <script defer>
            const sse = new EventSource(`/.hak/fs/sse`);
            sse.addEventListener('fs', (event) => {

                const payload = event.data;
                console.log(payload);

            });
        </script>
    </head>
    <body>
        <p>open console and have a look see</p>
    </body>
</html>
```

But since fasthak provides some convenience functions out the box, you may as well take advantage:

> index.js
```js
import { registerSSE } from "./.hak/js/sse.js";

registerSSE().then(sse => {
    
    //  your business logic
    const biz = fsEvent => {

        const payload = atob(fsEvent.data);
        console.log(payload);
        
    };

    //  attach your business logic to the fs event
    sse.addEventListener('fs', biz);

    //  it's polite to close your connection before leaving
    window.addEventListener("beforeunload", sse.close);

});
```

> index.html
```html
<!DOCTYPE html>
<html>
    <head>
        <script type="module" src="./index.js"></script>
    </head>
    <body>
        <p>open console and have a look see</p>
    </body>
</html>
```

Of course, that only registers the events. You must choose what to do in response. Here's what I do. It's the simplest live-reloader I can think of:

> index.js
```js
import { registerSSE } from "./.hak/js/sse.js";

registerSSE().then(sse => {
    sse.addEventListener('fs', () => {

        //  close the channel. We're about to reload
        sse.close();
        
        window.location.reload();
    
    }); 
});
```

## Why Server Sent Events?

SSE makes more sense than websockets, which is what traditional live-reloaders use. First off, the information does not need to be two-way. Your app has nothing to say to the server. Your server has much to say to your app. The duplex connection that websockets provide are overkill and wasted resources. Fasthak is therefore more correct and more efficient than LiveReload.

Secondly, fasthak provides filesystem events as [DOM Custom Events](https://developer.mozilla.org/en-US/docs/Web/API/CustomEvent/CustomEvent). You choose what you want to do with those events, which in the simplest case is to reload your browser, but could just as easily leverage [Hot Module Replacement](https://blog.bitsrc.io/webpacks-hot-module-replacement-feature-explained-43c13b169986), or some other action that only you can anticipate. LiveReload has some degree of HMR (stylesheets and images are reloaded via javascript), but it's brittle on non-configurable. FastHak gives you total control.


## What's Next?

FastHak could easily be extended to respond to events other than fileSystem events. For example, it could provide introspection capabilities to your otherwise static HTML site, or information about the server such as load and resource usage. There is no reason FastHak could not be used in production.

Additionally, it would be useful to provide hooks for common frameworks, such as React and Vue.

Also, I would like to have a `fasthak init` command that automatically generates the minimal scaffolding needed of your static sites.

Finally, the client code should be available as a browser extension, so that it's injected into the page but not a part of your codebase.

See [issues](https://github.com/sean9999/fasthak/issues) for more.
