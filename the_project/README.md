### Steps

- `k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2`
- `kubectl config use-context k3d-k3s-default`
- `docker build -t the_project ./app`
- `k3d image import the_project`
- `docker build -t backend ./backend`
- `k3d image import backend`
- `kubectl apply -f manifests/`
- Make sure to apply shared ingress and volumes from root folder `kubectl apply -f shared/`
- Access from `http://localhost:8081/`