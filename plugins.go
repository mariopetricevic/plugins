package plugins

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)
type customFilterPlugin struct{
	handle framework.Handle
}


//var _ framework.FilterPlugin = &customFilterPlugin{}
var _  = framework.FilterPlugin(&customFilterPlugin{})

func (p *customFilterPlugin) Name() string {
	return "CustomFilter"
}



func (p *customFilterPlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	// Implementirajte logiku filtriranja čvorova ovdje
	return framework.NewStatus(framework.Success)
}

func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &customFilterPlugin{}, nil
}

