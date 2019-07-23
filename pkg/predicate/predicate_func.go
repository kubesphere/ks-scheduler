package predicate

import "k8s.io/api/core/v1"

var (
	TruePredicate = Predicate{
		Name: "alwaysTrue",
		Func: func(pod v1.Pod, node v1.Node) (bool, error) {
			return true, nil
		},
	}
)
