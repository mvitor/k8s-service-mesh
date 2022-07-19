# k8s-service-mesh


## ingress-nginx
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.2.1/deploy/static/provider/cloud/deploy.yaml


## Manifests
kubectl apply -f video-web/

kubectl port-forward svc/videos-web 8080:80

kubectl apply -f playlist-api/

kubectl apply -f playlist-db/

kubectl port-forward svc/playlists-api 8080:80

kubectl apply -f videos-db/

https://github.com/mvitor/k8s-service-mesh.git