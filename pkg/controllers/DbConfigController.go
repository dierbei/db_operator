package controllers

import (
	"context"

	v1 "github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1"
	"github.com/shenyisyn/dbcore/pkg/builders"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DbConfigController struct {
	client.Client
}

func NewDbConfigController() *DbConfigController {
	return &DbConfigController{}
}

// Reconcile
// 事件循环：资源的 增、删、改 都会进入此方法
func (r *DbConfigController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	config := &v1.DbConfig{}

	// 查询对象，查询不到则返回
	if err := r.Get(ctx, req.NamespacedName, config); err != nil {
		return reconcile.Result{}, err
	}

	// 查询到对象则创建一个 builder
	builder, err := builders.NewDeployBuilder(config, req.Namespace, req.Name, r.Client)
	if err != nil {
		return reconcile.Result{}, err
	}

	// 执行最后同步操作
	if err := builder.Build(context.Background()); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, err
}

func (r *DbConfigController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}
