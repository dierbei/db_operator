package builders

const deptpl = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace}}
spec:
  selector:
    matchLabels:
      app: {{ .Namespace}}-{{ .Name }}
  replicas: 1
  template:
    metadata:
      labels:
        app: {{ .Namespace}}-{{ .Name }}
        version: v1
    spec:
      containers:
        - name: {{ .Namespace}}-{{ .Name }}-container
          image: docker.io/shenyisyn/dbcore:v1
          imagePullPolicy: IfNotPresent
          ports:
             - containerPort: 8081
             - containerPort: 8090


`
