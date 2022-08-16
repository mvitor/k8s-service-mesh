# Kubernetes Gateway API tutorial

## Kubernetes Gateway API 

Gateway API is an open source project managed by the Kubernetes Network Special Interest Group (SIG‑NETWORK) community to improve and standardize service networking in Kubernetes.

### Gateway API Beta Announcement 

Recently Kubernetes SIG-Network team has announced the v0.5.0 release of Gateway API, Kubernetes Gateway API is graduating to Beta.

### Why Gateway API

It's a decision you need to take when designing your applications. To discuss this decision is not the goal of this article then I'm sharing the following NGINX Blog post and Webinar with information about this design decision.

#### How Do I Choose? API Gateway vs. Ingress Controller vs. Service Mesh

This NGINX Blog post guides you through the decision about which technology to use for API gateway use cases, with sample scenarios for north‑south and east‑west API traffic.

[https://www.nginx.com/blog/how-do-i-choose-api-gateway-vs-ingress-controller-vs-service-mesh/](https://www.nginx.com/blog/how-do-i-choose-api-gateway-vs-ingress-controller-vs-service-mesh/)


#### API Gateway Use Cases for Kubernetes

NGINX Webinar Discussing the various tools and use cases, our experts demo how you can use an Ingress controller and service mesh to accomplish API gateway use cases.

[https://www.nginx.com/resources/webinars/api-gateway-use-cases-for-kubernetes/](https://www.nginx.com/resources/webinars/api-gateway-use-cases-for-kubernetes/)

## Tutorial Main Sections

This tutorial is separated in <three main sections which are requirements for the next steps.   

#### 1 - Build Golang APIs to route traffic 

We need to have Kubernetes HTTP Services so we can route traffic to them. In this section we're building two Golang APis and deploying them as Kubernetes Services. 

* If you already have APIs to route traffic in your K8s cluster this step can be skipped.  

#### 2 - Kubernetes Gateway API tutorial with NGINX controller


NGINX Kubernetes Gateway is an open-source project that provides an implementation of the Gateway API using NGINX. That project goal is to implement the core Kubernetes Gateway APIs functionalities which are being released by Kubernetes SIG Network team: Gateway, GatewayClass, HTTPRoute, Croute, TLSRoute, and UDPRoute which allow to configure an HTTP or TCP/UDP load balancer, reverse-proxy, or API gateway for applications running on Kubernetes.

The steps described on this section are taken from the official [nginx-kubernetes-gateway](https://github.com/nginxinc/nginx-kubernetes-gateway/) repository. It's going to create a Nginx Gateway API image and make it available to your cluster, install the controllers, gateway classes and finally setup Nginx proxy.

* If you already have an nginx-kubernetes-gateway image running and Gateway classes available this step can be skipped. 

 
#### 3 - Installing NGINX Kubernetes Gateway API Class and HTTP Routes

NGINX is an active contributor to the Kubernetes Gateway API project and is up to date with the most recent features released. We will be utilizing NGINX Kubernetes Gateway controller that implements the Kubernetes Gateway API specification. 

## Tutorial configuration files and scripts repository 

All YAML configuration file and scripts utilized on this tutorial are available on this Github repository folder [k8s-gateway-api-vbeta](https://github.com/mvitor/k8s-service-mesh/tree/main/k8s-gateway-api-vbeta).


## 1 -  Build Golang APIs to route traffic 

This step is when we're actually leveraging the Gateway API functionalities. We're creating the Gateway API and the HTTP Routes in different ways. 


###  Create Hi and Hello Golang APIs

We're using two different Golang APIs to route the traffic. I'm creating two hypothetical API. The first one should greet with a Hi and the name of the Pod, and the second one should greet with a Hello and the name of the Pod. With this we're able to differentiate between the APIs calls we will use to validate the routing works as expected.

### Create Kind Cluster

We will be using Kind to run a local cluster, the procedure works for any [A-Z]KS cluster. For EKS it's suggested using different Load Balancer configuration. 

#### Kind Manifest file
```
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
```
#### Kind Create command

```sh
kind create cluster --config kind.yaml
```

### Golang APIs

We will use Golang to return the Kubernetes pod name and greet which eventually can be 'Hi' and 'Hello'. I'm sharing below the code for both APIs which are similar. We will need to build also Hi and Hello APis Docker images and push it to leverage those images inside the Kubernetes cluster.

#### Golang Hi API
```sh
package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", HandleGet)
	http.ListenAndServe(":8080", nil)
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Fprintf(w, "GoLang Hi from Pod: %s", hostname)
}

```
### Golang Hello API 
```sh
package main
import (
	"fmt"
	"net/http"
	"os"
)
func main() {
	http.HandleFunc("/", HandleGet)
	http.ListenAndServe(":8080", nil)
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Fprintf(w, "GoLang Hello from Pod: %s", hostname)
}
```
### Golang APIs Docker images creation

As we will be using Docker as container runtime, we need to prepare our images, which means we need to build and publish our image in a public or in a private repository, in our case we will be pushing it to Dockerhub. I'm sharing below Docker build file and the commands use to push the image to the Dockerhub repository. 

#### Golang Docker file

Here, we're using [multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/) to have the smallest image size possible. First, we build the Golang application adding all dependencies then we deploy this application inside an Alpine runtime image which is going to host and run the Golang binary application. The Dockerfiles are the same for both APIs.

```sh
FROM golang:1.19rc2-alpine3.16 as dev

WORKDIR /work

##
## Build
##

FROM golang:1.19rc2-alpine3.16 as build

WORKDIR /api
COPY ./api/* /api/
RUN go build -o api

##
## Deploy build to image
##

FROM alpine as runtime 
COPY --from=build /api/api /
CMD ./api

EXPOSE 8080
```
#### Docker images creation and publishing to Dockerhub
##### 'Hi' Golang-api
```sh
> docker build ./hi-hostname-api . -t hi-hostname-golang-api
hi-hostname-api> docker build tag hi-hostname-golang-api mvitor/hi-hostname-golang-api
> docker build tag hi-hostname-golang-api mvitor/hi-hostname-golang-api
> docker push mvitor/hi-hostname-golang-api
```
##### 'Hello' Golang-api

```sh
> docker build ./hello-hostname-api -t hello-hostname-api
> docker build tag hello-hostname-api mvitor/hello-hostname-api
> docker push mvitor/hello-hostname-api
```
Your docker client must be logged in to Dockerhub in order to push the images, if you need help configuring this, this [article](https://www.techrepublic.com/article/how-to-successfully-log-in-to-dockerhub-from-the-command-line-interface/) might be helpful.


### Creating Hi and Hello Kubernetes Deployment 

Now we have our images deployed in Dockerhub we can deploy it in our Kubernetes cluster creating a Deployment object like the below:

#### Kubernetes Deployment - Hello API
```sh
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-hostname-api
  labels:
    app: hello-hostname-api
spec:
  selector:
    matchLabels:
      app: hello-hostname-api
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: hello-hostname-api
    spec:
      containers:
      - name: hello-hostname-api
        image: mvitor/hostname-golang-api:latest
        imagePullPolicy : Always
        ports:
        - containerPort: 8080
        env:
        - name: "ENVIRONMENT"
          value: "DEV"
---
apiVersion: v1
kind: Service
metadata:
  name: hello-hostname-api
spec:
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
    name: http
  selector:
    app: hello-hostname-api
```
#### Kubernetes Deployment - Hi API
```sh
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hi-hostname-api
  labels:
    app: hi-hostname-api
spec:
  selector:
    matchLabels:
      app: hi-hostname-api
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: hi-hostname-api
    spec:
      containers:
      - name: hi-hostname-api
        image: mvitor/hi-hostname-golang-api:latest
        imagePullPolicy : Always
        ports:
        - containerPort: 8080
        env:
        - name: "ENVIRONMENT"
          value: "DEV"
---
apiVersion: v1
kind: Service
metadata:
  name: hi-hostname-api
spec:
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
    name: http
  selector:
    app: hi-hostname-api
```

### Creating  Deployment Exposing Two Http services

```sh
kubectl apply -f hi-hostname-api/hi-hostname-api.yaml 
kubectl apply -f hello-hostname-api/hello-hostname-api.yaml
```
### Checking Deployment Cluster Status

Screenshot Cluster status

## 2 - NGINX Kubernetes Gateway Setup

### Build Nginx Gateway api Image 

#### Clone nginx-kubernetes-gateway 
```
git clone https://github.com/nginxinc/nginx-kubernetes-gateway.git
cd nginx-kubernetes-gateway
```
#### Build Image
```
make PREFIX=myregistry.example.com/nginx-kubernetes-gateway container
```
Set the PREFIX variable to the name of the registry you'd like to push the image to. By default, the image will be named nginx-kubernetes-gateway:0.0.1.
#### Push the image 
```
docker push myregistry.example.com/nginx-kubernetes-gateway:0.0.1
```
Make sure to substitute myregistry.example.com/nginx-kubernetes-gateway with your private registry.


### Kind Load Image 
#### Load the NGINX Kubernetes Gateway image onto your kind cluster
```
kind load docker-image mvitor/nginx-k8s-gateway:0.0.1
```
#### Install the Gateway CRDs
```
kubectl apply -k "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v0.4.2"

kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v0.5.0/standard-install.yaml

???

```
#### Create the nginx-gateway Namespace 
```
kubectl apply -f deploy/manifests/namespace.yaml
```
#### Create the njs-modules configmap
```
kubectl create configmap njs-modules --from-file=internal/nginx/modules/src/httpmatches.js -n nginx-gateway
```
#### Create the GatewayClass resource
```
kubectl apply -f deploy/manifests/gatewayclass.yaml
```
#### Deploy the NGINX Kubernetes Gateway:
```
kubectl apply -f deploy/manifests/nginx-gateway.yaml
```
#### Create Load Balancer Service 
```
kubectl apply -f  deploy/manifests/service/loadbalancer.yaml -n nginx-gateway
```
## 3 -  Create Gateway API and HTTP Rourts

### Create Gateway API Class
```
kubectl apply -f gateway-api/gateway-api.yaml
```

### Create HTTP Routes


```
kubectl apply -f gateway-api/gateway-api.yaml
```
## Access it using Port-forward 

```
kubectl port-forward svc/nginx-gateway 8080:80 -n nginx-gateway
```


curl --resolve mvitormais.com:8080:127.0.0.1 http://mvitormais.com:8080/greet


helm repo add kong https://charts.konghq.com
helm repo update
helm install kong kong/kong --values 

helm install --create-namespace --namespace kong kong kong/kong --set feature-gates=Gateway=true

helm install --values consul-values.yaml consul hashicorp/consul

## Metrics Yaml
kubectl apply -f ../metrics_components.yaml


Metrics is pulling data every 60s to save resources. 
https://gateway-api.sigs.k8s.io/guides/getting-started/#installing-a-gateway-controller

# Install Standard Channel¶
The standard release channel includes all resources that have graduated to beta, including GatewayClass, Gateway, and HTTPRoute. To install this channel, run the following kubectl command:

```
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v0.5.0/standard-install.yaml
```
# Deploying a simple Gateway¶

The simplest possible deployment is a Gateway and Route resource which are deployed together by the same owner. This represents a similar kind of model used for Ingress. In this guide, a Gateway and HTTPRoute are deployed which match all HTTP traffic and directs it to a single Service named foo-svc.


https://gateway-api.sigs.k8s.io/images/single-service-gateway.png

```
apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: prod-web
spec:
  gatewayClassName: acme-lb
  listeners:  
  - protocol: HTTP
    port: 80
    name: prod-web-gw
    allowedRoutes:
      namespaces:
        from: Same
```
# Deploying a simple Gateway¶

The Gateway represents the instantation of a logical load balancer. It's templated from a hypothetical acme-lb GatewayClass. The Gateway listens for HTTP traffic on port 80. This particular GatewayClass automatically assigns an IP address which will be shown in the Gateway.status after it has been deployed.

Route resources specify the Gateways they want to attach to using ParentRefs. As long as the Gateway allows this attachment (by default Routes from the same namespace are trusted), this will allow the Route to receive traffic from the parent Gateway. BackendRefs define the backends that traffic will be sent to. More complex bi-directional matching and permissions are possible and explained in other guides.

# Http to match all trafic

The following HTTPRoute defines how traffic from the Gateway listener is routed to backends. Because there are no host routes or paths specified, this HTTPRoute will match all HTTP traffic that arrives at port 80 of the load balancer and send it to the foo-svc Pods.

```
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: foo
spec:
  parentRefs:
  - name: prod-web
  rules:
  - backendRefs:
    - name: foo-svc
      port: 8080

```

## Clean up 

kubectl delete -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v0.5.0


kind delete cluster


### Links

https://github.com/nginxinc/nginx-kubernetes-gateway/


What is the Gateway API?
https://gateway-api.sigs.k8s.io

API Overview:
https://gateway-api.sigs.k8s.io/concepts/api-overview/#gatewayclass


Google podcast
https://kubernetespodcast.com/episode/186-gateway-api-beta/

Understanding the new Kubernetes Gateway API vs Ingress
https://www.youtube.com/watch?v=Zqlwn5TZknI&t=458s

Ingress:
https://kubernetes.io/docs/concepts/services-networking/ingress/

Kubernetes Networking 101 - Randy Abernethy, RX-M LLC
https://www.youtube.com/watch?v=cUGXu2tiZMc

Github 
https://github.com/kubernetes-sigs/gateway-api
[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)

   [dill]: <https://github.com/joemccann/dillinger>
   [git-repo-url]: <https://github.com/joemccann/dillinger.git>
   [john gruber]: <http://daringfireball.net>
   [df1]: <http://daringfireball.net/projects/markdown/>
   [markdown-it]: <https://github.com/markdown-it/markdown-it>
   [Ace Editor]: <http://ace.ajax.org>
   [node.js]: <http://nodejs.org>
   [Twitter Bootstrap]: <http://twitter.github.com/bootstrap/>
   [jQuery]: <http://jquery.com>
   [@tjholowaychuk]: <http://twitter.com/tjholowaychuk>
   [express]: <http://expressjs.com>
   [AngularJS]: <http://angularjs.org>
   [Gulp]: <http://gulpjs.com>

   [PlDb]: <https://github.com/joemccann/dillinger/tree/master/plugins/dropbox/README.md>
   [PlGh]: <https://github.com/joemccann/dillinger/tree/master/plugins/github/README.md>
   [PlGd]: <https://github.com/joemccann/dillinger/tree/master/plugins/googledrive/README.md>
   [PlOd]: <https://github.com/joemccann/dillinger/tree/master/plugins/onedrive/README.md>
   [PlMe]: <https://github.com/joemccann/dillinger/tree/master/plugins/medium/README.md>
   [PlGa]: <https://github.com/RahulHP/dillinger/blob/master/plugins/googleanalytics/README.md>




