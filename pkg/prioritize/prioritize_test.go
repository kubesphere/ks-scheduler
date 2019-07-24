package prioritize


//import (
//"path/filepath"
//"reflect"
//"testing"
//"k8s.io/api/core/v1"
//schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
//)
//
//func TestVolumeListSet(t *testing.T) {
//	pods := []struct {
//		Input    v1.Pod
//		Expected schedulerapi.HostPriorityList
//	} {
//		{v1.Pod{{"Pod","v1"},{"test1",""}},schedulerapi.HostPriorityList{}},
//	}
//	for _, test := range pods {
//		if len(test.Expected) != 0 {
//			test.Expected[0].Source = filepath.FromSlash(test.Expected[0].Source)
//		}
//		got := VolumeList{}
//		got.Set(test.Input)
//		if !reflect.DeepEqual(got, test.Expected) {
//			t.Errorf("On test %s, got %#v, expected %#v", test.Input, got, test.Expected)
//		}
//	}
//}