package plugins

import (
	"fmt"
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)
type customFilterPlugin struct{
	handle framework.Handle
}

const (
	// Name : name of plugin used in the plugin registry and configurations.
	Name = "CustomFilterPlugin"
)


//var _ framework.FilterPlugin = &customFilterPlugin{}
var _  = framework.FilterPlugin(&customFilterPlugin{})

func (p *customFilterPlugin) Name() string {
	return Name
}


//func (s *customFilterPlugin) PreFilter(ctx context.Context, pod *v1.Pod) *framework.Status {
//	return framework.NewStatus(framework.Success, "")
//}


func (p *customFilterPlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	
	fmt.Println("-----------------INFORMACIJE O ČVORU-----------------")
	//resursi cvora
	nodeCpu := nodeInfo.Node().Status.Capacity[v1.ResourceCPU]
	fmt.Println("node cpu je: ")
	fmt.Println(nodeCpu.String())
	fmt.Println(nodeCpu)
	fmt.Println("NODE NAME JE: ")
	fmt.Println(nodeInfo.Node().Name)
	fmt.Println("-----------------INFORMACIJE O ČVORU END---------------")
	fmt.Println()
	fmt.Println()

	fmt.Println("-----------------INFORMACIJE O PODU-----------------")
	podLabelValue := pod.Labels["scheduleon"]
	var podCPU resource.Quantity

	fmt.Println("LABELA PODA:")
	fmt.Println(podLabelValue)

	fmt.Println("-----------------INFORMACIJE O PODU END-----------------")
	fmt.Println()
	fmt.Println()

	fmt.Println("-----------------IZVRŠAVANJE KODA START-----------------")
	if podLabelValue == "agent2node" {

		fmt.Println("---agent2node je true")
		//ako se radi o cvoru na koji trebamo schedulat, provjeri njegove resurse
		for _, container := range pod.Spec.Containers {
			fmt.Println("------u for petlji za resurse: ")
			if cpu, ok := container.Resources.Requests[v1.ResourceCPU]; ok {
				podCPU.Add(cpu)
				fmt.Println("------Printam cpu poda:")
				fmt.Println(cpu)
			} else {
				fmt.Println("------nista ")
			}
		}

		//ako su resursi zadovoljavajuci stavi pod na taj node
		if nodeCpu.Cmp(podCPU) > 0 {
			fmt.Println("---ovdje da vratim success")
			fmt.Println(nodeCpu.Cmp(podCPU))
			return framework.NewStatus(framework.Success)
		} else {
			//ako nisu, onda logika za ping i trazi najblizi cvor drugi
			fmt.Println("---nedovoljno resursa, trazim drugi node ....")
			return framework.NewStatus(framework.Unschedulable, "nije moguce schedulat, trazi drugi cvor")
		}
	}

	fmt.Println("-----------------IZVRŠAVANJE KODA END-----------------")
	fmt.Println()
	fmt.Println()

	return framework.NewStatus(framework.Unschedulable, "Node is not masternode")
}

//func (s *customFilterPlugin) PreBind(ctx context.Context, pod *v1.Pod, nodeName string) *framework.Status {
//
//	return framework.NewStatus(framework.Success, "")
//}

func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &customFilterPlugin{}, nil
}
