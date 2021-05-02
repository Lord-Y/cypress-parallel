# developments

For our development we will use kind like this:

```bash
sudo kind create cluster --name cypress-parallel --image kindest/node:v1.19.4 --wait 5m --config kind-config.yaml
sudo kubectl cluster-info --context kind-cypress-parallel
sudo cp -r /root/.kube/ ~/ && sudo chown -R $USER:$USER ~/.kube
```

Kind documentation can be found [here](https://kind.sigs.k8s.io/docs/user/quick-start/)