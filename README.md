kubetrim
====================================================

Tidy up old Kubernetes clusters from kubeconfig.

kubetrim tries to connect to each cluster in the current kubeconfig file, and removes any that are unreachable, or which error.

Q&A:

* What if I like doing things the long way?

    You can combine `kubectl config get-clusters` with `kubectl config use-context` and `kubectl get nodes`, followed by `kubectl config delete-cluster` and `kubectl config delete-context` for each.

* Doesn't [my favourite tool] have an option to do this?

    Feel free to use your favourite tool instead. kubetrim is a memorable name, and a simple tool that just does one job, similar to `kubectx`

* What if I want to keep a cluster that is unreachable?

    This is not supported at this time, if you need that feature, open an issue.

* What if my cluster is valid, but kubetrim cannot detect it?

    Open an issue, and we can look at adding support for your use-case.

* `kubetrim` is great, how can I support you?

    Have a look at [arkade](https://github.com/alexellis/arkade) and [k3sup](https://github.com/alexellis/k3sup), you may like those too. You can also [sponsor me via GitHub](https://github.com/sponsors/alexellis).

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
