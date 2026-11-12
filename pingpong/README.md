### Steps

- `k3d cluster create -a 2 --no-lb`
- `kubectl config use-context k3d-k3s-default`
- `docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/kube`
- `docker build -t randomhash ./pingpong/`
- `k3d image import pingpong`
- `kubectl apply -f manifests/`
- Make sure shared resources are applied from root directory `kubectl apply -f shared/`
- `kubectl logs -f <pod-name>`
- Access from `http://localhost:8081/pingpong`
