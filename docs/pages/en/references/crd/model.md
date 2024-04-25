# `Model`

`Model` is a custom resource definition (CRD) that represents a Ollama server instance in the cluster.

When created, the creation of `Model` triggers the creation of [`Deployment`](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/), [`Service`](https://kubernetes.io/docs/concepts/services-networking/service/).

> If it is created for the first time, an additional StatefulSet and PersistentVolumeClaim will be created to store and persist the downloaded Ollama image.

## Full CRD Reference

```yaml
apiVersion: ollama.ayaka.io/v1
kind: Model
metadata:
  name: phi
spec:
  # Scale the model to 2 replicas
  replicas: 2
  # Use the model image `phi`
  image: phi
  imagePullPolicy: IfNotPresent
  resources:
    limits:
      cpu: 4
      memory: 8Gi
      nvidia.com/gpu: 1 # If you got GPUs
    requests:
      cpu: 4
      memory: 8Gi
      nvidia.com/gpu: 1 # If you got GPUs
  storageClassName: local-path
  # If you have your own PersistentVolumeClaim created
  persistentVolumeClaim: your-pvc
  # If you need to specify the access mode for the PersistentVolume
  persistentVolume:
    accessMode: ReadWriteOnce
```
