### Steps

- Use your own google cloud account for this and make sure to create a project and a cluster and a repository for images
- `docker build -t randomhash ./pingpong/`
- `docker tag pingpong <gcr-tag>`
- `docker push <gcr-tag>`
- `kubectl apply -f manifests/`
- Make sure shared resources are applied from root directory `kubectl apply -f shared/`
- `kubectl describe -n exercises ingress shared-ingress`
- Access from given ip