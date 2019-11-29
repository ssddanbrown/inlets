# Inlets & Kubernetes recipes

## Run as a deployment on Kubernetes

You can run the client inside Kubernetes to expose your local services to the Internet, or another network.

Here's an example showing how to get ingress into your cluster for your OpenFaaS gateway and for Prometheus:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inlets
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: inlets
  template:
    metadata:
      labels:
        app.kubernetes.io/name: inlets
    spec:
      containers:
      - name: inlets
        image: inlets/inlets:2.6.1
        imagePullPolicy: Always
        command: ["inlets"]
        args:
        - "client"
        - "--upstream=http://gateway.openfaas:8080,http://prometheus.openfaas:9090"
        - "--remote=your-public-ip"
```

Replace the line: `- "--remote=your-public-ip"` with the public IP belonging to your VPS.

Alternatively, see the unofficial helm chart from the community: [inlets-helm](https://github.com/paurosello/inlets_helm).

## Use authentication from a Kubernetes secret

In production, you should always use a secret to protect your exit-node. You will need a way of passing that to your server and inlets allows you to read a Kubernetes secret.

* Create a random secret

```
$ kubectl create secret generic inlets-token --from-literal token=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1)
secret/inlets-token created
```

* Or create a secret with the value from your remote server

```
$ export TOKEN=""
$ kubectl create secret generic inlets-token --from-literal token=${TOKEN}
secret/inlets-token created
```

* Bind the secret named `inlets-token` to the Deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inlets
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: inlets
  template:
    metadata:
      labels:
        app.kubernetes.io/name: inlets
    spec:
      containers:
      - name: inlets
        image: inlets/inlets:2.6.1
        imagePullPolicy: Always
        command: ["inlets"]
        args:
        - "client"
        - "--remote=ws://REMOTE-IP"
        - "--upstream=http://gateway.openfaas:8080"
        - "--token-from=/var/inlets/token"
        volumeMounts:
          - name: inlets-token-volume
            mountPath: /var/inlets/
      volumes:
        - name: inlets-token-volume
          secret:
            secretName: inlets-token
```

Optional tear-down:

```
$ kubectl delete deploy/inlets
$ kubectl delete secret/inlets-token
```

## Use your Kubernetes cluster as an exit-node

You can use a Kubernetes cluster which has public IP addresses, an IngressController, or a LoadBalancer to run one or more exit-nodes.

* Create a random secret

```
$ kubectl create secret generic inlets-token --from-literal token=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1)
secret/inlets-token created
```

* Or create a secret with the value from your remote server

```
$ export TOKEN=""
$ kubectl create secret generic inlets-token --from-literal token=${TOKEN}
secret/inlets-token created
```

* Create a `Service`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: inlets
  labels:
    app: inlets
spec:
  type: ClusterIP
  ports:
    - port: 8000
      protocol: TCP
      targetPort: 8000
  selector:
    app: inlets
```

* Create a `Deployment`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inlets
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: inlets
  template:
    metadata:
      labels:
        app.kubernetes.io/name: inlets
    spec:
      containers:
      - name: inlets
        image: inlets/inlets:2.6.1
        imagePullPolicy: Always
        command: ["inlets"]
        args:
        - "server"
        - "--token-from=/var/inlets/token"
        volumeMounts:
          - name: inlets-token-volume
            mountPath: /var/inlets/
      volumes:
        - name: inlets-token-volume
          secret:
            secretName: inlets-token
```

You can now create an `Ingress` record, or `LoadBalancer` to connect to your server.
Note that clients connecting to this server will have to specify port 8000 for their remote, as the default is 80.

## Try inlets with KinD (Kubernetes in Docker)

