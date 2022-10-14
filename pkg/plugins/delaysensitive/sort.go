package delaysensitive

import "k8s.io/kubernetes/pkg/scheduler/framework"

func (ds *DelaySensitive) Less(podInfo1, podInfo2 *framework.QueuedPodInfo) bool {
	return Less(podInfo1, podInfo2)
}

func Less(podInfo1, podInfo2 *framework.QueuedPodInfo) bool {
	return (SetPodPriority(podInfo1) > SetPodPriority(podInfo2)) || (SetPodPriority(podInfo1) > SetPodPriority(podInfo2) && podInfo1.Timestamp.Before(podInfo2.Timestamp))
}

func SetPodPriority(podInfo *framework.QueuedPodInfo) int {
	if podType, ok := podInfo.Pod.Labels["type"]; ok {
		switch podType {
		case "client":
			return 3
		case "logic":
			return 2
		case "storage":
			return 1
		default:
			return 0
		}
	}
	return 0
}
