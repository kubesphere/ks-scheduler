package prioritize

import (
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

var (
	PipelinePriority = Prioritize{
		Name: "pipelinePriority",
		Func: Pipeline,
	}
)

func Pipeline(pod v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error) {
	var priorityList schedulerapi.HostPriorityList
	priorityList = make([]schedulerapi.HostPriority, len(nodes))
	keys := parseMark(pod.Labels)

	for i, node := range nodes {
		score, err := Calculation(keys, node)
		if err != nil {
			panic(err)
		}

		priorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: score,
		}
	}
	return &priorityList, nil
}

func parseMark(labels map[string]string) ([]string) {}

func Calculation(keys []string, node v1.Node) (score int, err error) {}
