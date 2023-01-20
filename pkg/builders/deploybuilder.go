package builders

import (
	"bytes"
	"context"
	"text/template"

	configv1 "github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DeployBuilder
// 此对象主要用来控制 deployment
type DeployBuilder struct {
	deploy *v1.Deployment
	client.Client
	config *configv1.DbConfig //新增属性 。保存 config 对象
}

// deployName
// 给创建出来的 deployment 名称增加前缀
func deployName(name string) string {
	return "dbcore-" + name
}

func NewDeployBuilder(config *configv1.DbConfig, c client.Client) (*DeployBuilder, error) {
	deploy := &v1.Deployment{}

	if err := c.Get(context.Background(), types.NamespacedName{Name: deployName(config.Name), Namespace: config.Namespace}, deploy); err != nil {
		// 如果没有查询到 deployment 则进行第一次模板渲染
		deploy.Name, deploy.Namespace = deployName(config.Name), config.Namespace
		tpl, err := template.New("deploy").Parse(deptpl)
		if err != nil {
			return nil, err
		}

		var tplRet bytes.Buffer
		if err = tpl.Execute(&tplRet, deploy); err != nil {
			return nil, err
		}

		if err = yaml.Unmarshal(tplRet.Bytes(), deploy); err != nil {
			return nil, err
		}
	}

	return &DeployBuilder{
		deploy: deploy,
		Client: c,
		config: config,
	}, nil
}

// setOwner
// 设置 deployment 实体的 ownerreferences 属性，主要为了实现级联删除
func (this *DeployBuilder) setOwner() *DeployBuilder {
	this.deploy.OwnerReferences = append(this.deploy.OwnerReferences,
		metav1.OwnerReference{
			APIVersion: this.config.APIVersion,
			Kind:       this.config.Kind,
			Name:       this.config.Name,
			UID:        this.config.UID,
		})
	return this
}

// apply
// 修改 deployment replicas
func (this *DeployBuilder) apply() *DeployBuilder {
	*this.deploy.Spec.Replicas = int32(this.config.Spec.Replicas)
	return this
}

// Replicas
// 修改 deployment replicas，目前次方法已没使用
func (this *DeployBuilder) Replicas(r int) *DeployBuilder {
	*this.deploy.Spec.Replicas = int32(r)
	return this
}

// Build
// 最后的创建 deployment 动作
// CreationTimestamp 为空表示为第一次创建
// 否则则是更新对象
func (this *DeployBuilder) Build(ctx context.Context) error {
	if this.deploy.CreationTimestamp.IsZero() {
		this.apply().setOwner()
		err := this.Create(ctx, this.deploy)
		if err != nil {
			return err
		}
	} else {
		patch := client.MergeFrom(this.deploy.DeepCopy())
		this.apply() //同步  所需要的属性 如 副本数
		err := this.Patch(ctx, this.deploy, patch)
		if err != nil {
			return err
		}
	}
	return nil
}
