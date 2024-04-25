# `Model` 模型资源

`Model` 是一个自定义资源定义（CRD），代表集群中的一个 Ollama 推理服务实例。

创建时，`Model` 的创建会触发对 [`Deployment`](https://kubernetes.io/zh-cn/docs/concepts/workloads/controllers/deployment/)，[`Service`](https://kubernetes.io/docs/concepts/services-networking/service/) 的创建。

> 如果是首次创建，还会额外创建一个用于存储和持久化已下载 Ollama 镜像的 StatefulSet 和其对应的 PersistentVolumeClaim。

## 完整 CRD 参考

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
