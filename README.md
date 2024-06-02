# cli

## Features (WIP)

- **Deploy Kubernetes**
- **Deploy Shared Resources**
- **Deploy Services**

## Basic Use

### Preliminaries

1. Setup an S3 bucket backend and Dynamodb Lock Table.
2. Set up a local KinD cluster
   - see `./kind-config.yml`
   - run `kind create cluster --name <whatever> --config kind-config.yml`
3. Grab a Digital Ocean account and token to pass to the cli with `-d`
4. Write some platform configs (see `dev-config.yml` and `prod-config.yml`) to pass to the cli with `-c`.

### Manage Resources of KinD cluster

In this situation, make sure the path to the config points at a file where `provider: kind`.

```
cli k8s <plan|apply|destroy> -b <bucket> -t <table> -c <config>
```

### Manage DO K8s Cluster

In this situation, make sure the path to the config points at a file where `provider: do`.

```
cli infra <plan|apply|destroy> -b <bucket> -t <table> -c <config> -d <do-token>
```

### Manage Resources of DO K8s Cluster

```
cli k8s <plan|apply|destroy> -b <bucket> -t <table> -c <config> -d <do-token>
```

## Write Your Own Terraform

If you don't want to use the terraform found at `github.com/w-h-a`, you can write your own and put them up on a public repository.

Then you can pass in the base source to the cli with `-s`.

### Couple of Assumptions about Repo Naming

The kubernetes terraform must be in a repo with the name:

```
<base url>/kubernetes-<provider>
```

E.g., "github.com/w-h-a/kubernetes-do".

Similarly, the same can be said for:

```
<base url>/kubeconfig
```

and

```
<base url>/kubernetes-namespaces
```
