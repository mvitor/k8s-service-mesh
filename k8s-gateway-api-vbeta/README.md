# Gateway API Announcement 

Recently Kubernet SIG Network team has announced the alpha release of Gateway API

# Tutorial Sections

This tutorial is seggregated in three main sections which is requirements for the next steps.   

## Build Golang APIs to route traffic 

We need to have Kubernetes HTTP Services so we can route traffic to them. In this section we're building two Golang APis and deploying them as Kubernetes Services. 

If you already have APIs to route traffic in your K8s cluster this step can be skipped.  

## Install NGINX Kubernetes Gateway

NGINX Kubernetes Gateway is an open-source project that provides an implementation of the Gateway API using NGINX. That project goal is to implement the core Kubernetes Gateway APIs funcionalities which are being released by Kubernetes SIG Network team: Gateway, GatewayClass, HTTPRoute, TCPRoute, TLSRoute, and UDPRoute which allow to configure an HTTP or TCP/UDP load balancer, reverse-proxy, or API gateway for applications running on Kubernetes.

The steps described on this section are taken from the official [nginx-kubernetes-gateway](https://github.com/nginxinc/nginx-kubernetes-gateway/) repository. It's going to create a Nginx Gateway API image and make it available to your cluster, install the controllers, gateway classes and finally setup Nginx proxy.

If you are already have an nginx-kubernetes-gateway image running and Gateway classes available this step can be skipped. 

## Create Gateway API and Http Routes 

This step is where we're actually levereging the Gateway API funcionalities. We're creating the Gateway API and the HTTP Routes in different ways. 
#  Create Hi and Hello Golang APIs

We're using two different Golang APIs to route the traffic. I'm creating two hyphotetical API. FIrst one should Greet with a Hi and the name of the Pod, second one should greet with a Hello and the name of the Pod. With this we're able to differentiate between the APIs calls we will use to validate the routing works as expected.

## Kind

We're using Kind to run a local cluster but the procedure works for any [A-Z]KS cluster. :) For EKS it's suggested different Load Balancer configuration. 

### Kind Manifest file
```
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
```
### Kind Create command

```
kind create cluster --config kind.yaml
```

##
## Exposing Two Http services

kubectl apply -f hi-hostname-api/hi-hostname-api.yaml 
kubectl apply -f hello-hostname-api/hello-hostname-api.yaml

docker build . -t hi-hostname-golang-api

### Golang Hi API

### Golang Hello API 
# NGINX Kubernetes Gateway Setup
## Build Nginx Gateway api Image 

### Clone nginx-kubernetes-gateway 
```
git clone https://github.com/nginxinc/nginx-kubernetes-gateway.git
cd nginx-kubernetes-gateway
```
### Build Image

make PREFIX=myregistry.example.com/nginx-kubernetes-gateway container

Set the PREFIX variable to the name of the registry you'd like to push the image to. By default, the image will be named nginx-kubernetes-gateway:0.0.1.
### Push the image 
docker push myregistry.example.com/nginx-kubernetes-gateway:0.0.1

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

kubectl apply -f deploy/manifests/namespace.yaml

#### Create the njs-modules configmap
kubectl create configmap njs-modules --from-file=internal/nginx/modules/src/httpmatches.js -n nginx-gateway

#### Create the GatewayClass resource

kubectl apply -f deploy/manifests/gatewayclass.yaml

#### Deploy the NGINX Kubernetes Gateway:

kubectl apply -f deploy/manifests/nginx-gateway.yaml

#### Create Load Balancer Service 

kubectl apply -f  deploy/manifests/service/loadbalancer.yaml -n nginx-gateway

# Create Gateway API
## Create Gateway API Class

kubectl apply -f gateway-api/gateway-api.yaml
## Create Gateway API API
### Create HTTP Routes
###   
kubectl apply -f gateway-api/gateway-api.yaml
## Access it using Port-forward 
kubectl port-forward svc/nginx-gateway 8080:80 -n nginx-gateway

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

# Clean up 

kubectl delete -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v0.5.0

A
Links:

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

This post shows how to use Shared Access Signature Authentication in Ansible using the native REST API, but the concept utilized here can be applied to any language and/or platform. The same SAS procedure/script can be used for any Azure Storage API integration like Tables and Queues.

<!--more-->
