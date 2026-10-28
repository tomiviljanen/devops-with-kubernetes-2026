### Steps

- `k3d cluster create -a 2 --no-lb`
- `kubectl config use-context k3d-k3s-default`
- `docker build -t the_project ./`
- `k3d image import the_project`
- `kubectl apply -f manifests/deployment.yaml`
- `kubectl logs -f <pod-name>`
