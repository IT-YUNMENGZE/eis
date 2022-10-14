package delaysensitive

import (
	"Edge-Scheduler/pkg/util"
	"context"
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	len "my.domain/Len/api/v1"
)

func (ds *DelaySensitive) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	// Get Node Info of the given node name
	nodeInfo, err := ds.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v", nodeName, err))
	}

	// Get the list of NodeInfos
	nodeInfosList, err := ds.handle.SnapshotSharedLister().NodeInfos().List()
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("%q: getting the list of NodeInfos from Snapshot: %v", nodeName, err))
	}

	// Get Len Info
	currentLen := &len.Len{}
	err = ds.cache.Get(ctx, types.NamespacedName{Name: nodeName}, currentLen)
	if err != nil {
		klog.Errorf("Get Len Error: %v", err)
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("Score Node: %v Error: %v", nodeInfo.Node().Name, err))
	}

	uNodeScore, err := CalculateScore(currentLen, p, nodeInfosList)
	if err != nil {
		klog.Errorf("CalculateScore Error: %v", err)
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("Score Node: %v Error: %v", nodeInfo.Node().Name, err))
	}
	// uint64 => int64
	nodeScore := util.Uint64ToInt64(uNodeScore)

	klog.Infof("Calculate latency score for pod %s in namespace %s on node %s successfully!", p.Name, p.Namespace, nodeName)
	return nodeScore, framework.NewStatus(framework.Success, "")
}

func (ds *DelaySensitive) NormalizeScore(_ context.Context, _ *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
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
		klog.Infof("Node: %v, Score: %v in Plugin: DelaySensitive When scheduling Pod: %v/%v", scores[i].Name, scores[i].Score, pod.GetNamespace(), pod.GetName())
	}
	return nil
}

func (ds *DelaySensitive) ScoreExtensions() framework.ScoreExtensions {
	return ds
}

func CalculateScore(l *len.Len, pod *v1.Pod, nodeInfosList []*framework.NodeInfo) (uint64, error) {
	if appType, ok := pod.GetLabels()["type"]; ok {
		switch appType {
		case "client":
			return CalculateClientScore(l), nil
		case "logic":
			if appName, ok := pod.GetLabels()["app"]; ok {
				return CalculateLogicScore(l, pod, nodeInfosList, appName), nil
			}
			return 0, errors.New("pod app label is not set")
		case "storage":
			if appName, ok := pod.GetLabels()["app"]; ok {
				return CalculateStorageScore(l, pod, nodeInfosList, appName), nil
			}
			return 0, errors.New("pod app label is not set")
		default:
			return 0, errors.New("the pod type label is illegal")
		}
	}
	return 0, errors.New("pod type label is not set")
}

func CalculateClientScore(l *len.Len) uint64 {
	for _, lantency := range l.Status.LatencyList {
		if lantency.NodeName == "gateway" {
			gatewayLantency := lantency.Latency
			return (1 / uint64(gatewayLantency)) * 100
		}
	}
	return 0
}

func CalculateLogicScore(l *len.Len, pod *v1.Pod, nodeInfosList []*framework.NodeInfo, appName string) uint64 {
	for _, nodeInfo := range nodeInfosList {
		for _, pod := range nodeInfo.Pods {
			if podAppLabel, ok1 := pod.Pod.GetLabels()["app"]; ok1 {
				if podTypeLabel, ok2 := pod.Pod.GetLabels()["type"]; ok2 {
					if podAppLabel == appName && podTypeLabel == "client" {
						clientPodNode := pod.Pod.Spec.NodeName
						for _, lantency := range l.Status.LatencyList {
							if lantency.NodeName == clientPodNode {
								clientLantency := lantency.Latency
								return (1 / uint64(clientLantency)) * 100
							}
						}
					}
				}
			}
		}
	}
	return 0
}

func CalculateStorageScore(l *len.Len, pod *v1.Pod, nodeInfosList []*framework.NodeInfo, appName string) uint64 {
	for _, nodeInfo := range nodeInfosList {
		for _, pod := range nodeInfo.Pods {
			if podAppLabel, ok1 := pod.Pod.GetLabels()["app"]; ok1 {
				if podTypeLabel, ok2 := pod.Pod.GetLabels()["type"]; ok2 {
					if podAppLabel == appName && podTypeLabel == "logic" {
						logicPodNode := pod.Pod.Spec.NodeName
						for _, lantency := range l.Status.LatencyList {
							if lantency.NodeName == logicPodNode {
								logicLantency := lantency.Latency
								return (1 / uint64(logicLantency)) * 100
							}
						}
					}
				}
			}

		}
	}
	return 0
}
