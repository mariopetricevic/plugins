package plugins

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

type MyK3SPlugin struct {
	handle framework.Handle
}

type NodeNameDistance struct {
	Name     string
	Distance int
}

const (
	// Name : name of plugin used in the plugin registry and configurations.
	Name = "MyK3SPlugin"
)

var _ = framework.FilterPlugin(&MyK3SPlugin{})
var _ = framework.ScorePlugin(&MyK3SPlugin{})

func (p *MyK3SPlugin) Name() string {
	return Name
}

func (p *MyK3SPlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	//
	// Filter nodes that can run the pod. If node cannot run the pod
	// method returns framework status Unschedulable.
	// Additionaly, method filters nodes that already have pod running
	// with the same application name.
	//
	fmt.Println("filter pod:", pod.Name, ", application name: ", pod.Labels["applicationName"], ", Node: ", nodeInfo.Node().Name)
	fmt.Println("Node resources: CPU:", nodeInfo.Node().Status.Allocatable.Cpu(), ", memory: ", nodeInfo.Node().Status.Allocatable.Memory())

	//availableResources := nodeInfo.Node().Status.Allocatable
	nodeCpu := nodeInfo.Node().Status.Capacity[v1.ResourceCPU]

	//calculate allocated CPU for running all containters in current pod
	var podCPU resource.Quantity
	for _, container := range pod.Spec.Containers {
		if cpu, ok := container.Resources.Requests[v1.ResourceCPU]; ok {
			podCPU.Add(cpu)
		}
	}

	//If requred resources for running pod are less than available resources on a node
	if nodeCpu.Cmp(podCPU) > 0 {

		// Check if node already runs application with the same name
		pods := nodeInfo.Pods
		for _, p := range pods {
			labels := p.Pod.GetLabels()
			//if any pod running on current node already has pod with the label applicationName equal to pod given as parameter, then this app
			// already exists on this node.
			if applicationName, found := labels["applicationName"]; found && applicationName == pod.Labels["applicationName"] {

				fmt.Println("Application:", pod.Labels["applicationName"], "Already exists on node: ", nodeInfo.Node().Name)

				return framework.NewStatus(framework.Unschedulable, "Application already exists on this node")
			}
		}
		return framework.NewStatus(framework.Success)
	} else {
		return framework.NewStatus(framework.Unschedulable, "Not enough resources to run application: ", pod.Labels["applicationName"])
	}
}

func (p *MyK3SPlugin) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {

	fmt.Println("--------------------SCORE-------------------------------------")
	fmt.Println()

	//deployments need to have label requestFrom which indicates on which node request for pod has came
	requestFromNode := pod.Labels["requestFrom"] //scheduleon

	nodes, err := p.handle.SnapshotSharedLister().NodeInfos().List()
	if err != nil {
		fmt.Println("Error occured while getting all nodes in cluster.")
	}

	for _, node := range nodes {
		fmt.Println("NodeName :", node.Node().Name, "Address: ", node.Node().Status.Addresses[0].Address)
	}

	fmt.Println("Calculating score for current node:", nodeName)

	//Get node from which request came
	requesterNodeInfo, err := p.handle.SnapshotSharedLister().NodeInfos().Get(requestFromNode)
	requesterNode := requesterNodeInfo.Node()
	requesterNodeLabels := requesterNode.GetLabels()

	//map of node names and distances
	mapNodeNameDistance := make(map[string]int)

	// get labels with node names in cluster and their distances
	for label, value := range requesterNodeLabels {
		if strings.HasPrefix(label, "ping-") {

			distanceToNode, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Error while converting.")
			}

			nodeName := strings.TrimPrefix(label, "ping-")
			mapNodeNameDistance[nodeName] = distanceToNode
		}
	}

	// Store node distances
	var sortedNodeDistances []NodeNameDistance
	for label, value := range mapNodeNameDistance {
		sortedNodeDistances = append(sortedNodeDistances, NodeNameDistance{label, value})
	}

	//sort node distances from lowest to highest
	sort.Slice(sortedNodeDistances, func(i, j int) bool {
		return sortedNodeDistances[i].Distance < sortedNodeDistances[j].Distance
	})

	fmt.Println("----------Sorted node distances-------------")
	for _, nodeNameDistance := range sortedNodeDistances {
		fmt.Printf("Label: %s, Value: %d\n", nodeNameDistance.Name, nodeNameDistance.Distance)
	}
	fmt.Println("-------End sorted node distances-------------")

	//Scoring nodes
	var score int = 0
	if requestFromNode == nodeName {
		score = 100
		fmt.Println("Requested node: ", requestFromNode, " can run the pod. Scoring node with maximum score.")
		return 100, nil
	} else {
		for _, nodeNameDistance := range sortedNodeDistances {
			//find first (with lowest distance) node that can run a pod
			if nodeNameDistance.Name == nodeName {

				score = 90 - nodeNameDistance.Distance
				fmt.Println("Scoring node: ", nodeNameDistance.Name, " score: ", score)
				break
			}
		}
	}

	fmt.Println("--------------------SCORE END-------------------------------------")
	return int64(score), nil
}

func (p *MyK3SPlugin) ScoreExtensions() framework.ScoreExtensions {
	return p
}

func (p *MyK3SPlugin) NormalizeScore(_ context.Context, _ *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
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
	return &MyK3SPlugin{handle: handle}, nil
}
