package dataintensive

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const (
	Name = "dataintensive"
)

var _ framework.ScorePlugin = &DataIntensive{}

type DataIntensive struct {
	handle framework.Handle
}

func New(_ runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &DataIntensive{
		handle: handle,
	}, nil
}

func (di *DataIntensive) Name() string {
	return Name
}