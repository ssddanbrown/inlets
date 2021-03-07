## Quickstart tutorial

This tutorial shows you how the inlets client and server components work together, using your laptop to run both parts of the tunnel.

* On the *exit-server* (or server)

Start the tunnel server on a machine with a publicly-accessible IPv4 IP address such as a VPS.

Example with a token for client authentication:

```bash
export token=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1)
inlets server \
  --port 8000 \
  --control-port 8001 \
  --token "${token}"
```

A note on security:

> inlets OSS doesn't enable any form of encryption for its control-plane, so you can upgrade to [inlets PRO](https://inlets.dev/) or do this yourself with a separate reverse proxy.

Note down your public IPv4 IP address.

* Head over to your machine where you are running a sample service, or something you want to expose.

Run a HTTP server like [Python's built-in HTTP server](https://docs.python.org/2/library/simplehttpserver.html):

```sh
mkdir -p /tmp/inlets-test/
cd /tmp/inlets-test/
touch hello-world
python -m SimpleHTTPServer 8000
```

* On the same machine, start the inlets client

Start the tunnel client:

```sh
export EXIT_SERVER_IP="192.168.0.100"
export UPSTREAM="http://127.0.0.1:8000"

export CONTROL_PORT="8001"
inlets client \
 --url ws://$EXIT_SERVER_IP:$CONTROL_PORT \
 --upstream=$UPSTREAM \
 --token=$TOKEN
```

Note the warning you receive about using `ws://` instead of `wss://`, then add `--insecure` if you understand what this means and want to continue testing.

* Replace the `--url` with the address where your exit-server is running `inlets server`.
* Replace the `--token` with the value from your server

We now have three processes:

* example service running (hash-browns) or Python's webserver
* an exit-server running the tunnel server (`inlets server`)
* a client running the tunnel client (`inlets client`)

Now connect to the tunnel using the data-port `8000` and the server's IP:

```bash
echo Browse your files at: http://$EXIT_SERVER_IP:8000
```

You can change the data-port from 8000 to 80, however if you wish to serve traffic to clients over the Internet, you will need to configure TLS for the server using a separate reverse proxy.

## See also:

A [personal license for inlets PRO](https://inlets.dev/) can be used at work and at home. It can configures encryption for you using TLS, without the need to install any additional software.

* [Quick-start: expose a local websites with HTTPS with inlets PRO](https://docs.inlets.dev/#/get-started/quickstart-caddy)
* [Tunnel SSH with inlets PRO](https://docs.inlets.dev/#/get-started/quickstart-tcp-ssh)
* [Quick-start: Tunnel a private database over inlets PRO](https://docs.inlets.dev/#/get-started/quickstart-tcp-database)
* [Use inlets PRO with Kubernetes](docs/kubernetes.md)
