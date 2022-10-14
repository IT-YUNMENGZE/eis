package computeintensive

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const (
	Name = "computeintensive"

	CoreWeight      = 3 //CPU逻辑核数量权重
	FrequencyWeight = 2 //CPU睿频权重
	CacheWeight     = 1 //CPU缓存权重
)

var _ framework.ScorePlugin = &ComputeIntensive{}

type ComputeIntensive struct {
	handle framework.Handle
}

func New(_ runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &ComputeIntensive{
		handle: handle,
	}, nil
}

func (ci *ComputeIntensive) Name() string {
	return Name
}