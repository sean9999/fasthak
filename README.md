<img style="opacity: 0.75" src="fast_hak.png" alt="FastHak" title="FastHak" />


[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/sean9999/fasthak/graphs/commit-activity)

[![Go Reference](https://pkg.go.dev/badge/github.com/sean9999/fasthak.svg)](https://pkg.go.dev/github.com/sean9999/fasthak)

[![Go Report Card](https://goreportcard.com/badge/github.com/sean9999/fasthak)](https://goreportcard.com/report/github.com/sean9999/fasthak)

[![Go version](https://img.shields.io/github/go-mod/go-version/sean9999/fasthak.svg)](https://github.com/sean9999/fasthak)

FastHak is a web server written in Go designed for rapid front-end development. It uses Server Sent Events for live-reload, and automatically injects the necessary javascript files, allowing you to get straight to developing your awesome web app.

It is designed to serve on localhost using HTTPS, because modern web-apps need HTTPS to [do](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#sending_events_from_the_server) [awesome](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API#specifications) [stuff](https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API). You need to provide your own certs. I recommend [mkcert](https://github.com/FiloSottile/mkcert#readme).


## Getting Started

Install and run [mkcert](https://github.com/FiloSottile/mkcert), which will drop a localhost.pem and localhost-key.pem in your current directory. You should make sure you understand the [ramifications](https://github.com/FiloSottile/mkcert#installation) of what mkcert does.

Install the certs

```shell
$ mkcert
```

Install fasthak:

```shell
go install github.com/sean9999/fasthak@latest
```

assuming you want to start your server against the ./public subdirectory:

```
$ fasthak -dir=./public
```

which is the equivalent of:

```
$ fasthak \
    -pubkey=./localhost.pem \
    -privkey=./localhost-key.pem \
    -dir=./public \
    -port=9443
```

It will serve your app on `https://localhost:9443`, as you might expect. You'll want to at least have an `index.html` there.

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
import { hak, registerSSE } from "./.hak/js/sse.js";

registerSSE().then(sse => {
    sse.addEventListener('fs', (event) => {

        const payload = event.data;
        console.log(payload);
    
    }); 
});
```
> index.html
```html
<!DOCTYPE html>
<html>
    <head>
        <script defer src="./index.js"></script>
    </head>
    <body>
        <p>open console and have a look see</p>
    </body>
</html>
```

Of course, that only registers the events. You must choose what to do in response. Here's what I do. It's the simplest live-reloader I can think of:

> index.js
```js
registerSSE().then(sse => {
    sse.addEventListener('fs', (event) => {

        //  close the channel
        //  not strictly necessary, but polite
        sse.close();

        //  reload browser window
        //  sse will reconnect on DOMContentLoaded
        window.location.reload();
    
    }); 
});

```

## Why Server Sent Events?

SSE makes more sense than websockets, which is what traditional live-reloaders use. First off, the information does not need to be two-way. Your app has nothing to say to the server. Your server has much to say to your app. The duplex connection that websockets provide are overkill and wasted resources. Fasthak is therefore more correct and more efficient than LiveReload.

Secondly, fasthak provides filesystem events as [DOM Custom Events](https://developer.mozilla.org/en-US/docs/Web/API/CustomEvent/CustomEvent). You choose what you want to do with those events, which in the simplest case is to reload your browser, but could just as easily leverage [Hot Module Replacement](https://blog.bitsrc.io/webpacks-hot-module-replacement-feature-explained-43c13b169986), or some other action that only you can anticipate. LiveReload has some degree of HMR (stylesheets and images are reloaded via javascript), but it's brittle on non-configurable. FastHak gives you total control.


## Should I switch from LiveReload?

Meh, probably not, if you can't easily see how it would improve your workflow. This is a very niche improvement, since performance optimisation rarely matters in development mode. For me, FastHak was mainly an excuse to write a server in Go. That said, I use it all the time for net new web projects, like [my blog](https://www.seanmacdonald.ca).

## How does it work

[Rebouncer](https://github.com/sean9999/rebouncer) does all the heavy-lifting. Fasthak simply wraps a static server around it.

## What's Next?

Due to it's design, FastHak could easily be extended to respond to events other than fileSystem events. For example, it could provide introspection capabilities to your otherwise static HTML site, or information about the server such as load and resource usage. There is no reason FastHak could not be used in production.

Additionally, it would be useful to provide hooks for common frameworks, such as React and Vue.

Furthermore, I would like to have a `fasthak init` command that automatically generates the minimal scaffolding needed of your static sites, and possibly generates certs too, freeing you from having to wrestle with mkcert.

Finally, the client code should be available as a browser extension, so that it's injected into the page but not a part of your codebase.

See [issues](https://github.com/sean9999/fasthak/issues) for more.
