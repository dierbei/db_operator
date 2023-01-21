package builders

import (
	"bytes"
	"context"
	"log"
	"text/template"

	configv1 "github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigMapBuilder struct {
	cm     *corev1.ConfigMap
	config *configv1.DbConfig
	client.Client
	DataKey string
}

func NewConfigMapBuilder(config *configv1.DbConfig, client client.Client) (*ConfigMapBuilder, error) {
	cm := &corev1.ConfigMap{}
	if err := client.Get(context.Background(), types.NamespacedName{Namespace: config.Namespace, Name: deployName(config.Name)}, cm); err != nil {
		//只需要 赋值 name 和namespace data不管，在apply 函数中处理
		cm.Name, cm.Namespace = deployName(config.Name), config.Namespace
		cm.Data = make(map[string]string) //搞一个空的map即可
	}
	return &ConfigMapBuilder{cm: cm, Client: client, config: config}, nil
}

func (this *ConfigMapBuilder) setOwner() *ConfigMapBuilder {
	this.cm.OwnerReferences = append(this.cm.OwnerReferences,
		metav1.OwnerReference{
			APIVersion: this.config.APIVersion,
			Kind:       this.config.Kind,
			Name:       this.config.Name,
			UID:        this.config.UID,
		})
	return this
}

const configMapKey = "app.yml"

func (this *ConfigMapBuilder) parseKey() *ConfigMapBuilder {

	if appData, ok := this.cm.Data[configMapKey]; ok {
		this.DataKey = Md5(appData)
		return this
	}
	this.DataKey = ""
	return this
}

func (this *ConfigMapBuilder) apply() *ConfigMapBuilder {
	tpl, err := template.New("appyaml").Delims("[[", "]]").Parse(cmtpl)
	if err != nil {
		log.Println(err)
		return this
	}

	var tplRet bytes.Buffer
	err = tpl.Execute(&tplRet, this.config.Spec)
	if err != nil {
		log.Println(err)
		return this
	}

	this.cm.Data[configMapKey] = tplRet.String()
	return this
}

func (this *ConfigMapBuilder) Build(ctx context.Context) error {
	if this.cm.CreationTimestamp.IsZero() {
		this.apply().setOwner().parseKey()
		err := this.Create(ctx, this.cm)
		if err != nil {
			return err
		}
	} else {
		patch := client.MergeFrom(this.cm.DeepCopy())
		this.apply().parseKey()
		err := this.Patch(ctx, this.cm, patch)
		if err != nil {
			return err
		}

	}
	return nil
}
