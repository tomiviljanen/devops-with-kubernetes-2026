### Steps

- `k3d cluster create -a 2 --no-lb`
- `kubectl config use-context k3d-k3s-default`
- `docker build -t randomhash ./randomhash/`
- `k3d image import randomhash`
- `kubectl apply -f manifests/`
- Make sure shared ingress is applied from root directory `kubectl apply -f shared/`
- `kubectl logs -f <pod-name>`
- Access from `http://localhost:8082`
