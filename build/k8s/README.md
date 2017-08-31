# Deploy elsd in K8s

Setup a clusters for testing

```shell
gcloud beta container clusters create europe-west1-b   --cluster-version 1.7.3   --zone=europe-west1-b   --scopes "cloud-platform,storage-ro,logging-write,monitoring-write,service-control,service-management,https://www.googleapis.com/auth/ndev.clouddns.readwrite"
```

Get credentials

```shell
gcloud container clusters get-credentials europe-west1-b  \
  --zone europe-west1-b
```

Tag and publish the container if you haven done

```shell
docker tag hpccp/elsd:latest gcr.io/print-cloud-software/elsd:latest
```

```shell
gcloud docker -- push gcr.io/print-cloud-software/elsd:latest
```

Deploy

```shell
kubectl apply -f build/k8s/els-deploy.yaml
```