docker build -t video-processing-app:latest .

kubectl create secret generic aws-secrets \
  --from-literal=access-key-id=$(aws configure get aws_access_key_id) \
  --from-literal=secret-access-key=$(aws configure get aws_secret_access_key) \
  --from-literal=access-session-token=$(aws configure get aws_session_token)
kubectl apply -f kubernetes/video-processing-deployment.yaml
kubectl apply -f kubernetes/video-processing-service.yaml
kubectl apply -f kubernetes/video-processing-hpa.yaml