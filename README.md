## inlets is a Cloud Native Tunnel written in Go

<img src="docs/inlets-logo-sm.png" width="150px">

Expose your local endpoints to the Internet or within a remote network, without touching firewalls.

[![Build Status](https://travis-ci.com/inlets/inlets.svg?branch=master)](https://travis-ci.com/inlets/inlets)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/inlets/inlets)](https://goreportcard.com/report/github.com/inlets/inlets)
[![Documentation](https://godoc.org/github.com/inlets/inlets?status.svg)](http://godoc.org/github.com/inlets/inlets)
![GitHub All Releases](https://img.shields.io/github/downloads/inlets/inlets/total)

[Follow @inletsdev on Twitter](https://twitter.com/inletsdev)

[English](./README.md) | [中文文档](./README_CN.md)

## Intro

inlets &reg; is how you connect services between different networks. You won't have to think about managing firewalls, NAT or VPNs again. Services can be tunnelled securely over a websocket and accessed on a remote network privately, or exposed on the Internet using an exit-server (5-10USD / mo).

Why do we need this project? Similar tools such as [ngrok](https://ngrok.com/) and [Argo Tunnel](https://developers.cloudflare.com/argo-tunnel/) from [Cloudflare](https://www.cloudflare.com/) are closed-source, have limits built-in, can work out expensive, and have limited support for arm/arm64, Docker and Kubernetes. Ngrok's domain is also often banned by corporate firewall policies meaning it can be unusable. Other open-source tunnel tools are designed to only set up a single static tunnel.

With inlets you can set up your own self-hosted tunnel, copy over the static binary and start tunnelling traffic without any arbitrary limits or restrictions. When used with TLS, inlets can be used with most corporate HTTP proxies.

![Conceptual diagram](docs/inlets.png)

*Conceptual diagram for inlets*

## Do you use inlets? Sponsor the author

Alex is the maintainer of inlets, if you use the project in some way, support it by becoming a sponsor on GitHub.

<a href="https://github.com/sponsors/inlets/">
<img alt="Sponsor this project" src="https://github.com/alexellis/alexellis/blob/master/sponsor-today.png" width="90%">
</a>

You can become a sponsor as an individual or as a corporation, check out the tiers to find which is right for you.

[Find out more now](https://github.com/sponsors/inlets/)

## About inlets OSS

inlets OSS uses a websocket to create a tunnel between a client and a server. The server is typically a machine with a public IP address, and the client is on a private network with no public address.

inlets OSS is for hobbyists only and has no guarantee or warranty of any kind.

For a commercially-supported solution, see [inlets PRO](https://inlets.dev/), which has secure defaults and thorough testing.

### Private or public tunnels?

* A private tunnel is where you start a tunnel to a server and only expose it on the server's LAN address (this can replace the use-cases where you would use a VPN or Kubernetes federation)
* A public tunnel is where you expose the private service to users via the server's public IP

### Features

* Tunnel HTTP or websockets
* Expose multiple sites on same port through the use of DNS entries and a `Host` header
* Upgrade to link encryption using TLS for websockets (`wss://`) with an external add-on or [inlets PRO](https://inlets.dev)
* Shared authentication token for the client and server
* Automatic reconnects for when the connection drops

Distribution:

* Binaries and Docker images for multiple architecture - Intel and ARM
* Dockerfile
* systemd unit file for client/server

inlets PRO has a native Kubernetes LoadBalancer integration through the [inlets-operator](https://github.com/inlets/inlets-operator).

### Going to production with inlets PRO

The following features / use-cases are covered by [inlets PRO](https://inlets.dev):

* Tunnel L4 TCP traffic such as websockets, databases, reverse proxies, remote desktop, VNC and SSH
* Tunnel L7 HTTPS / REST traffic - with automated Let's Encrypt support 
* Expose multiple ports from the same client - i.e. 80 and 443
* Kubernetes helm charts, YAML and [Operator](https://github.com/inlets/inlets-operator)
* Automated TLS for the control-plane
* Commercial services & support
* Documentation, blog posts, tutorials and videos

### inlets projects

inlets is a Cloud Native Tunnel and is [listed on the Cloud Native Landscape](https://landscape.cncf.io/category=service-proxy&format=card-mode&grouping=category&sort=stars) under *Service Proxies*.

* [inlets PRO](https://inlets.dev) - Secure HTTP(s) and TCP tunnels with automated TLS encryption. Replaces inlets OSS
* [inlets-operator](https://github.com/inlets/inlets-operator) - Public IPs for your private Kubernetes Services and CRD using inlets PRO
* [inletsctl](https://github.com/inlets/inletsctl) - The fastest way to create self-hosted exit-servers using inlets PRO
* [inlets](https://github.com/inlets/inlets) - Cloud Native Tunnel for HTTP only - **no** tutorials, automation, TLS, TCP or Kubernetes integration available. Superseded by inlets PRO

## Get inlets OSS

You can install the CLI with a `curl` utility script, `brew` or by downloading the binary from the releases page. Once installed you'll get the `inlets` command.

### Install the CLI

> Note: `inlets` is made available free-of-charge, but you can support its ongoing development and sign up for updates through [GitHub Sponsors](https://github.com/sponsors/alexellis/) 💪

Utility script with `curl`:

```bash
# Install to local directory
curl -sLS https://get.inlets.dev | sh

# Install to /usr/local/bin/
curl -sLS https://get.inlets.dev | sudo sh
```

> Note: the `brew` distribution is maintained by the brew team, so it may lag a little behind the GitHub release.

Binaries are made available on the [releases page](https://github.com/inlets/inlets/releases) for Linux (x86_64, armhf & arm64), Windows (experimental), and for Darwin (MacOS). You will also find SHA checksums available if you want to verify your download.

Windows users are encouraged to use [git bash](https://git-scm.com/downloads) to install the inlets binary.

## inlets PRO

inlets PRO superceds the original inlets OSS project and includes: commercial support, automation, secure defaults, and HTTPS/TCP tunnels.

In "A tale of two networks" Alex Ellis and Johan Siebens explore blog posts, use-cases and show demos of inlets PRO.

[A tale of two networks - demos and use-cases for inlets tunnels](https://www.youtube.com/watch?v=AFMA1xA4zts&feature=youtu.be)

[![https://img.youtube.com/vi/AFMA1xA4zts/hqdefault.jpg](https://img.youtube.com/vi/AFMA1xA4zts/maxresdefault.jpg)](https://youtu.be/AFMA1xA4zts)

### Quickstart tutorial

You can run inlets between any two computers with connectivity, these could be containers, VMs, bare metal or even "loop-back" on your own laptop.

Try the [quickstart tutorial now](./docs/quickstart.md) on your local computer.

### Documentation & tutorials

inlets and inlets PRO have their own documentation site:

Official docs: [docs.inlets.dev](https://docs.inlets.dev)

* Docs: [inlets PRO Kubernetes charts, operator and manifests]([./docs/kubernetes.md](https://github.com/inlets/inlets-pro/tree/master/chart))
* Tutorial: [Get a LoadBalancer for your private Kubernetes cluster with inlets-operator](https://blog.alexellis.io/ingress-for-your-local-kubernetes-cluster/)
* Docs: [Quickstart tutorial on your laptop](./docs/quickstart.md)

Looking for a Docker image? Try [inlets PRO](https://github.com/inlets/inlets-pro)

### What are people saying about inlets?

Read community tutorials, the launch posts on Hacker News, and send a PR if you have written about inlets or inlets PRO:

* [Community and tutorials page](docs/community.md)

> You can share about inlets using `@inletsdev`, `#inletsdev`, and `https://inlets.dev`.

### Do you use inlets or inlets PRO?

Add an entry to the [ADOPTERS.md](./ADOPTERS.md) file with your use-case.

See [inlets PRO testimonials](https://inlets.dev/)

## Disclaimer

Developers wishing to use inlets within a corporate network are advised to seek approval from their administrators or management before using the tool. By downloading, using, or distributing inlets, you agree to the [LICENSE](./LICENSE) terms & conditions. No warranty or liability is provided.

### Development

See [CONTRIBUTING.md](./CONTRIBUTING.md)

#### Other Kubernetes port-forwarding tooling

* [`kubectl port-forward`](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster/) - built into the Kubernetes CLI, forwards a single port to the local computer.
* [inlets PRO](https://inlets.dev/) - exit-server automation, multiple ports, TCP and automatic Let's Encrypt
* [kubefwd](https://github.com/txn2/kubefwd) - Kubernetes utility to port-forward multiple services to your local computer.
* [kurun](https://github.com/banzaicloud/kurun) - Run main.go in Kubernetes with one command, also port-forward your app into Kubernetes.

inlets&reg; is a registered trademark of OpenFaaS Ltd. All rights reserved, registered company in the UK: 11076587
