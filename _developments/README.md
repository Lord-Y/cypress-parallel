# developments

## Creation
For our development we will use kind like this:

```bash
sudo kind create cluster --name cypress-parallel --image kindest/node:v1.22.0 --wait 5m --config kind-config.yaml
sudo kubectl cluster-info --context kind-cypress-parallel
sudo cp -r /root/.kube/ ~/ && sudo chown -R $USER:$USER ~/.kube

# to fix dns issues
kubectl -n kube-system apply -f configmap.yaml
kubectl -n kube-system rollout restart deploy coredns
```

Kind documentation can be found [here](https://kind.sigs.k8s.io/docs/user/quick-start/)

## Docker images

Load docker image into `Kind`:
```bash
CLI_IMAGE=ghcr.io/lord-y/cypress-parallel-docker-images/cypress-parallel-docker-images:10.10.0-0.3.0
KIND_CLUSTER_NAME=cypress-parallel
sudo docker pull ${CLI_IMAGE}
time for i in cypress-parallel-worker cypress-parallel-worker2; do sudo kind load docker-image ${CLI_IMAGE} --name ${KIND_CLUSTER_NAME} --nodes $i;done
```

## Cleaning

```bash
sudo kind delete cluster --name cypress-parallel
```
