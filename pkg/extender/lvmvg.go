package extender

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func getLVMNodeToAllocatableVolumeMap() map[string]int64 {
	lvsMap := make(map[string]int64)
	lvmClient, err := getLVMClient()
	if err != nil {
		log.Print("could not initialize lvm client : ", err)
		return lvsMap
	}
	lvs, err := lvmClient.LocalV1alpha1().LVMNodes("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Print("could not get lvm nodes using lvm client")
		return lvsMap
	}
	for _, lv := range lvs.Items {
		for _, vg := range lv.VolumeGroups {
			if vg.Name == "lvmvg" {
				freeSpace, err := extractIntegralPartFromVolumeInBytes(ConvertToIBytes(vg.Free.String()))
				if err != nil {
					log.Print("could not find volume group allocatable size : ", err)
					lvsMap[lv.Name] = 0
				} else {
					lvsMap[lv.Name] = freeSpace
				}
				break
			} else {
				lvsMap[lv.Name] = 0
			}
		}
	}
	return lvsMap
}
