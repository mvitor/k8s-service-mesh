# k8s-service-mesh


q
kubectl apply -f video-web.yml

kubectl port-forward svc/videos-web 8080:80

kubectl apply -f playlist-api/playlist-api.yaml

kubectl apply -f playlist-db/

kubectl port-forward svc/playlists-api 8080:80

kubectl apply -f videos-db/

https://github.com/mvitor/k8s-service-mesh.git