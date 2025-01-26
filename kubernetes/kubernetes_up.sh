docker build -t fastfood-app:latest .
kubectl apply -f kubernetes/postgres-pv.yaml
kubectl apply -f kubernetes/postgres-pvc.yaml
kubectl apply -f kubernetes/postgres-deploy.yaml
kubectl apply -f kubernetes/postgres-service.yaml
kubectl apply -f kubernetes/fastfood-deployment.yaml
kubectl apply -f kubernetes/fastfood-service.yaml
kubectl apply -f kubernetes/fastfood-hpa.yaml