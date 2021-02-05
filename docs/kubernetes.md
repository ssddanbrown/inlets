# Inlets & Kubernetes recipes

## 1) Run as a deployment on Kubernetes

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
      app: inlets
  template:
    metadata:
      labels:
        app: inlets
    spec:
      containers:
      - name: inlets
        image: inlets/inlets:2.7.4
        imagePullPolicy: Always
        command: ["inlets"]
        args:
        - "client"
        - "--upstream=http://gateway.openfaas:8080,http://prometheus.openfaas:9090"
        - "--url=your-public-ip"
```

Replace the line: `- "--url=your-public-ip"` with the public IP belonging to your VPS.

Alternatively, see the unofficial helm chart from the community: [inlets-helm](https://github.com/paurosello/inlets_helm).

## 2) Use authentication from a Kubernetes secret

In production, you should always use a secret to protect your exit-server. You will need a way of passing that to your server and inlets allows you to read a Kubernetes secret.

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
      app: inlets
  template:
    metadata:
      labels:
        app: inlets
    spec:
      containers:
      - name: inlets
        image: inlets/inlets:2.7.4
        imagePullPolicy: Always
        command: ["inlets"]
        args:
        - "client"
        - "--url=ws://REMOTE-IP"
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

## Use your Kubernetes cluster as an exit-server

You can use a Kubernetes cluster which has public IP addresses, an IngressController, or a LoadBalancer to run one or more exit-servers.

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
      app: inlets
  template:
    metadata:
      labels:
        app: inlets
    spec:
      containers:
      - name: inlets
        image: inlets/inlets:2.7.4
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

## 3) Try inlets with KinD (Kubernetes in Docker)

Try this guide to expose services running in a [KinD cluster](https://github.com/kubernetes-sigs/kind).

You'll expose Kubernetes ClusterIP services with [inlets.dev](https://inlets.dev) to your local computer.

Note: this can now be simplified with `inletsctl kfwd` - [go and check it out](https://github.com/inlets/inletsctl/pull/11)

### Get KinD:

* Download a binary release:

```sh
# Linux

sudo curl -Lo /usr/local/bin/kind \
 https://github.com/kubernetes-sigs/kind/releases/download/v0.6.0/kind-linux-amd64

# MacOS

sudo curl -Lo /usr/local/bin/kind \
 https://github.com/kubernetes-sigs/kind/releases/download/v0.6.0/kind-darwin-amd64
```

* Create the cluster

```sh
kind create cluster
```

* Now switch to the `kind` cluster with `kubectl`:

```sh
export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"
```

### Create a sample service

We'll deploy a HTTP server that runs the `figlet` binary to generate ASCII logos.

* Define a service by running the below:

```yaml
kubectl apply -f - <<EOF
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
EOF
```

* Define a Deployment by running the following:

```yaml
kubectl apply -f - <<EOF
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
EOF
```

Save both files and create the objects with `kubectl create -f filename.yaml`

### Now get inlets and inletsctl

If you don't like running `sh` as `sudo`, then remove `sudo` and simply move the resulting binaries into your `/usr/local/bin/` folder manually after the download.

```sh
curl -sLSf https://get.inlets.dev | sudo sh
curl -sLSf https://raw.githubusercontent.com/inlets/inletsctl/master/get.sh | sudo sh
```

### Start an exit-server on your laptop (inlets server)

Find the IP of your Ethernet or WiFi connection via `ifconfig`, mine was `192.168.0.14`

```sh
inletsctl kfwd --if 192.168.0.14 --from openfaas-figlet:8080
```

That's it, you can now access the figlet microservice via `http://192.168.0.14:8080` or `http://127.0.0.1:8080`

```
curl -d "inlets" 192.168.0.14:8080
 _       _      _       
(_)_ __ | | ___| |_ ___ 
| | '_ \| |/ _ \ __/ __|
| | | | | |  __/ |_\__ \
|_|_| |_|_|\___|\__|___/
                        
```

If you need more flexibility or if want to forward more than one service you can deploy an `inlets client` Deployment manually or edit the one created by `kfwd`.

* Now try another service

Run Nginx and expose it:

```
kubectl run static-web --image nginx --port 80
kubectl expose deploy/static-web --port 8080 --target-port 80
```

```sh
inletsctl kfwd --if 192.168.0.14 --from static-web:8080
```

View the site with the URL printed or at `http://127.0.0.1:8080`.
