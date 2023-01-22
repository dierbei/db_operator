package controllers

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shenyisyn/dbcore/pkg/mymetrics"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/event"

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

	// 增加指标
	mymetrics.MyReconcileTotal.With(prometheus.Labels{
		"controller": "dbconfig",
	}).Inc()

	// 查询对象，查询不到则返回
	if err := r.Get(ctx, req.NamespacedName, config); err != nil {
		return reconcile.Result{}, err
	}

	// 查询到对象则创建一个 builder
	builder, err := builders.NewDeployBuilder(config, r.Client)
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

func (r *DbConfigController) OnDelete(event event.DeleteEvent,
	limitingInterface workqueue.RateLimitingInterface) {
	for _, ref := range event.Object.GetOwnerReferences() {

		if ref.Kind == "DbConfig" && ref.APIVersion == "api.jtthink.com/v1" {
			limitingInterface.Add(
				reconcile.Request{
					types.NamespacedName{
						Name: ref.Name, Namespace: event.Object.GetNamespace(),
					},
				})
		}
	}
}

func (r *DbConfigController) OnUpdate(event event.UpdateEvent,
	limitingInterface workqueue.RateLimitingInterface) {
	for _, ref := range event.ObjectNew.GetOwnerReferences() {
		log.Println("deployment update")
		if ref.Kind == "DbConfig" && ref.APIVersion == "api.jtthink.com/v1" {
			limitingInterface.Add(
				reconcile.Request{
					types.NamespacedName{
						Name: ref.Name, Namespace: event.ObjectNew.GetNamespace(),
					},
				})
		}
	}
}
