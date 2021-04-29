## Advanced usage of inlets

### Kubernetes operator

Automate tunnels for LoadBalancers and Ingress: [inlets-operator](https://github.com/inlets/inlets-operator)

See also: [Helm chart for inlets PRO](https://github.com/inlets/inlets-pro/tree/master/chart) client and server tunnels

### Docker images

Docker images are published as multi-arch for `x86_64`, `arm64` and `armhf` on the GitHub Container Registry

* [See the latest versions here](https://github.com/orgs/inlets/packages/container/package/inlets)

### Bind to a different adapter, or to localhost

By default the inlets server will bind to all adapters and addresses on your machine.

At times, you may wish to change this, so that you can "hide" the HTTP websocket behind a reverse proxy, adding TLS termination and link-level encryption without exposing the plain HTTP port to the network or Internet.

The diagram below shows how inlets can act as a VPN, when only binding tunnelled services to the local adapter or a private network:

![Tunnelling, but not exposing a service](./inlets-private.png)

> Tunnelling, but not exposing a service.

Usage:

* `--control-addr 127.0.0.1`
* `--control-addr 10.0.101.20`

### Bind a different port for the control-plane

You can bind two separate TCP ports for the user-facing port and the tunnel.

* `--port` - the port for users to connect to and for serving data, i.e. the *Data Plane*
* `--control-port` - the port for the websocket to connect to i.e. the *Control Plane*

### Strict forwarding policy

By default, the server code can access any host. The client specifies a number of upstream hosts via `--upstream`. If you want these to be the only hosts that the server can connect to, then enable strict forwarding.

* `--strict-forwarding`

This is off by default, however when set to true, only hosts in `--upstream` can be accessed by the server. It could prevent a bad actor from accessing other hosts on your network.

### Tunnelling multiple services

You can expose multiple hosts through the `--upstream` flag using a comma-delimited list.

```bash
inlets client --url ws://$IP:8080 \
  --upstream "openfaas.example.com=http://127.0.0.1:8080,prometheus.example.com=http://127.0.0.1:9090"
```

You can also forward everything to a single host such as:

```bash
inlets client --url ws://$IP:8080 \
  --upstream "http://nginx.svc.default"
```
