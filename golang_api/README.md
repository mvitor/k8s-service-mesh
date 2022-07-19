Creation of Golang API hosted on Kubernetes with Docker images

Steps:
1 - Install kind using kind.yaml
2 - Prepare Golang API code to GET data (Data from node will be returned)
3 - Prepare Docker build image to build and host Golang API
4 - Create Kubernetes Deployment with 2 replicas
5 - Create Ingress Service for Clusterip
6 - Access browser

Roadmap to be implemented in next steps:

1 - Implement Github action workflow
2 - Integrate ArgoCD
3 - Implement Autoscaling
4 - Implement NGINX rules
5 - Implement DB with Persistent Volume
