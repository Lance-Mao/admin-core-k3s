# Kubernetes 集群运维手册

完整的 K3s 集群运维、监控、故障排查和日志查询指南。

## 目录

- [集群状态查看](#集群状态查看)
- [应用状态监控](#应用状态监控)
- [日志查询](#日志查询)
- [故障排查](#故障排查)
- [性能监控](#性能监控)
- [资源管理](#资源管理)
- [数据备份与恢复](#数据备份与恢复)
- [应急处理](#应急处理)
- [常见问题解决方案](#常见问题解决方案)

---

## 集群状态查看

### 查看集群信息

```bash
# 查看集群基本信息
kubectl cluster-info

# 查看集群详细信息（包括健康检查）
kubectl cluster-info dump

# 查看 K3s 版本
kubectl version --short

# 查看 K3s 服务状态
systemctl status k3s
```

### 查看节点状态

```bash
# 查看所有节点
kubectl get nodes

# 查看节点详细信息
kubectl get nodes -o wide

# 查看节点资源使用情况
kubectl top nodes

# 查看节点详细配置和状态
kubectl describe node <node-name>

# 查看节点标签
kubectl get nodes --show-labels
```

### 查看命名空间

```bash
# 列出所有命名空间
kubectl get namespaces
# 或简写
kubectl get ns

# 查看某个命名空间的详细信息
kubectl describe namespace default
```

---

## 应用状态监控

### 查看 Pod 状态

```bash
# 查看所有 Pod
kubectl get pods

# 查看所有命名空间的 Pod
kubectl get pods --all-namespaces
# 或简写
kubectl get pods -A

# 查看 Pod 详细信息（包括 IP、节点等）
kubectl get pods -o wide

# 实时监控 Pod 状态变化
kubectl get pods -w

# 查看特定标签的 Pod
kubectl get pods -l app.kubernetes.io/component=api

# 查看 Pod 详细描述（包括事件）
kubectl describe pod <pod-name>

# 查看 Pod 的 YAML 配置
kubectl get pod <pod-name> -o yaml

# 查看 Pod 资源使用情况
kubectl top pods

# 按资源使用排序
kubectl top pods --sort-by=cpu
kubectl top pods --sort-by=memory
```

### 查看 Deployment 状态

```bash
# 查看所有 Deployment
kubectl get deployments
# 或简写
kubectl get deploy

# 查看 Deployment 详细信息
kubectl get deployments -o wide

# 查看 Deployment 历史版本
kubectl rollout history deployment/<deployment-name>

# 查看 Deployment 滚动更新状态
kubectl rollout status deployment/<deployment-name>

# 查看 Deployment 详细描述
kubectl describe deployment <deployment-name>
```

### 查看 StatefulSet 状态

```bash
# 查看所有 StatefulSet
kubectl get statefulsets
# 或简写
kubectl get sts

# 查看 StatefulSet 详细信息
kubectl describe statefulset <statefulset-name>
```

### 查看 Service 状态

```bash
# 查看所有服务
kubectl get services
# 或简写
kubectl get svc

# 查看服务详细信息（包括 Endpoints）
kubectl get svc -o wide

# 查看服务的 Endpoints
kubectl get endpoints

# 查看服务详细描述
kubectl describe service <service-name>
```

### 查看 Ingress 状态

```bash
# 查看所有 Ingress
kubectl get ingress
# 或简写
kubectl get ing

# 查看 Ingress 详细信息
kubectl describe ingress <ingress-name>

# 查看 Ingress Controller 状态
kubectl get pods -n ingress-nginx
```

### 查看持久化存储

```bash
# 查看 PersistentVolumeClaim
kubectl get pvc

# 查看 PersistentVolume
kubectl get pv

# 查看 StorageClass
kubectl get storageclass
# 或简写
kubectl get sc

# 查看 PVC 详细信息
kubectl describe pvc <pvc-name>
```

### 查看 ConfigMap 和 Secret

```bash
# 查看 ConfigMap
kubectl get configmaps
# 或简写
kubectl get cm

# 查看 ConfigMap 内容
kubectl get configmap <configmap-name> -o yaml

# 查看 Secret
kubectl get secrets

# 查看 Secret 详细信息（已编码）
kubectl get secret <secret-name> -o yaml

# 解码 Secret 的某个值
kubectl get secret <secret-name> -o jsonpath='{.data.password}' | base64 -d
```

---

## 日志查询

### 基础日志查询

```bash
# 查看 Pod 日志
kubectl logs <pod-name>

# 查看 Pod 中特定容器的日志
kubectl logs <pod-name> -c <container-name>

# 实时跟踪日志（类似 tail -f）
kubectl logs -f <pod-name>

# 查看最近的 N 行日志
kubectl logs <pod-name> --tail=100

# 查看过去某段时间的日志
kubectl logs <pod-name> --since=1h
kubectl logs <pod-name> --since=30m
kubectl logs <pod-name> --since=2023-01-01T00:00:00Z

# 查看前一个已崩溃容器的日志
kubectl logs <pod-name> --previous
# 或简写
kubectl logs <pod-name> -p
```

### 多容器日志查询

```bash
# 查看 Pod 所有容器的日志
kubectl logs <pod-name> --all-containers=true

# 查看 Init Container 的日志
kubectl logs <pod-name> -c <init-container-name>
```

### 按标签查询日志

```bash
# 查看某个标签的所有 Pod 日志
kubectl logs -l app.kubernetes.io/component=api

# 查看某个标签的所有 Pod 日志并实时跟踪
kubectl logs -l app.kubernetes.io/component=api -f

# 查看最近 100 行
kubectl logs -l app.kubernetes.io/component=api --tail=100
```

### 按 Deployment 查询日志

```bash
# 查看 Deployment 的日志（自动选择一个 Pod）
kubectl logs deployment/<deployment-name>

# 查看 Deployment 所有 Pod 的日志
kubectl logs deployment/<deployment-name> --all-containers=true

# 实时跟踪 Deployment 日志
kubectl logs -f deployment/<deployment-name>
```

### 日志搜索和过滤

```bash
# 查看日志并搜索关键字
kubectl logs <pod-name> | grep "error"
kubectl logs <pod-name> | grep -i "database"

# 查看日志并排除某些内容
kubectl logs <pod-name> | grep -v "debug"

# 统计错误日志数量
kubectl logs <pod-name> | grep -c "error"

# 查看包含错误的日志及其上下文
kubectl logs <pod-name> | grep -A 5 -B 5 "error"
```

### 导出日志到文件

```bash
# 导出日志到文件
kubectl logs <pod-name> > pod-logs.txt

# 导出最近 1000 行日志
kubectl logs <pod-name> --tail=1000 > pod-logs.txt

# 导出带时间戳的日志
kubectl logs <pod-name> --timestamps=true > pod-logs.txt
```

### K3s 系统日志

```bash
# 查看 K3s 服务日志
journalctl -u k3s -n 100

# 实时跟踪 K3s 日志
journalctl -u k3s -f

# 查看过去 1 小时的 K3s 日志
journalctl -u k3s --since "1 hour ago"

# 查看包含错误的 K3s 日志
journalctl -u k3s | grep -i error

# 导出 K3s 日志
journalctl -u k3s > k3s-logs.txt
```

---

## 故障排查

### Pod 启动失败排查

#### 1. ImagePullBackOff / ErrImagePull

**症状**：Pod 一直处于 `ImagePullBackOff` 状态

```bash
# 查看 Pod 详细信息
kubectl describe pod <pod-name>

# 常见原因和解决方案
```

**原因1：镜像不存在或名称错误**
```bash
# 检查镜像名称
kubectl get pod <pod-name> -o jsonpath='{.spec.containers[*].image}'

# 解决：修改 Deployment 使用正确的镜像
kubectl set image deployment/<deployment-name> <container-name>=<correct-image>
```

**原因2：镜像是私有的，缺少认证**
```bash
# 检查是否有 imagePullSecrets
kubectl get pod <pod-name> -o jsonpath='{.spec.imagePullSecrets}'

# 解决：将 GitHub Packages 设为 Public，或创建 imagePullSecrets
# 1. 前往 GitHub Packages 设置为 Public（推荐）
# 2. 或创建 Docker registry secret
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=<github-username> \
  --docker-password=<github-token> \
  --docker-email=<email>

# 然后在 Deployment 中引用这个 secret
```

**原因3：网络问题无法拉取镜像**
```bash
# 在节点上手动拉取镜像测试
ssh <node-ip>
docker pull <image-name>
```

#### 2. CrashLoopBackOff

**症状**：Pod 启动后立即崩溃，不断重启

```bash
# 查看 Pod 状态
kubectl get pods

# 查看 Pod 详细信息
kubectl describe pod <pod-name>

# 查看容器日志
kubectl logs <pod-name>

# 查看上一次崩溃的日志
kubectl logs <pod-name> --previous

# 查看所有容器的日志（包括 Init Container）
kubectl logs <pod-name> --all-containers=true
```

**常见原因**：
- 应用启动失败（配置错误、依赖服务不可用）
- 健康检查失败
- 资源不足（OOM）

```bash
# 检查资源限制
kubectl describe pod <pod-name> | grep -A 10 "Limits"

# 检查事件
kubectl get events --sort-by='.lastTimestamp' | grep <pod-name>
```

#### 3. Pending

**症状**：Pod 一直处于 `Pending` 状态

```bash
# 查看 Pod 详细信息，关注 Events 部分
kubectl describe pod <pod-name>

# 常见原因
```

**原因1：资源不足**
```bash
# 查看节点资源使用情况
kubectl top nodes
kubectl describe nodes

# 查看 Pod 请求的资源
kubectl describe pod <pod-name> | grep -A 5 "Requests"

# 解决：减少资源请求或增加节点资源
```

**原因2：PVC 无法绑定**
```bash
# 查看 PVC 状态
kubectl get pvc

# 查看 PVC 详细信息
kubectl describe pvc <pvc-name>

# 查看 PV 可用性
kubectl get pv

# 解决：检查 StorageClass 和 PV 配置
```

**原因3：节点选择器/污点无法满足**
```bash
# 查看 Pod 的节点选择器
kubectl get pod <pod-name> -o jsonpath='{.spec.nodeSelector}'

# 查看节点标签
kubectl get nodes --show-labels

# 查看节点污点
kubectl describe nodes | grep Taints
```

#### 4. Init 容器失败

```bash
# 查看 Init 容器状态
kubectl describe pod <pod-name> | grep -A 20 "Init Containers"

# 查看 Init 容器日志
kubectl logs <pod-name> -c <init-container-name>

# 常见场景：等待依赖服务
# 例如：API 等待 RPC 就绪
```

### 服务访问问题排查

#### 1. Service 无法访问

```bash
# 检查 Service 是否存在
kubectl get svc <service-name>

# 检查 Service Endpoints
kubectl get endpoints <service-name>

# 如果 Endpoints 为空，说明没有匹配的 Pod
# 检查 Service selector 和 Pod labels
kubectl get svc <service-name> -o jsonpath='{.spec.selector}'
kubectl get pods -l <label-selector> --show-labels

# 检查 Pod 是否就绪
kubectl get pods -l <label-selector>

# 测试服务连通性（从另一个 Pod）
kubectl run test-pod --image=busybox --rm -it -- sh
# 在 Pod 中执行
wget -O- http://<service-name>:<port>
nslookup <service-name>
```

#### 2. Ingress 无法访问

```bash
# 检查 Ingress 配置
kubectl get ingress
kubectl describe ingress <ingress-name>

# 检查 Ingress Controller 状态
kubectl get pods -n ingress-nginx
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller

# 检查 Ingress 后端服务
kubectl describe ingress <ingress-name> | grep -A 10 "Backend"

# 测试从服务器直接访问
curl http://localhost
curl http://<server-ip>
```

### 数据库连接问题排查

```bash
# 1. 检查 MySQL Pod 状态
kubectl get pods | grep mysql
kubectl describe pod <mysql-pod-name>

# 2. 检查 MySQL 日志
kubectl logs <mysql-pod-name>

# 3. 测试数据库连接
kubectl exec <mysql-pod-name> -- mysql -uroot -p<password> -e "SELECT 1;"

# 4. 检查 Secret 配置
kubectl get secret k3s-app-secret -o yaml
kubectl get secret k3s-app-secret -o jsonpath='{.data.mysql-password}' | base64 -d

# 5. 测试从应用 Pod 连接数据库
kubectl exec -it <app-pod-name> -- sh
# 在 Pod 中测试
ping <mysql-service-name>
telnet <mysql-service-name> 3306

# 6. 检查网络策略
kubectl get networkpolicies
```

### 磁盘空间问题

```bash
# 检查节点磁盘使用情况
df -h

# 查看 Docker/Containerd 占用空间
du -sh /var/lib/rancher/k3s

# 清理未使用的镜像
k3s crictl rmi --prune

# 查看 PVC 使用情况
kubectl get pvc
kubectl describe pvc <pvc-name>

# 查看 Pod 磁盘使用（需要进入 Pod）
kubectl exec <pod-name> -- df -h
```

---

## 性能监控

### 资源使用监控

```bash
# 查看节点资源使用
kubectl top nodes

# 查看所有 Pod 资源使用
kubectl top pods

# 查看特定命名空间的 Pod 资源使用
kubectl top pods -n <namespace>

# 按 CPU 使用排序
kubectl top pods --sort-by=cpu

# 按内存使用排序
kubectl top pods --sort-by=memory

# 查看容器级别的资源使用
kubectl top pods --containers
```

### 查看资源配额

```bash
# 查看 Pod 的资源请求和限制
kubectl describe pod <pod-name> | grep -A 10 "Limits"

# 查看 Deployment 的资源配置
kubectl get deployment <deployment-name> -o yaml | grep -A 10 "resources"

# 查看命名空间资源配额
kubectl get resourcequotas
kubectl describe resourcequota <quota-name>
```

### 性能指标

```bash
# 查看 API Server 性能
kubectl get --raw /metrics

# 查看 Pod 重启次数
kubectl get pods -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.containerStatuses[0].restartCount}{"\n"}{end}'

# 查看 Pod 运行时长
kubectl get pods -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.startTime}{"\n"}{end}'
```

---

## 资源管理

### 扩缩容操作

```bash
# 手动扩容 Deployment
kubectl scale deployment <deployment-name> --replicas=3

# 手动缩容
kubectl scale deployment <deployment-name> --replicas=1

# 扩容 StatefulSet
kubectl scale statefulset <statefulset-name> --replicas=3

# 查看扩缩容进度
kubectl get pods -w
```

### 滚动更新

```bash
# 更新镜像
kubectl set image deployment/<deployment-name> <container-name>=<new-image>

# 查看滚动更新状态
kubectl rollout status deployment/<deployment-name>

# 暂停滚动更新
kubectl rollout pause deployment/<deployment-name>

# 恢复滚动更新
kubectl rollout resume deployment/<deployment-name>

# 回滚到上一个版本
kubectl rollout undo deployment/<deployment-name>

# 回滚到特定版本
kubectl rollout undo deployment/<deployment-name> --to-revision=2

# 查看滚动更新历史
kubectl rollout history deployment/<deployment-name>
```

### 重启应用

```bash
# 重启 Deployment（滚动重启）
kubectl rollout restart deployment/<deployment-name>

# 强制删除 Pod（会自动重建）
kubectl delete pod <pod-name>

# 删除所有指定标签的 Pod
kubectl delete pod -l app.kubernetes.io/component=api

# 重启 StatefulSet（按顺序重启）
kubectl rollout restart statefulset/<statefulset-name>

# 强制删除卡住的 Pod
kubectl delete pod <pod-name> --force --grace-period=0
```

### 资源清理

```bash
# 删除 Deployment
kubectl delete deployment <deployment-name>

# 删除 Service
kubectl delete service <service-name>

# 删除 PVC（会删除数据！）
kubectl delete pvc <pvc-name>

# 删除所有已完成的 Pod
kubectl delete pods --field-selector=status.phase=Succeeded

# 删除所有失败的 Pod
kubectl delete pods --field-selector=status.phase=Failed

# 清理 Evicted 状态的 Pod
kubectl get pods | grep Evicted | awk '{print $1}' | xargs kubectl delete pod
```

---

## 数据备份与恢复

### MySQL 数据备份

```bash
# 获取 MySQL root 密码
ROOT_PWD=$(kubectl get secret k3s-app-secret -o jsonpath='{.data.mysql-root-password}' | base64 -d)

# 备份单个数据库
kubectl exec k3s-app-mysql-0 -- mysqldump \
  -uroot -p"$ROOT_PWD" \
  simple_admin > backup-$(date +%Y%m%d-%H%M%S).sql

# 备份所有数据库
kubectl exec k3s-app-mysql-0 -- mysqldump \
  -uroot -p"$ROOT_PWD" \
  --all-databases > backup-all-$(date +%Y%m%d-%H%M%S).sql

# 压缩备份
kubectl exec k3s-app-mysql-0 -- mysqldump \
  -uroot -p"$ROOT_PWD" \
  simple_admin | gzip > backup-$(date +%Y%m%d-%H%M%S).sql.gz
```

### MySQL 数据恢复

```bash
# 从备份恢复
kubectl exec -i k3s-app-mysql-0 -- mysql \
  -uroot -p"$ROOT_PWD" \
  simple_admin < backup.sql

# 从压缩备份恢复
gunzip < backup.sql.gz | kubectl exec -i k3s-app-mysql-0 -- mysql \
  -uroot -p"$ROOT_PWD" \
  simple_admin
```

### ConfigMap 和 Secret 备份

```bash
# 备份 ConfigMap
kubectl get configmap <configmap-name> -o yaml > configmap-backup.yaml

# 备份 Secret
kubectl get secret <secret-name> -o yaml > secret-backup.yaml

# 备份所有 ConfigMap 和 Secret
kubectl get configmap,secret -o yaml > configs-backup.yaml
```

### PVC 数据备份

```bash
# 创建临时 Pod 挂载 PVC 进行备份
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: backup-pod
spec:
  containers:
  - name: backup
    image: busybox
    command: ["sleep", "3600"]
    volumeMounts:
    - name: data
      mountPath: /data
  volumes:
  - name: data
    persistentVolumeClaim:
      claimName: <pvc-name>
EOF

# 等待 Pod 就绪
kubectl wait --for=condition=ready pod/backup-pod

# 将数据复制出来
kubectl exec backup-pod -- tar czf /tmp/backup.tar.gz /data
kubectl cp backup-pod:/tmp/backup.tar.gz ./backup.tar.gz

# 清理临时 Pod
kubectl delete pod backup-pod
```

---

## 应急处理

### 快速重启所有服务

```bash
# 重启所有应用 Deployment
kubectl rollout restart deployment -l app.kubernetes.io/instance=k3s-app

# 重启 MySQL
kubectl delete pod k3s-app-mysql-0

# 重启 Redis
kubectl rollout restart deployment k3s-app-redis
```

### 紧急扩容

```bash
# 快速扩容 API
kubectl scale deployment k3s-app-api --replicas=5

# 快速扩容 Frontend
kubectl scale deployment k3s-app-frontend --replicas=5
```

### 应急降级

```bash
# 回滚到上一个稳定版本
kubectl rollout undo deployment/k3s-app-api
kubectl rollout undo deployment/k3s-app-rpc

# 查看回滚状态
kubectl rollout status deployment/k3s-app-api
```

### 紧急排水节点（驱逐 Pod）

```bash
# 标记节点为不可调度
kubectl cordon <node-name>

# 驱逐节点上的所有 Pod
kubectl drain <node-name> --ignore-daemonsets --delete-emptydir-data

# 恢复节点为可调度
kubectl uncordon <node-name>
```

### K3s 服务重启

```bash
# 重启 K3s 服务
systemctl restart k3s

# 查看 K3s 状态
systemctl status k3s

# 查看 K3s 启动日志
journalctl -u k3s -n 100
```

---

## 常见问题解决方案

### 问题1：Pod 无法启动，OOMKilled

**症状**：Pod 被杀死，日志显示 `OOMKilled`

```bash
# 查看 Pod 状态
kubectl describe pod <pod-name> | grep -A 5 "State"

# 查看 Pod 资源限制
kubectl describe pod <pod-name> | grep -A 10 "Limits"
```

**解决方案**：
```bash
# 方案1：增加内存限制（编辑 values.yaml 或 Deployment）
kubectl edit deployment <deployment-name>
# 修改 resources.limits.memory

# 方案2：临时扩容（重新部署后失效）
kubectl set resources deployment <deployment-name> -c=<container-name> --limits=memory=512Mi
```

### 问题2：磁盘空间不足

```bash
# 检查磁盘使用
df -h

# 清理 Docker 镜像
k3s crictl images | grep '<none>' | awk '{print $3}' | xargs k3s crictl rmi

# 清理日志
journalctl --vacuum-time=7d

# 清理 Pod 日志（需要进入节点）
find /var/log/pods/ -name "*.log" -mtime +7 -delete
```

### 问题3：Service 无法解析

```bash
# 测试 DNS 解析
kubectl run -it --rm debug --image=busybox --restart=Never -- nslookup <service-name>

# 查看 CoreDNS 日志
kubectl logs -n kube-system -l k8s-app=kube-dns

# 重启 CoreDNS
kubectl rollout restart -n kube-system deployment/coredns
```

### 问题4：证书过期

```bash
# 查看证书有效期
openssl x509 -in /etc/rancher/k3s/k3s.yaml -noout -dates

# K3s 会自动更新证书，如果需要手动更新
k3s certificate rotate
systemctl restart k3s
```

### 问题5：Helm 部署卡住

```bash
# 查看 Helm Release 状态
helm list

# 查看 Helm Release 历史
helm history k3s-app

# 回滚 Helm Release
helm rollback k3s-app <revision>

# 强制删除卡住的 Release
helm delete k3s-app --no-hooks
```

---

## 常用命令速查表

### 快速查看

```bash
# 一键查看所有关键资源
kubectl get all

# 查看所有资源（包括 ConfigMap、Secret 等）
kubectl get all,cm,secret,pvc

# 查看最近的事件
kubectl get events --sort-by='.lastTimestamp' | tail -20

# 查看问题 Pod
kubectl get pods | grep -v "Running\|Completed"
```

### 快速进入 Pod

```bash
# 进入 Pod 的 Shell
kubectl exec -it <pod-name> -- sh
# 或
kubectl exec -it <pod-name> -- bash

# 执行单个命令
kubectl exec <pod-name> -- ls -la

# 进入 MySQL
kubectl exec -it k3s-app-mysql-0 -- mysql -usimple_admin -p
```

### 端口转发（本地调试）

```bash
# 将 Pod 端口转发到本地
kubectl port-forward pod/<pod-name> 8080:80

# 将 Service 端口转发到本地
kubectl port-forward service/<service-name> 8080:80

# 转发 MySQL 端口（在后台运行）
kubectl port-forward k3s-app-mysql-0 3306:3306 &

# 然后可以用本地工具连接
mysql -h 127.0.0.1 -P 3306 -usimple_admin -p
```

### 文件复制

```bash
# 从 Pod 复制文件到本地
kubectl cp <pod-name>:/path/to/file ./local-file

# 从本地复制文件到 Pod
kubectl cp ./local-file <pod-name>:/path/to/file

# 复制整个目录
kubectl cp <pod-name>:/path/to/dir ./local-dir
```

### 临时调试 Pod

```bash
# 创建临时调试 Pod
kubectl run debug --image=busybox --rm -it --restart=Never -- sh

# 创建带网络工具的调试 Pod
kubectl run debug --image=nicolaka/netshoot --rm -it --restart=Never -- bash

# 在调试 Pod 中测试服务连通性
wget -O- http://<service-name>:<port>
curl http://<service-name>:<port>
nslookup <service-name>
ping <service-name>
```

---

## 监控脚本

### 健康检查脚本

保存为 `health-check.sh`：

```bash
#!/bin/bash

echo "=== Cluster Health Check ==="
echo ""

echo "--- Nodes Status ---"
kubectl get nodes

echo ""
echo "--- Pods Status ---"
kubectl get pods -o wide | grep -v "Running.*1/1\|Running.*2/2"

echo ""
echo "--- System Pods ---"
kubectl get pods -n kube-system

echo ""
echo "--- Recent Events (Last 10) ---"
kubectl get events --sort-by='.lastTimestamp' | tail -10

echo ""
echo "--- Resource Usage ---"
kubectl top nodes
kubectl top pods

echo ""
echo "--- PVC Status ---"
kubectl get pvc

echo ""
echo "--- Service Endpoints ---"
kubectl get endpoints

echo ""
echo "=== Check Complete ==="
```

使用方法：
```bash
chmod +x health-check.sh
./health-check.sh
```

---

## 总结

本文档涵盖了 Kubernetes 集群的日常运维操作，包括：

- ✅ 集群和应用状态监控
- ✅ 完整的日志查询方案
- ✅ 系统化的故障排查流程
- ✅ 性能监控和资源管理
- ✅ 数据备份与恢复
- ✅ 应急处理方案
- ✅ 常见问题解决方案

**建议**：
1. 将常用命令添加到 shell 别名
2. 定期备份关键数据
3. 监控集群资源使用情况
4. 建立告警机制

**相关文档**：
- [部署教程](DEPLOYMENT.md)
- [详细配置](SETUP.md)
- [项目说明](README.md)
