### Steps

- Use your own google cloud account for this and make sure to create a project and a cluster
- Make sure shared resources are applied from root directory `kubectl apply -f shared/`
- `docker build -t randomhash ./pingpong/`
- `docker tag pingpong <gcr-tag>`
- `kubectl apply -f manifests/`
- `kubectl get  svc -n exercises` 
- Access from given ip