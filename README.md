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
```

## 参考链接
```shell
# code-generator 生成代码
https://github.com/kubernetes/code-generator
```
