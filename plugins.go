package plugins

import (
	"context"
	"fmt"
	"time"

	"github.com/go-ping/ping"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

type customFilterPlugin struct {
	handle framework.Handle
}

const (
	// Name : name of plugin used in the plugin registry and configurations.
	Name = "CustomFilterPlugin"
)

// var _ framework.FilterPlugin = &customFilterPlugin{}
var _ = framework.FilterPlugin(&customFilterPlugin{})

var _ = framework.ScorePlugin(&customFilterPlugin{})

func (p *customFilterPlugin) Name() string {
	return Name
}

func (p *customFilterPlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {

	fmt.Println("------------------FILTER-------------------------------------------")

	fmt.Println("-----------------INFORMACIJE O ČVORU-----------------")
	//resursi cvora
	nodeCpu := nodeInfo.Node().Status.Capacity[v1.ResourceCPU]
	fmt.Println("---node cpu, cpu, name, adresa: ")
	fmt.Println("----", nodeCpu.String(), nodeCpu, nodeInfo.Node().Name, nodeInfo.Node().Status.Addresses[0].Address)
	fmt.Println("-----------------INFORMACIJE O ČVORU END---------------")
	fmt.Println()
	fmt.Println()

	fmt.Println("-----------------INFORMACIJE O PODU-----------------")
	podLabelValue := pod.Labels["scheduleon"]
	var podCPU resource.Quantity
	fmt.Println("---LABELA PODA:", podLabelValue)
	fmt.Println("-----------------INFORMACIJE O PODU END-----------------")
	fmt.Println()
	fmt.Println()

	fmt.Println("-----------------IZVRŠAVANJE KODA START-----------------")

	for _, container := range pod.Spec.Containers {
		fmt.Println("------u for petlji za resurse: ")
		if cpu, ok := container.Resources.Requests[v1.ResourceCPU]; ok {
			podCPU.Add(cpu)
			fmt.Println("------Printam cpu poda:", cpu)
		} else {
			fmt.Println("------nista ")
		}
	}

	//smh := false
	//ako su resursi zadovoljavajuci stavi pod na taj node
	if nodeCpu.Cmp(podCPU) > 0 { //smh dio koda treba maknuti jer trenutno ne radi ovo s provjerom resursa
		fmt.Println("---DOVOLJNO RESURSA", nodeCpu.Cmp(podCPU))
		fmt.Println("------------------FILTER END-------------------------------------------")
		return framework.NewStatus(framework.Success)
	} else {
		fmt.Println("---NEDOVOLJNO RESURSA")
		fmt.Println("------------------FILTER END-------------------------------------------")
		return framework.NewStatus(framework.Unschedulable, "NEDOVOLJNO RESURSA")
	}

	// if podLabelValue == "agent2node" && nodeInfo.Node().Name == "agent2node" {

	// 	fmt.Println("---agent2node je true")
	// 	//ako se radi o cvoru na koji trebamo schedulat, provjeri njegove resurse
	// 	// for _, container := range pod.Spec.Containers {
	// 	// 	fmt.Println("------u for petlji za resurse: ")
	// 	// 	if cpu, ok := container.Resources.Requests[v1.ResourceCPU]; ok {
	// 	// 		podCPU.Add(cpu)
	// 	// 		fmt.Println("------Printam cpu poda:")
	// 	// 		fmt.Println(cpu)
	// 	// 	} else {
	// 	// 		fmt.Println("------nista ")
	// 	// 	}
	// 	// }

	// 	smh := false
	// 	//ako su resursi zadovoljavajuci stavi pod na taj node
	// 	if nodeCpu.Cmp(podCPU) > 0 && smh { //smh dio koda treba maknuti jer trenutno ne radi ovo s provjerom resursa
	// 		fmt.Println("---ovdje da vratim success")
	// 		fmt.Println(nodeCpu.Cmp(podCPU))
	// 		return framework.NewStatus(framework.Success)
	// 	} else {
	// 		//ako nisu, onda logika za ping i trazi najblizi cvor drugi
	// 		fmt.Println("---nedovoljno resursa, trazim drugi node ....")

	// 		if p.handle == nil {
	// 			fmt.Println("---handle je nulll....")
	// 		}

	// 		nodeLister := p.handle.SnapshotSharedLister().NodeInfos()
	// 		if nodeLister == nil {
	// 			fmt.Println("---nodes lister je nullll....")
	// 		}

	// 		nodes, err := p.handle.SnapshotSharedLister().NodeInfos().List()
	// 		if err != nil {
	// 			fmt.Println("---ovo tu je null sta li...")
	// 			return framework.NewStatus(framework.Error, "Error getting node list")
	// 		}
	// 		fmt.Println("---pronaso informacije o nodovima. Ispis nodova:")

	// 		for _, node := range nodes {
	// 			fmt.Println(node.Node().Name)
	// 			fmt.Println(node.Node().Status.Addresses[0].Address)
	// 		}
	// 		fmt.Println("gotov ispis nodova")

	// 		// var closestNode *framework.NodeInfo
	// 		// minRTT := time.Duration(math.MaxInt64)
	// 		// for _, node := range nodes {
	// 		// 	rtt, err := pingNode(node.Node().Status.Addresses[0].Address)
	// 		// 	fmt.Println("ispis rtta")
	// 		// 	fmt.Println(node.Node().Name)
	// 		// 	fmt.Println(rtt)
	// 		// 	if err != nil {
	// 		// 		fmt.Println("error dohvcanja rtta")
	// 		// 		continue
	// 		// 	}
	// 		// 	if rtt < minRTT {
	// 		// 		fmt.Println("ispis rtta")
	// 		// 		fmt.Println(rtt)
	// 		// 		minRTT = rtt
	// 		// 		closestNode = node
	// 		// 	}
	// 		// 	fmt.Println("proso po svima")
	// 		// }

	// 		// fmt.Println("-----INFORMACIJE O NAJBLIZEM ČVORU-----------------")
	// 		// //resursi cvora
	// 		// nodeCpu := closestNode.Node().Status.Capacity[v1.ResourceCPU]
	// 		// fmt.Println("node cpu je: ")
	// 		// fmt.Println(nodeCpu.String())
	// 		// fmt.Println(nodeCpu)
	// 		// fmt.Println("NODE NAME JE: ")
	// 		// fmt.Println(closestNode.Node().Name)
	// 		// fmt.Println("ADRESA JE: ")
	// 		// fmt.Println(closestNode.Node().Status.Addresses[0].Address)
	// 		// fmt.Println("-----INFORMACIJE O NAJBLIZEM ČVORU END---------------")

	// 		//&& closestNode.Node().Name == nodeInfo.Node().Name
	// 		// Ako je najbliži čvor trenutni čvor, vrati Success
	// 		// if closestNode != nil {
	// 		// 	fmt.Println("---pronaden najblizi cvor pomocu pinga")
	// 		// 	fmt.Println("-----INFORMACIJE O NAJBLIZEM ČVORU-----------------")
	// 		// 	//resursi cvora
	// 		// 	nodeCpu := closestNode.Node().Status.Capacity[v1.ResourceCPU]
	// 		// 	fmt.Println("node cpu je: ")
	// 		// 	fmt.Println(nodeCpu.String())
	// 		// 	fmt.Println(nodeCpu)
	// 		// 	fmt.Println("NODE NAME JE: ")
	// 		// 	fmt.Println(closestNode.Node().Name)
	// 		// 	fmt.Println("ADRESA JE: ")
	// 		// 	fmt.Println(closestNode.Node().Status.Addresses[0].Address)
	// 		// 	fmt.Println("-----INFORMACIJE O NAJBLIZEM ČVORU END---------------")
	// 		// 	return framework.NewStatus(framework.Success)
	// 		// } else {
	// 		// 	fmt.Println("nemoguce schedulatt")
	// 		// 	return framework.NewStatus(framework.Unschedulable, "nije moguce schedulat, trazi drugi cvor")
	// 		// }

	// 		//return framework.NewStatus(framework.Unschedulable, "nije moguce schedulat, trazi drugi cvor")
	// 	}
	// }

	// fmt.Println("-----------------IZVRŠAVANJE KODA END-----------------")
	// fmt.Println()
	// fmt.Println()
	// fmt.Println("------------------FILTER  END-------------------------------------------")
	// return framework.NewStatus(framework.Unschedulable, "Node is not schedulable")
}

func (p *customFilterPlugin) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {

	fmt.Println("--------------------SCORE-------------------------------------")
	fmt.Println()

	if p.handle == nil {
		fmt.Println("---handle je nulll....")
	}

	nodes, err := p.handle.SnapshotSharedLister().NodeInfos().List()
	if err != nil {
		fmt.Println("---ovo tu je null sta li...")
	}
	fmt.Println("---pronaso informacije o nodovima. Ispis nodova:")

	for _, node := range nodes {
		fmt.Println("---", node.Node().Name)
		fmt.Println("---", node.Node().Status.Addresses[0].Address)
	}
	fmt.Println("---gotov ispis nodova")

	podLabelValue := pod.Labels["scheduleon"]

	fmt.Println()
	fmt.Println()

	if nodeName == podLabelValue {
		fmt.Println("---ispis scorea za node")
		fmt.Println("---", nodeName)
		fmt.Println("--------------------SCORE END-------------------------------------")
		return 100, nil
	}

	nodeInfo, err := p.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		fmt.Println("error")
		fmt.Println("--------------------SCORE END-------------------------------------")
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("Error getting node %s from Snapshot: %v", nodeName, err))
	}

	rtt, err := pingNode(nodeInfo.Node().Status.Addresses[0].Address)
	fmt.Println("---ISPIS RTTA:")
	fmt.Println("----", nodeInfo.Node().Name)
	fmt.Println("----", rtt)
	fmt.Println("---ISPIS RTTA END")
	if err != nil {
		fmt.Println("error")
		fmt.Println("--------------------SCORE END-------------------------------------")
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("Error pinging node %s: %v", nodeName, err))
	}

	// Ovdje bi trebalo pretvoriti vrijeme odziva u ocjenu. Niže vrijeme odziva bi trebalo rezultirati višom ocjenom.
	// Ovaj primjer samo pretvara vrijeme odziva u milisekunde i oduzima ga od nekog maksimalnog broja (na primjer, 10000)
	// kako bi niže vrijeme odziva rezultiralo višom ocjenom. Trebali biste prilagoditi ovu logiku prema svojim potrebama.
	score := 95 - int64(rtt.Milliseconds())
	fmt.Println("---ispis scorea za node:")
	fmt.Println("----", nodeInfo.Node().Name)
	fmt.Println("----", score)

	fmt.Println("--------------------SCORE END-------------------------------------")
	return score, nil
}

func (p *customFilterPlugin) ScoreExtensions() framework.ScoreExtensions {
	return p
}

func (p *customFilterPlugin) NormalizeScore(_ context.Context, _ *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	fmt.Println("----------NORMALIZE  SCORE -------------------------------------")
	var (
		highest int64 = 0
		lowest        = scores[0].Score
	)

	for _, nodeScore := range scores {
		if nodeScore.Score < lowest {
			lowest = nodeScore.Score
		}
		if nodeScore.Score > highest {
			highest = nodeScore.Score
		}
	}

	if highest == lowest {
		lowest--
	}

	// Set Range to [0-100]
	for i, nodeScore := range scores {
		scores[i].Score = (nodeScore.Score - lowest) * framework.MaxNodeScore / (highest - lowest)
		fmt.Println(scores[i].Name, scores[i].Score, pod.GetNamespace(), pod.GetName())
	}

	fmt.Println("----------NORMALIZE  SCORE END-------------------------------------")
	return nil
}

func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &customFilterPlugin{handle: handle}, nil
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
