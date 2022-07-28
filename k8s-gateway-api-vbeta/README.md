
#  Useful commands

docker build . -t hi-hostname-golang-api

kind create cluster --config kind.yaml

kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v0.5.0/standard-install.yaml


kubectl apply -f gateway-api/gateway-api.yaml
kubectl apply -f hi-hostname-api/hi-hostname-api.yaml
kubectl apply -f hello-hostname-api/hello-hostname-api.yaml

helm repo add kong https://charts.konghq.com
helm repo update
helm install --create-namespace --namespace kong kong kong/kong --set feature-gates=Gateway=true

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

Links:

What is the Gateway API?
https://gateway-api.sigs.k8s.io


Ingress:
https://kubernetes.io/docs/concepts/services-networking/ingress/

API Overview:
https://gateway-api.sigs.k8s.io/concepts/api-overview/#gatewayclass


Google podcast
https://kubernetespodcast.com/episode/186-gateway-api-beta/

Understanding the new Kubernetes Gateway API vs Ingress
https://www.youtube.com/watch?v=Zqlwn5TZknI&t=458s

Kubernetes Networking 101 - Randy Abernethy, RX-M LLC
https://www.youtube.com/watch?v=cUGXu2tiZMc

Github 
https://github.com/kubernetes-sigs/gateway-api

This post shows how to use Shared Access Signature Authentication in Ansible using the native REST API, but the concept utilized here can be applied to any language and/or platform. The same SAS procedure/script can be used for any Azure Storage API integration like Tables and Queues.

<!--more-->
