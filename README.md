# NFS Minio Operator
A Kubernetes operator to access your NFS data using Minio's S3 compatible API.

In development.

## Deployment
```bash
kubectl create -f deploy/crds/nfsminiooperator_v1alpha1_nfsminio_crd.yaml
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/operator.yaml
```

## Devlopment
```bash
export GO111MODULE=on
```

### Development Depoyment
Run this one time
```bash
kubectl create -f deploy/crds/nfsminiooperator_v1alpha1_nfsminio_crd.yaml
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/operator.yaml

export OPERATOR_NAME=nfs-minio-operator
```

Then for testing outside the cluster (you need connectivity with your Kubernetes cluster using kubectl) the controller run:
```bash
operator-sdk up local --namespace=default
```

