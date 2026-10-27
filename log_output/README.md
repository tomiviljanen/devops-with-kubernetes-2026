### Steps

- `k3d cluster create -a 2 --no-lb`
- `kubectl config use-context k3d-k3s-default`
- `docker build -t randomhash ./randomhash/`
- `k3d image import randomhash`
- `kubectl create deployment randomhash --image=randomhash`
- `kubectl edit deployment randomhash`
  - edit `imagePullPolicy` to `Never`
- `kubectl logs -f <pod-name>`
