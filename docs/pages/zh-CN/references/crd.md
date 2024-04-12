### 完整 CRD 参考

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
  storageClassName: local-path
  # If you have your own PersistentVolumeClaim created
  persistentVolumeClaim: your-pvc
  # If you need to specify the access mode for the PersistentVolume
  persistentVolume:
    accessMode: ReadWriteOnce
```
