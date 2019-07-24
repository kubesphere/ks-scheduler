package prioritize

import (
	"k8s.io/api/core/v1"
	"github.com/soulseen/ks-schduler/pkg/sqlite"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

type Prioritize struct {
	Name string
	Func func(pod v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error)
	Datatable sqlite.KeyNodeTable
}

func (p Prioritize) Handler(args schedulerapi.ExtenderArgs) (*schedulerapi.HostPriorityList, error) {
	return p.Func(*args.Pod, args.Nodes.Items)
}
