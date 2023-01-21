## 参考文档
```shell
# k8s 官方校验文档
https://kubernetes.io/zh/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#validation

# githup 校验规范
https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md
```

## 基本写法
#### 最小值 & 最大值
```yaml
properties:
replicas:
  type: integer
  minimum: 1
  maximum: 20
```

#### 默认值
```yaml
dsn:
  type: string
  default: "abcde"
```

#### 设置必填
```yaml
spec:
  type: object
  properties:
    replicas:
      type: integer
      minimum: 1
      maximum: 20
    dsn:
      type: string
      default: "abcde"
  required:
    - replicas
    - dsn
```

#### 设置状态
```yaml
schema:
    openAPIV3Schema:
      type: object
      properties:
        spec:
          type: object
          properties:
            replicas:
              type: integer
              minimum: 1
              maximum: 20
            dsn:
              type: string
              default: "abcde"
          required:
            - replicas
            - dsn
        status:
          type: object
          properties:
            replicas:
              type: string
subresources:
  status: { }
```

#### kubectl get 显示属性
```yaml
versions:
- name: v1
  # 是否有效
  served: true
  storage: true
  additionalPrinterColumns:
    - name: replicas
      type: string
      jsonPath: .spec.status.replicas
    - name: Age
      type: date
      jsonPath: .metadata.creationTimestamp
```

#### kubectl scale 支持
```yaml
# kubectl scale --replicas=5 dc/mydbconfig
subresources:
    status: { }
    scale:
      # specReplicasPath 定义定制资源中对应 scale.spec.replicas 的 JSON 路径
      specReplicasPath: .spec.replicas
      # statusReplicasPath 定义定制资源中对应 scale.status.replicas 的 JSON 路径
      statusReplicasPath: .status.replicas
```

## 增加 kubectl get 打印字段步骤
```yaml
# 1. 增加属性
schema:
  openAPIV3Schema:
    type: object
    properties:
      spec:
        type: object
        properties:
          replicas:
            type: integer
            minimum: 1
            maximum: 20
          maxOpenConn:
            type: integer
            minimum: 1
            maximum: 2000
            default: 15
          maxIdleConn:
            type: integer
            minimum: 1
            maximum: 2000
            default: 5
          maxLifeTime:
            type: integer
            minimum: 60 #最小60秒
            default: 600  #默认10分钟
          dsn:
            type: string

# 2. 增加打印
versions:
  - name: v1
    # 是否有效
    served: true
    storage: true
    additionalPrinterColumns:
      - name: Ready
        type: string
        jsonPath: .status.ready
      - name: Age
        type: date
        jsonPath: .metadata.creationTimestamp
      - name: 最大连接数
        type: integer
        jsonPath: .spec.maxOpenConn

# 3. 生效
kubectl apply -f crd.yaml

# 4. 修改自定义资源结构体定义
... ...
```