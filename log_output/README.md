### Steps

- `k3d cluster create -a 2 --no-lb`
- `kubectl config use-context k3d-k3s-default`
- `docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/kube`
- Make sure shared resources are applied from root directory `kubectl apply -f shared/`
- `docker build -t randomhash ./randomhash/`
- `k3d image import randomhash`
- `docker build -t randomhash-server ./randomhash-server/`
- `k3d image import randomhash-server`
- `kubectl apply -f manifests/`
- `kubectl logs -f <pod-name>`
- Access from `http://localhost:8081/randomhash`
