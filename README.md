kubetrim
====================================================

Tidy up any old or broken Kubernetes clusters and contexts from your kubeconfig.

kubetrim tries to connect to each cluster in the current kubeconfig file, and removes any that are unreachable, or which error.

## Usage

```bash
$ kubectx

default
do-lon1-openfaas-cluster
kind-2
kind-ingress

$ kubetrim

kubetrim (dev) by Alex Ellis 

Loaded: /home/alex/.kube/config. Checking..
  - kind-2: ✅
  - kind-ingress: ❌ - (failed to connect to cluster: Get "https://127.0.0.1:40349/api/v1/nodes": dial tcp 127.0.0.1:40349: connect: connection refused)
  - default: ✅
  - do-lon1-openfaas-cluster: ❌ - (failed to connect to cluster: Get "https://da39a3ee5e6b4b0d3255bfef95601890afd80709.k8s.ondigitalocean.co.uk/api/v1/nodes": dial tcp: lookup da39a3ee5e6b4b0d3255bfef95601890afd80709.k8s.ondigitalocean.co.uk on 127.0.0.53:53: no such host)
Updated: /home/alex/.kube/config (in 364ms).

$ kubectx

default
kind-2
```

Try out kubetrim without writing changes to the kubeconfig file:

```bash
$ kubetrim --dry-run
```

Use a different kubeconfig file:

```bash
$ KUBECONFIG=$HOME/.kube/config.bak kubetrim
```

## Installation

Getting `kubetrim` with arkade:

```bash
curl -sfLS https://get.arkade.dev | sh

arkade get kubetrim
```

## License

MIT license
