package dataintensive

import (
	"Edge-Scheduler/pkg/util"
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

func (di *DataIntensive) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := di.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil || nodeInfo.Node() == nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v, node is nil: %v", nodeName, err, nodeInfo.Node() == nil))
	}

	iopsLabel := nodeInfo.Node().Labels["iops"]
	iopsScore := int64(0)
	
	if strings.Contains(iopsLabel, "kB_s") {
		iopsScore = int64(1)
	} else if strings.Contains(iopsLabel, "MB_s") {
		iopsScore = util.StrToInt64(util.Split(iopsLabel))
		t := iopsScore / 10
		if t < 1 {	// <=10MB/s
			iopsScore = 10
		} else if t > 100 {	// >=1000MB/s
			iopsScore = 100
		} else {  // 10MB/s ~ 1000MB/s
			if t < 20 { // 10MB/s ~ 200MB/s
				t = 20
			}
			iopsScore = t			
		}
	} else if strings.Contains(iopsLabel, "GB_s") {
		iopsScore = int64(100)
	} else {
		iopsScore = int64(0)
	}

	// switch {
	// case strings.Contains(iopsLabel, "kB_s"):
	// 	iopsScore = int64(1)
	// case strings.Contains(iopsLabel, "MB_s"):
	// 	iopsScore = util.StrToInt64(util.Split(iopsLabel))
	// 	t := iopsScore / 10
	// 	if t < 1 {	// <=10MB/s
	// 		iopsScore = 10
	// 	} else if t > 100 {	// >=1000MB/s
	// 		iopsScore = 100
	// 	} else {  // 10MB/s ~ 1000MB/s
	// 		if t < 20 { // 10MB/s ~ 200MB/s
	// 			t = 20
	// 		}
	// 		iopsScore = t			
	// 	}
    // case strings.Contains(iopsLabel, "GB_s"):
	// 	iopsScore = int64(100)
	// default:
	// 	iopsScore = int64(0)
	// }

	throughputLabel := nodeInfo.Node().Labels["throughput"]
	throughputScore := int64(0)
	if strings.Contains(throughputLabel, "kB_s") {
		throughputScore = int64(0)
	} else if strings.Contains(throughputLabel, "MB_s") {
		throughputScore = util.StrToInt64(util.Split(throughputLabel))
		t := throughputScore / 10
		if t < 1 {	// <=10MB/s
			throughputScore = 1
		} else if t > 100 {	// >=1000MB/s
			throughputScore = 100
		} else {  // 10MB/s ~ 1000MB/s
			throughputScore = t	
		}
	} else if strings.Contains(throughputLabel, "GB_s") {
		throughputScore = int64(100)
	} else {
		throughputScore = int64(0)
	}
	klog.Infof("Calculate data score for pod %s in namespace %s on node %s successfully!", p.Name, p.Namespace, nodeName)
	return (iopsScore + throughputScore), framework.NewStatus(framework.Success, "")
}

func (di *DataIntensive) NormalizeScore(ctx context.Context, cycleState *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
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

func (di *DataIntensive) ScoreExtensions() framework.ScoreExtensions {
	return di
}
