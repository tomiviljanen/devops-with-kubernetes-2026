### Steps

- `k3d cluster create -a 2 --no-lb`
- `kubectl config use-context k3d-k3s-default`
- `docker build -t the_project ./`
- `k3d image import the_project`
- `kubectl create deployment the-project --image=the_project`
- `kubectl edit deployment the-project`
  - edit `imagePullPolicy` to `Never`
- `kubectl logs -f <pod-name>`