Try this guide to expose services running in a [KinD cluster](https://github.com/kubernetes-sigs/kind).

You'll expose Kubernetes ClusterIP services with [inlets.dev](https://inlets.dev) to your local computer.

Note: this can now be simplified with `inletsctl kfwd` - [go and check it out](https://github.com/inlets/inletsctl/pull/11)

### Get KinD:

```sh
# Linux

sudo curl -Lo /usr/local/bin/kind \
 https://github.com/kubernetes-sigs/kind/releases/download/v0.4.0/kind-linux-amd64

# MacOS

sudo curl -Lo /usr/local/bin/kind \
 https://github.com/kubernetes-sigs/kind/releases/download/v0.4.0/kind-darwin-amd64
```

Create the cluster

```
kind create cluster
```

### Switch to the `kind` cluster with `kubectl`

```sh
export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"
```

### Create a sample service

We'll deploy a HTTP server that runs the `figlet` binary to generate ASCII logos

* Define a service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: openfaas-figlet
  labels:
    app: openfaas-figlet
spec:
  type: ClusterIP
  ports:
    - port: 8080
      protocol: TCP
      targetPort: 8080
  selector:
    app: openfaas-figlet
```

Define a Deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: openfaas-figlet
  labels:
   app: openfaas-figlet
spec:
  replicas: 1
  selector:
    matchLabels:
      app: openfaas-figlet
  template:
    metadata:
      labels:
        app: openfaas-figlet
    spec:
      containers:
      - name: openfaas-figlet
        image: functions/figlet:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
          protocol: TCP
```

Save both files and create the objects with `kubectl create -f`

### Now get inlets

```sh
cd /tmp/

# Download to local directory
curl -sLS https://get.inlets.dev | sh

chmod +x ./inlets
sudo mv inlets /usr/local/bin/

inlets --version
Version: 2.1.0
Git Commit: c23f6993892a1b4e398e8acf61e3dc7bfcb7c6ed
```

### Start an exit-node on your laptop (inlets server)

Our Kubernetes cluster will connect to this server.

```
export token=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1)
inlets server --port=8090 --token="$token" --print-token=true
```

Note your `token` when the server starts up.

### Run the inlets client as a Kubernetes Deployment

Create a secret for the inlets client:

```sh
export TOKEN="" # Use the value from earlier
kubectl create secret generic inlets-token --from-literal token=${TOKEN}
```

Apply the Deployment YAML file, with `kubectl apply -f`.

Change the following two parameters:

Use your laptop's IP in place of `REMOTE-IP`:

```sh
- "--remote=ws://REMOTE-IP"
```

My IP for my WiFi interface is `192.168.1.51`.

> Note: your "exit-node" could be any PC that has reachability including a VPS with a public IPv4 address.

Specify the service, or services which you want to expose:

```
- "--upstream=http://openfaas-figlet.default:8080"
```

This is the sample YAML:

```yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: inlets
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: inlets
    spec:
      containers:
      - name: inlets
        image: alexellis2/inlets:2.1.0
        imagePullPolicy: Always
        command: ["inlets"]
        args:
        - "client"
        - "--remote=ws://REMOTE-IP"
        - "--upstream=http://openfaas-figlet:8080"
        - "--token-from=/var/inlets/token"
        volumeMounts:
          - name: inlets-token-volume
            mountPath: /var/inlets/
      volumes:
        - name: inlets-token-volume
          secret:
            secretName: inlets-token
```

### Access your service

You can now access the service inside the KinD cluster, from the inlets server port and IP.

```sh
curl 192.168.1.51:8090 -d "inlets.dev"
 _       _      _            _            
(_)_ __ | | ___| |_ ___   __| | _____   __
| | '_ \| |/ _ \ __/ __| / _` |/ _ \ \ / /
| | | | | |  __/ |_\__ \| (_| |  __/\ V / 
|_|_| |_|_|\___|\__|___(_)__,_|\___| \_/  
```

You could also use `127.0.0.1:8090` on your local machine.

### Access multiple services

Run Nginx and expose it:

```
kubectl run static-web --image nginx --port 80
kubectl expose deploy/static-web --port 80 --target-port 80
```

Edit the upstream parameter (`kubectl edit deploy/inlets`):

```
        - "--upstream=openfaas-figlet.local=http://openfaas-figlet:8080,static-web.local=http://static-web:80"
```

Now setup two hosts file entries in `/etc/hosts`:

```
127.0.0.1  openfaas-figlet.local
127.0.0.1  static-web.local
```

Now access either:

```
curl -d hi http://127.0.0.1:8090 -H "Host: openfaas-figlet.local"
curl http://127.0.0.1:8090 -H "Host: static-web.local"
```