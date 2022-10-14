package computeintensive

import (
	"Edge-Scheduler/pkg/util"
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

func (ci *ComputeIntensive) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := ci.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	                                                                                                                                                                                                              
	if err != nil || nodeInfo.Node() == nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v, node is nil: %v", nodeName, err, nodeInfo.Node() == nil))
	}

	allocatableCPU := nodeInfo.Allocatable.MilliCPU
	turboFrequency := nodeInfo.Node().Labels["frequency"]
	cache := nodeInfo.Node().Labels["cache"]

	klog.Infof("Calculate computing score for pod %s in namespace %s on node %s successfully!", p.Name, p.Namespace, nodeName)
	return (allocatableCPU * CoreWeight + (util.StrToInt64(turboFrequency) / 1000) * FrequencyWeight + (util.StrToInt64(cache) / 1024) * CacheWeight), framework.NewStatus(framework.Success, "")
}

func (ci *ComputeIntensive) NormalizeScore(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
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
		klog.Infof("Node: %v, Score: %v in Plugin: Mandalorian When scheduling Pod: %v/%v", scores[i].Name, scores[i].Score, pod.GetNamespace(), pod.GetName())
	}
	return nil
}

func (ci *ComputeIntensive) ScoreExtensions() framework.ScoreExtensions {
	return ci
}
