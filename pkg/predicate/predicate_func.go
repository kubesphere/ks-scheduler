package predicate

import (
	log "github.com/golang/glog"
	"k8s.io/api/core/v1"
)

func AlwaysTrue(pod v1.Pod, node v1.Node) (bool, error) {
	log.Info("podddd")
	log.Info(pod)
	return true, nil
}
