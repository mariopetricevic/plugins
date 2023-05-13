package plugins

import (
	"fmt"
	"context"
	"time"
	"math"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"github.com/go-ping/ping"
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

func (p *customFilterPlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	
	fmt.Println("-----------------INFORMACIJE O ČVORU-----------------")
	//resursi cvora
	nodeCpu := nodeInfo.Node().Status.Capacity[v1.ResourceCPU]
	fmt.Println("node cpu je: ")
	fmt.Println(nodeCpu.String())
	fmt.Println(nodeCpu)
	fmt.Println("NODE NAME JE: ")
	fmt.Println(nodeInfo.Node().Name)
	fmt.Println("ADRESA JE: ")
	fmt.Println(nodeInfo.Node().Status.Addresses[0].Address)
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
	if podLabelValue == "agent2node" && nodeInfo.Node().Name == "agent2node" {

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

		smh := false
		//ako su resursi zadovoljavajuci stavi pod na taj node
		if nodeCpu.Cmp(podCPU) > 0 && smh{
			fmt.Println("---ovdje da vratim success")
			fmt.Println(nodeCpu.Cmp(podCPU))
			return framework.NewStatus(framework.Success)
		} else {
			//ako nisu, onda logika za ping i trazi najblizi cvor drugi
			fmt.Println("---nedovoljno resursa, trazim drugi node ....")
			nodes, err := p.handle.SnapshotSharedLister().NodeInfos().List()
			if err != nil {
				return framework.NewStatus(framework.Error, "Error getting node list")
			}
			fmt.Println("---pronaso informacije o nodovima")
			var closestNode *framework.NodeInfo
			minRTT := time.Duration(math.MaxInt64)
			for _, node := range nodes {
				rtt, err := pingNode(node.Node().Status.Addresses[0].Address)
				if err != nil {
					continue
				}
				if rtt < minRTT {
					minRTT = rtt
					closestNode = node
				}
			}
			
			
			// Ako je najbliži čvor trenutni čvor, vrati Success
			if closestNode != nil && closestNode.Node().Name == nodeInfo.Node().Name {
				fmt.Println("---pronaden najblizi cvor pomocu pinga")
				fmt.Println("-----INFORMACIJE O NAJBLIZEM ČVORU-----------------")
				//resursi cvora
				nodeCpu := closestNode.Node().Status.Capacity[v1.ResourceCPU]
				fmt.Println("node cpu je: ")
				fmt.Println(nodeCpu.String())
				fmt.Println(nodeCpu)
				fmt.Println("NODE NAME JE: ")
				fmt.Println(closestNode.Node().Name)
				fmt.Println("ADRESA JE: ")
				fmt.Println(closestNode.Node().Status.Addresses[0].Address)
				fmt.Println("-----INFORMACIJE O NAJBLIZEM ČVORU END---------------")
				return framework.NewStatus(framework.Success)
			} else {
				return framework.NewStatus(framework.Unschedulable, "nije moguce schedulat, trazi drugi cvor")
			}
			
			
			
			
			//return framework.NewStatus(framework.Unschedulable, "nije moguce schedulat, trazi drugi cvor")
		}
	}

	fmt.Println("-----------------IZVRŠAVANJE KODA END-----------------")
	fmt.Println()
	fmt.Println()

	return framework.NewStatus(framework.Unschedulable, "Node is not masternode")
}

func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &customFilterPlugin{}, nil
}

func pingNode(ip string) (time.Duration, error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return 0, err
	}
	pinger.Count = 3
	pinger.Timeout = time.Second * 5
	pinger.SetPrivileged(true)
	pinger.Run()
	stats := pinger.Statistics()
	return stats.AvgRtt, nil
}
