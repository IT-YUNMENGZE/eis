package delaysensitive

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	len "my.domain/Len/api/v1"
)

const (
	Name = "delaysensitive"
)

var (
	_ framework.QueueSortPlugin = &DelaySensitive{}
	_ framework.ScorePlugin     = &DelaySensitive{}
	_ framework.ScoreExtensions = &DelaySensitive{}

	scheme = runtime.NewScheme()
)

type DelaySensitive struct {
	handle framework.Handle
	cache  cache.Cache
}

func (ds *DelaySensitive) Name() string {
	return Name
}

func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	mgrConfig := ctrl.GetConfigOrDie()
	mgrConfig.QPS = 1000
	mgrConfig.Burst = 1000

	if err := len.AddToScheme(scheme); err != nil {
		klog.Error(err)
		return nil, err
	}

	mgr, err := ctrl.NewManager(mgrConfig, ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: "",
		LeaderElection:     false,
		Port:               9443,
	})
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	go func() {
		if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			klog.Error(err)
			panic(err)
		}
	}()

	lenCache := mgr.GetCache()

	if lenCache.WaitForCacheSync(context.TODO()) {
		return &DelaySensitive{
			handle: h,
			cache:  lenCache,
		}, nil
	} else {
		return nil, errors.New("cache not sync! ")
	}
}
