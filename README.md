## 技术栈
- code-generator v0.20.14
- golang 1.18.10

## 环境搭建
#### 1. 复制 code-genetator
```shell
# 下载 code-generator 源码 并复制到 $GOPATH/src/k8s.io/code-generator && go mod download

# 增加 shell 脚本可执行权限
chmod +x generate-groups.sh && ls
```

#### 2. 生成 deepcopy & client
```shell
# 复制 pkg go.mod main.go 到服务器目录 $GOPATH/src/github.com/shenyisyn/dbcore && go mod tidy

# 执行脚本
$GOPATH/src/k8s.io/code-generator/generate-groups.sh all  github.com/shenyisyn/dbcore/pkg/client github.com/shenyisyn/dbcore/pkg/apis dbconfig:v1

# 不能生成在 main.go 中加入 import _ "k8s.io/code-generator" && go mod tidy
```

## 控制器关键点
```text
1. Reconcile 函数（增删改查）
2. 启动 manager，并使用控制器
3. 进行资源监听（如何实现更新、删除、创建）
```

## 前台级联删除
```shell
kubectl delete dc mydb1 --cascade=foreground
```

## 应用程序启动参数
```shell
# 第一种
./dbcore --replicas=5 --dns=xxx ......

# 第二种
利用 ConfigMap 挂在到 Pod 里面
```

## 软件配置更新
```text
1. 第一种方法是软件支持实时读取（需要修改软件，并且万一该软件不在K8s中运行呢？）
2. 第二宗方法是控制器发现配置更新，删除Pod
3. 第三种方法是更新Deployment annotation
```

## 旧事件 & 子资源
```text
# 旧事件
可以使用名字 & UID 进行过滤

# 子资源
判断事件的 ID 是否为 Deployment 的 UID
```

## 指标查看
```shell
http://localhost:8082/metrics
```

## 参考链接
```shell
# code-generator 生成代码
https://github.com/kubernetes/code-generator
```
