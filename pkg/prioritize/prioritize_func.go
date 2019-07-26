package prioritize

import (
	"github.com/soulseen/ks-pipeline-schduler/pkg/sqlite"
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

var (
	PipelinePriority = Prioritize{
		Name:      "pipeline",
		Func:      Pipeline,
		Datatable: sqlite.InitKeyNodeTable(),
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

func parseMark(labels map[string]string) ([]string) {

	var test []string
	test = append(test, "aa")
	return test

}

func Calculation(keys []string, node v1.Node) (score int, err error) {

	score = 1
	return score, nil
}
