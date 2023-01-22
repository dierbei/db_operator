package k8sconfig

import (
	"github.com/shenyisyn/dbcore/pkg/mymetrics"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	v1 "github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1"
	"github.com/shenyisyn/dbcore/pkg/controllers"
	"github.com/shenyisyn/dbcore/pkg/dashboard"

	appv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func InitManager() {
	logf.SetLogger(zap.New())
	mgr, err := manager.New(K8sRestConfig(), manager.Options{
		Logger: logf.Log.WithName("dbcore"),

		// 下面是选主参数
		LeaderElection:          true,
		LeaderElectionID:        "mydbcore-lock", // ConfigMap & lease 名字
		LeaderElectionNamespace: "default",       // 命名空间
		MetricsBindAddress:      ":8082",         // 端口，本地测试启动多个需要需改
	})
	if err != nil {
		mgr.GetLogger().Error(err, "unable to set up manager")
		os.Exit(1)
	}

	// 加入自定义指标
	metrics.Registry.MustRegister(mymetrics.MyReconcileTotal)

	// 加入 schema
	if err = v1.SchemeBuilder.AddToScheme(mgr.GetScheme()); err != nil {
		mgr.GetLogger().Error(err, "unable add scheme")
		os.Exit(1)
	}

	// 加入 controller
	dbconfigController := controllers.NewDbConfigController()
	if err = builder.ControllerManagedBy(mgr).For(&v1.DbConfig{}).Watches(&source.Kind{Type: &appv1.Deployment{}},
		handler.Funcs{
			DeleteFunc: dbconfigController.OnDelete,
			UpdateFunc: dbconfigController.OnUpdate,
		},
	).Complete(controllers.NewDbConfigController()); err != nil {
		mgr.GetLogger().Error(err, "unable to create manager")
		os.Exit(1)
	}

	// 加入 dashboard
	if err = mgr.Add(dashboard.NewAdminUi(mgr.GetClient())); err != nil {
		mgr.GetLogger().Error(err, "unable to create dashborad")
		os.Exit(1)
	}

	if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
		mgr.GetLogger().Error(err, "unable to start manager")
		os.Exit(1)
	}

}
