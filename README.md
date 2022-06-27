<img src="fast_hak.png" alt="FastHak" title="FastHak" />

FastHak is a web server written in Go designed for rapid front-end development. It uses Server Sent Events for live-reload, and automatically injects the necessary javascript files, allowing you to get straight to developing your awesome web app.

It is designed to serve on localhost using HTTPS, because modern web-apps need HTTPS to [do](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#sending_events_from_the_server) [awesome](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API#specifications) [stuff](https://developer.mozilla.org/en-US/docs/Web/API/WebRTC_API). You need to provide your own certs. I recommend [mkcert](https://github.com/FiloSottile/mkcert#readme).

FastHak is a small binary you can drop directly in your project root and invoke like so:

```
$ fasthak -dir=public --port=9443
```

Which will serve your app on `https://localhost:9443`.

The client-side code is available at `https://localhost:9443/.hak/js/`. You do not need to include these files in your project. They are embedded in the server itself. But you will want to point to them _from_ your project.

## Why Server Sent Events?

SSE makes more sense than websockets, which is what traditional live-reloaders use. First off, the information does not need to be two-way. Your app has nothing to say to the server. Your server has much to say to your app. The duplex connection that websockets provide are overkill and wasted resources. Fasthak is therefore more correct and more efficient than LiveReload.

Secondly, fasthak provides filesystem events as Javascript Events. You choose what you want to do with those events, which in the simplest case is to reload your browser, but could just as easily leverage [Hot Module Replacement](https://blog.bitsrc.io/webpacks-hot-module-replacement-feature-explained-43c13b169986). LiveReload has some degree of this (stylesheets and images are reloaded via javascript), but it's brittle on non-configurable. FastHak gives you total control.

## Should I switch from LiveReload?

Meh, probably not, if you can't easily see how it would improve your workflow. This is a very niche improvement, since performance optimisation rarely matters in development mode. For me, FastHak was mainly an excuse to write a server in Go. That said, I use it all the time for net new web projects, like [my blog](https://www.seanmacdonald.ca).

## What's Next?

Due to it's design, FastHak could easily be extended to respond to events other than fileSystem events. For example, it could provide introspection capabilities to your otherwise static HTML site, or information about the server such as load and resource usage. There is no reason FastHak could not be used in production.

Additionally, it would be useful to provide hooks for common frameworks, such as React and Vue.

Finally, I would like to have a `fasthak init` command that automatically generates the minimal scaffolding needed of your static sites, and possibly generates certs too, freeing you from having to use mkcert. 
