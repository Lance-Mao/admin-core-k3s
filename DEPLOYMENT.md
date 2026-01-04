# Simple Admin Core - 完整部署教程

从零开始在新服务器上部署 Simple Admin Core 项目的完整指南。

## 目录

- [系统要求](#系统要求)
- [架构概览](#架构概览)
- [第一步：服务器准备](#第一步服务器准备)
- [第二步：安装 K3s](#第二步安装-k3s)
- [第三步：配置 GitHub](#第三步配置-github)
- [第四步：部署应用](#第四步部署应用)
- [第五步：验证部署](#第五步验证部署)
- [故障排查](#故障排查)
- [日常运维](#日常运维)

---

## 系统要求

### 服务器配置

- **操作系统**: Ubuntu 20.04/22.04 LTS 或 Debian 11/12
- **CPU**: 最低 2 核（推荐 4 核）
- **内存**: 最低 4GB（推荐 8GB）
- **硬盘**: 最低 40GB（推荐 80GB+）
- **网络**: 公网 IP 地址

### 本地环境

- Git 客户端
- GitHub 账号
- SSH 客户端

---

## 架构概览

```
┌─────────────────────────────────────────────────────────────────┐
│                    VPS 服务器 (K3s 集群)                          │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │           Nginx Ingress Controller                       │    │
│  │              (端口 80/443)                                │    │
│  └────────────────────┬────────────────────────────────────┘    │
│                       │                                          │
│          ┌────────────┴──────────────┐                          │
│          │ /                         │ /api/                    │
│          ▼                           ▼                          │
│   ┌─────────────┐            ┌─────────────┐                   │
│   │  Frontend   │            │   API       │                   │
│   │  (Vue.js)   │            │  (Go)       │                   │
│   │   Nginx     │            │  Port 9100  │                   │
│   └─────────────┘            └──────┬──────┘                   │
│                                     │                           │
│                              ┌──────┴──────┐                    │
│                              │             │                    │
│                              ▼             ▼                    │
│                       ┌─────────────┐ ┌─────────────┐          │
│                       │    RPC      │ │   MySQL     │          │
│                       │   (Go)      │ │  Database   │          │
│                       │  Port 9101  │ │  Port 3306  │          │
│                       └──────┬──────┘ └─────────────┘          │
│                              │                                  │
│                              ▼                                  │
│                       ┌─────────────┐                          │
│                       │   Redis     │                          │
│                       │   Cache     │                          │
│                       │  Port 6379  │                          │
│                       └─────────────┘                          │
└─────────────────────────────────────────────────────────────────┘
```

### 组件说明

- **Frontend**: Vue.js 前端应用（Nginx 服务）
- **API**: Go 语言编写的 RESTful API 服务
- **RPC**: Go 语言编写的 gRPC 服务
- **MySQL**: 数据持久化存储
- **Redis**: 缓存和会话存储
- **Nginx Ingress**: 统一入口和负载均衡

---

## 第一步：服务器准备

### 1.1 连接服务器

```bash
# 使用 SSH 连接到你的 VPS
ssh root@YOUR_SERVER_IP
```

### 1.2 更新系统

```bash
# 更新软件包列表
apt update && apt upgrade -y

# 安装必要工具
apt install -y curl wget git vim ufw
```

### 1.3 配置防火墙

```bash
# 启用防火墙
ufw --force enable

# 允许 SSH（重要！避免被锁在外面）
ufw allow 22/tcp

# 允许 HTTP/HTTPS
ufw allow 80/tcp
ufw allow 443/tcp

# 允许 Kubernetes API（用于 CI/CD 访问）
ufw allow 6443/tcp

# 查看防火墙状态
ufw status
```

预期输出：
```
Status: active

To                         Action      From
--                         ------      ----
22/tcp                     ALLOW       Anywhere
80/tcp                     ALLOW       Anywhere
443/tcp                    ALLOW       Anywhere
6443/tcp                   ALLOW       Anywhere
```

---

## 第二步：安装 K3s

### 2.1 安装 K3s

K3s 是轻量级 Kubernetes 发行版，非常适合单节点部署。

```bash
# 安装 K3s（禁用 Traefik，我们使用 Nginx Ingress）
curl -sfL https://get.k3s.io | sh -s - --disable traefik

# 等待 K3s 启动
sleep 30

# 检查 K3s 状态
systemctl status k3s
```

### 2.2 配置 kubectl

```bash
# 创建 .kube 目录
mkdir -p ~/.kube

# 复制 K3s 配置文件
cp /etc/rancher/k3s/k3s.yaml ~/.kube/config

# 设置正确的权限
chown $(id -u):$(id -g) ~/.kube/config

# 设置环境变量
echo 'export KUBECONFIG=~/.kube/config' >> ~/.bashrc
source ~/.bashrc

# 测试 kubectl
kubectl get nodes
```

预期输出：
```
NAME         STATUS   ROLES                  AGE   VERSION
your-node    Ready    control-plane,master   1m    v1.28.x+k3s1
```

### 2.3 安装 Helm

```bash
# 下载 Helm 安装脚本
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# 验证安装
helm version
```

### 2.4 准备 Kubeconfig 用于 GitHub Actions

```bash
# 获取你的服务器公网 IP
SERVER_IP=$(curl -s ifconfig.me)
echo "服务器 IP: $SERVER_IP"

# 生成用于 GitHub Actions 的 kubeconfig（base64 编码）
cat /etc/rancher/k3s/k3s.yaml | \
  sed "s/127.0.0.1/$SERVER_IP/g" | \
  base64 -w 0

# 复制输出的 base64 字符串，稍后会用到
```

**重要：** 保存这个 base64 字符串，我们将在 GitHub Secrets 中使用它。

---

## 第三步：配置 GitHub

### 3.1 Fork 或 Clone 项目

如果你还没有项目代码：

```bash
# 在本地克隆项目
git clone https://github.com/YOUR_USERNAME/admin-core-k3s.git
cd admin-core-k3s
```

### 3.2 配置 GitHub Secrets

前往 GitHub 仓库页面：`Settings` → `Secrets and variables` → `Actions`

点击 `New repository secret` 添加以下 secrets：

| Secret 名称 | 说明 | 示例值 |
|------------|------|--------|
| `KUBECONFIG` | K3s 集群配置（必需） | 第二步获取的 base64 字符串 |
| `MYSQL_ROOT_PASSWORD` | MySQL root 密码（必需） | `MyS3cureR00tP@ssw0rd!2024` |
| `MYSQL_USERNAME` | MySQL 用户名（可选） | `simple_admin`（默认值） |
| `MYSQL_PASSWORD` | MySQL 用户密码（必需） | `AppDB@ss123` |
| `JWT_SECRET` | JWT 密钥（必需） | `D8eGr6LRnF2dSM8aDRcyzFaRLfkoNxWBwOTkRaIm5NPkks6xzx0fnEKGLrionGEG` |

**生成安全密码的方法：**

```bash
# 生成 MySQL root 密码
openssl rand -base64 32

# 生成 MySQL 用户密码
openssl rand -base64 24

# 生成 JWT Secret（64字符）
openssl rand -base64 48 | tr -d '\n' | head -c 64
```

### 3.3 修改 values.yaml（如果需要）

编辑 `helm/k3s-app/values.yaml`：

```yaml
# 镜像仓库配置（修改为你的 GitHub 用户名）
frontend:
  image:
    repository: YOUR_GITHUB_USERNAME/core-frontend

api:
  image:
    repository: YOUR_GITHUB_USERNAME/core-api

rpc:
  image:
    repository: YOUR_GITHUB_USERNAME/core-rpc
```

---

## 第四步：部署应用

### 4.1 首次部署

```bash
# 提交更改（如果有修改 values.yaml）
git add .
git commit -m "feat: initial deployment setup"
git push origin main
```

### 4.2 监控部署过程

1. **查看 GitHub Actions**
   - 前往仓库的 `Actions` 标签页
   - 查看工作流运行状态

2. **在服务器上监控 Pod 创建**

```bash
# 实时查看 Pod 状态
watch kubectl get pods

# 查看 Helm 部署状态
helm list

# 查看详细事件
kubectl get events --sort-by='.lastTimestamp'
```

### 4.3 处理 ImagePullBackOff（首次部署常见）

首次部署时，GitHub 默认将容器镜像设为私有，需要改为公开：

1. 访问 `https://github.com/YOUR_USERNAME?tab=packages`
2. 点击每个包（`core-frontend`、`core-api`、`core-rpc`）
3. 进入 **Package settings** → **Danger Zone** → **Change visibility**
4. 选择 **Public** → 确认

然后重新运行 GitHub Action 或重新推送代码。

### 4.4 等待部署完成

```bash
# 等待所有 Pod 就绪
kubectl wait --for=condition=ready pod --all --timeout=600s

# 查看最终状态
kubectl get pods,svc
```

预期输出：
```
NAME                                    READY   STATUS    RESTARTS   AGE
pod/k3s-app-api-xxx                     1/1     Running   0          2m
pod/k3s-app-api-yyy                     1/1     Running   0          2m
pod/k3s-app-frontend-xxx                1/1     Running   0          2m
pod/k3s-app-frontend-yyy                1/1     Running   0          2m
pod/k3s-app-mysql-0                     1/1     Running   0          2m
pod/k3s-app-redis-xxx                   1/1     Running   0          2m
pod/k3s-app-rpc-xxx                     1/1     Running   0          2m
pod/k3s-app-rpc-yyy                     1/1     Running   0          2m

NAME                       TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)
service/k3s-app-api        ClusterIP   10.43.x.x       <none>        9100/TCP
service/k3s-app-frontend   ClusterIP   10.43.x.x       <none>        80/TCP
service/k3s-app-mysql      ClusterIP   10.43.x.x       <none>        3306/TCP
service/k3s-app-redis      ClusterIP   10.43.x.x       <none>        6379/TCP
service/k3s-app-rpc        ClusterIP   10.43.x.x       <none>        9101/TCP
```

---

## 第五步：验证部署

### 5.1 检查 Ingress

```bash
# 查看 Ingress 配置
kubectl get ingress

# 获取外部 IP（可能需要等待几分钟）
kubectl get svc -n ingress-nginx
```

### 5.2 访问应用

在浏览器中访问：

```
http://YOUR_SERVER_IP/
```

你应该能看到前端页面。

### 5.3 测试 API

```bash
# 测试 API 健康检查
curl http://YOUR_SERVER_IP/api/health

# 测试 API 基本功能
curl http://YOUR_SERVER_IP/api/v1/core/init/database
```

### 5.4 初始化数据库

首次部署后，需要初始化数据库：

```bash
# 方法1：通过 API 初始化
curl -X POST http://YOUR_SERVER_IP/api/v1/core/init/database

# 方法2：通过前端界面
# 访问 http://YOUR_SERVER_IP/ 并按照提示进行初始化
```

---

## 故障排查

### 问题 1：Pod 一直处于 Pending 状态

**原因**：资源不足或存储问题

**解决方案**：
```bash
# 查看 Pod 详细信息
kubectl describe pod POD_NAME

# 检查节点资源
kubectl top nodes
kubectl describe nodes
```

### 问题 2：Pod 处于 CrashLoopBackOff

**原因**：应用启动失败

**解决方案**：
```bash
# 查看 Pod 日志
kubectl logs POD_NAME

# 查看前一次崩溃的日志
kubectl logs POD_NAME --previous

# 常见问题：
# - MySQL 连接失败：检查 Secret 中的密码是否正确
# - RPC 服务未就绪：等待 RPC Pod 先启动
```

### 问题 3：MySQL 认证失败

**错误信息**：`Access denied for user 'simple_admin'`

**原因**：数据库密码与 Secret 不匹配

**解决方案**：
```bash
# 删除 MySQL 数据（测试环境）
kubectl scale statefulset k3s-app-mysql --replicas=0
kubectl delete pvc data-k3s-app-mysql-0
kubectl scale statefulset k3s-app-mysql --replicas=1

# 等待 MySQL 重新初始化
kubectl wait --for=condition=ready pod/k3s-app-mysql-0 --timeout=120s

# 重启 API Pod
kubectl delete pod -l app.kubernetes.io/component=api
```

### 问题 4：无法访问应用（404）

**原因**：Ingress 配置问题

**解决方案**：
```bash
# 检查 Ingress 状态
kubectl describe ingress k3s-app

# 检查 Ingress Controller
kubectl get pods -n ingress-nginx

# 重启 Ingress Controller
kubectl rollout restart deployment -n ingress-nginx ingress-nginx-controller
```

### 问题 5：GitHub Actions 部署超时

**原因**：API Server 连接问题

**解决方案**：
```bash
# 在服务器上检查
kubectl cluster-info

# 检查防火墙
ufw status | grep 6443

# 确保 KUBECONFIG Secret 中的 IP 地址正确
# 重新生成并更新 GitHub Secret
```

---

## 日常运维

### 查看应用状态

```bash
# 查看所有 Pod
kubectl get pods

# 查看服务
kubectl get svc

# 查看 Ingress
kubectl get ingress

# 查看 PVC（持久化存储）
kubectl get pvc
```

### 查看日志

```bash
# 查看 API 日志
kubectl logs -f deployment/k3s-app-api

# 查看 Frontend 日志
kubectl logs -f deployment/k3s-app-frontend

# 查看 RPC 日志
kubectl logs -f deployment/k3s-app-rpc

# 查看 MySQL 日志
kubectl logs k3s-app-mysql-0

# 查看最近的日志（最后100行）
kubectl logs --tail=100 POD_NAME
```

### 重启服务

```bash
# 重启 API
kubectl rollout restart deployment k3s-app-api

# 重启 Frontend
kubectl rollout restart deployment k3s-app-frontend

# 重启 RPC
kubectl rollout restart deployment k3s-app-rpc

# 重启 Redis
kubectl rollout restart deployment k3s-app-redis

# 重启 MySQL（谨慎操作）
kubectl delete pod k3s-app-mysql-0
```

### 扩缩容

```bash
# 扩展 API 副本数
kubectl scale deployment k3s-app-api --replicas=3

# 扩展 Frontend 副本数
kubectl scale deployment k3s-app-frontend --replicas=3

# 或者修改 values.yaml 并重新部署
```

### 更新应用

```bash
# 本地修改代码后
git add .
git commit -m "feat: your changes"
git push origin main

# GitHub Actions 会自动构建并部署
```

### 手动部署（不通过 GitHub Actions）

```bash
# 在服务器上
cd /path/to/admin-core-k3s

# 拉取最新代码
git pull

# 更新 Helm 部署
helm upgrade k3s-app ./helm/k3s-app \
  --set secrets.mysql.rootPassword="YOUR_ROOT_PASSWORD" \
  --set secrets.mysql.password="YOUR_PASSWORD" \
  --set secrets.jwt.accessSecret="YOUR_JWT_SECRET"
```

### 备份 MySQL 数据

```bash
# 导出数据库
kubectl exec k3s-app-mysql-0 -- mysqldump \
  -uroot -p"YOUR_ROOT_PASSWORD" \
  simple_admin > backup-$(date +%Y%m%d).sql

# 复制备份到本地
kubectl cp k3s-app-mysql-0:/backup.sql ./mysql-backup.sql
```

### 恢复 MySQL 数据

```bash
# 复制备份到 Pod
kubectl cp ./mysql-backup.sql k3s-app-mysql-0:/tmp/restore.sql

# 恢复数据
kubectl exec k3s-app-mysql-0 -- mysql \
  -uroot -p"YOUR_ROOT_PASSWORD" \
  simple_admin < /tmp/restore.sql
```

### 卸载应用

```bash
# 卸载 Helm Release
helm uninstall k3s-app

# 删除持久化数据（可选）
kubectl delete pvc --all

# 卸载 K3s（完全清理）
/usr/local/bin/k3s-uninstall.sh
```

---

## 性能优化建议

### 1. 调整副本数

根据实际负载调整各组件的副本数：

```yaml
# helm/k3s-app/values.yaml
frontend:
  replicaCount: 3  # 建议 2-4

api:
  replicaCount: 3  # 建议 2-4

rpc:
  replicaCount: 2  # 建议 2-3
```

### 2. 配置资源限制

```yaml
# helm/k3s-app/values.yaml
api:
  resources:
    requests:
      memory: "128Mi"
      cpu: "100m"
    limits:
      memory: "256Mi"
      cpu: "500m"
```

### 3. 启用 Redis 持久化（可选）

如果需要 Redis 数据持久化：

```yaml
redis:
  persistence:
    enabled: true
    size: 5Gi
```

### 4. MySQL 性能调优

```yaml
mysql:
  config:
    dbMaxOpenConn: 100
    # 根据实际情况调整
```

---

## 安全建议

1. **修改默认密码**：确保所有 Secret 使用强密码
2. **启用 HTTPS**：使用 Let's Encrypt 配置 SSL 证书（参考 SETUP.md）
3. **限制访问**：配置防火墙规则，只允许必要的端口
4. **定期更新**：保持系统和组件更新到最新版本
5. **监控日志**：定期检查应用日志，发现异常行为

---

## 监控和告警（可选）

### 安装 Prometheus + Grafana

```bash
# 添加 Helm 仓库
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# 安装监控栈
helm install monitoring prometheus-community/kube-prometheus-stack \
  --namespace monitoring --create-namespace \
  --set grafana.service.type=NodePort \
  --set grafana.service.nodePort=30030 \
  --set grafana.adminPassword=admin123

# 开放端口
ufw allow 30030/tcp

# 访问 Grafana
# http://YOUR_SERVER_IP:30030
# 用户名: admin
# 密码: admin123
```

---

## 常见问题 FAQ

### Q1: 部署需要多长时间？

A: 首次部署通常需要 10-15 分钟，包括镜像下载和容器启动时间。

### Q2: 可以在本地开发环境部署吗？

A: 可以，使用 minikube 或 kind 创建本地 Kubernetes 集群，然后按相同步骤部署。

### Q3: 如何查看应用使用的资源？

A: 使用 `kubectl top pods` 和 `kubectl top nodes` 命令。

### Q4: 数据会丢失吗？

A: MySQL 使用 PVC（持久化卷），数据不会因 Pod 重启而丢失。但建议定期备份。

### Q5: 如何更新 Kubernetes 版本？

A: K3s 可以通过重新运行安装脚本升级：
```bash
curl -sfL https://get.k3s.io | sh -s - --disable traefik
```

---

## 获取帮助

- **项目文档**: [README.md](README.md)
- **详细配置**: [SETUP.md](SETUP.md)
- **问题反馈**: 提交 GitHub Issue

---

## 总结

现在你已经成功部署了 Simple Admin Core 项目！

**下一步建议：**
1. ✅ 配置域名和 SSL 证书（参考 SETUP.md）
2. ✅ 设置监控和告警
3. ✅ 配置自动备份
4. ✅ 优化性能参数

祝你使用愉快！ 🎉
