package controller

import (
	"fmt"
	"github.com/golang/glog"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/staging/src/k8s.io/client-go/util/workqueue"
	"k8s.io/sample-controller/pkg/signals"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"time"
)

const (
	// maxRetries is the number of times a service will be retried before it is dropped out of the queue.
	// With the current rate-limiter in use (5ms*2^(maxRetries-1)) the following numbers represent the
	// sequence of delays between successive queuings of a service.
	//
	// 5ms, 10ms, 20ms, 40ms, 80ms, 160ms, 320ms, 640ms, 1.3s, 2.6s, 5.1s, 10.2s, 20.4s, 41s, 82s
	maxRetries = 15

	defaultResync = 600 * time.Second

	DefaultLabelKey   = "scheduler"
	DefaultLabelValue = "ks-scheduler"
)

type Controller struct {
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
	indexer  cache.Indexer
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) {
	// catch error
	defer utilruntime.HandleCrash()
	// close queue and shutdown work.
	defer c.queue.ShutDown()

	glog.Info("start controller...")

	// run Informer
	go c.informer.Run(stopCh)

	// waiting for cache successful
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("cache time out."))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	glog.Info("Stopping Pod controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {

	// get next.
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	// remove Key
	defer c.queue.Done(key)

	// process Key
	err := c.processItem(key.(string))
	c.handleErr(err, key)

	return true
}

func (c *Controller) processItem(key string) error {
	glog.Infof("Start procrss event %s", key)
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		// Below we will warm up our cache with a Pod, so that we will see a delete for one pod
		fmt.Printf("Pod %s does not exist anymore\n", key)
	} else {
		// Note that you also have to check the uid if you have a local controlled resource, which
		// is dependent on the actual instance, to detect that a Pod was recreated with the same name
		fmt.Printf("Sync/Add/Update for Pod %s\n", obj.(*apiv1.Pod).GetName())
	}
	return nil
}

func RunController() {
	// Get a config to talk to the apiserver
	glog.Info("setting up client for manager")

	cfg, err := config.GetConfig()
	if err != nil {
		glog.Error(err, "unable to set up client config")
		os.Exit(1)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("build clientset error: %s", err.Error())
	}

	//TODO
	//setLabel := map[string]string{"scheduler": "ks-scheduler"}
	//
	//selectPod := fields.SelectorFromSet(setLabel)

	podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", apiv1.NamespaceDefault, fields.Everything())

	stopCh := signals.SetupSignalHandler()

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	indexer, informer := cache.NewIndexerInformer(podListWatcher, &apiv1.Pod{}, defaultResync, cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool {
			switch t := obj.(type) {
			case *apiv1.Pod:
				return assignedPod(t)
			default:
				utilruntime.HandleError(fmt.Errorf("unable to handle object : %T", obj))
				return false
			}
		},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				switch t := obj.(type) {
				case *apiv1.Pod:
					res := assignedPod(t)
					if res == true {
						key, err := cache.MetaNamespaceKeyFunc(obj)
						if err == nil {
							queue.Add(key)
						}
					}
				default:
					utilruntime.HandleError(fmt.Errorf("unable to handle object : %T", obj))
					return
				}
			},
		},
	}, cache.Indexers{})

	ctrl := Controller{
		queue,
		informer,
		indexer,
	}

	//ctrl.Run(stopCh)
	go ctrl.Run(3, stopCh)

	// Wait forever
	select {}

}

// assignedPod selects pods that are assigned (scheduled and running).
func assignedPod(pod *apiv1.Pod) bool {
	if len(pod.Spec.NodeName) == 0 {
		glog.Info("not have binding node")
		return false
	}
	if pod.Status.Phase != apiv1.PodRunning {
		glog.Info("Status not Running")
		return false
	}
	if pod.Labels[DefaultLabelKey] != DefaultLabelValue {
		glog.Info("no scheduler label")
		return false
	}
	return true
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < maxRetries {
		glog.Info("Error syncing pod %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	glog.Infof("Dropping pod %q out of the queue: %v", key, err)
}
