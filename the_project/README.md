### Steps

- `k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2`
- `kubectl config use-context k3d-k3s-default`
- `docker build -t the_project ./`
- `k3d image import the_project`
- `kubectl apply -f manifests/deployment.yaml`
- `kubectl apply -f manifests/service.yaml`
- Access from `http://localhost:8082`