::: tip 没有现成的 Kubernetes 集群吗？

运行以下命令以在您的本地机器上安装 Docker 和 kind 并创建一个 Kubernetes 集群：

::: code-group

```shell [macOS]
brew install --cask docker
brew install docker kind kubectl
wget https://raw.githubusercontent.com/nekomeowww/ollama-operator/main/hack/kind-config.yaml
kind create cluster --config kind-config.yaml
```

```powershell [Windows]
Invoke-WebRequest  -OutFile "./Docker Desktop Installer.exe"
Start-Process 'Docker Desktop Installer.exe' -Wait install
start /w "" "Docker Desktop Installer.exe" install

# If you use Scoop command line installer
scoop install docker kubectl go
# Alternatively, if you use Chocolatey as package manager
choco install docker-desktop kubernetes-cli golang

go install sigs.k8s.io/kind@latest
wget https://raw.githubusercontent.com/nekomeowww/ollama-operator/main/hack/kind-config.yaml
kind create cluster --config kind-config.yaml
```

```shell [Linux]
# refer to Install Docker Engine on Debian | Docker Docs https://docs.docker.com/engine/install/debian/
# and Install and Set Up kubectl on Linux | Kubernetes https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/
```

:::
