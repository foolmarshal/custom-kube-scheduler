package extender

import (
	"context"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func getPVCNamesFromPodSpec(pod v1.Pod) []string {
	var claimNames []string
	for _, volume := range pod.Spec.Volumes {
		log.Print("VolumeSource is equal to PVCVS : ", volume.VolumeSource.PersistentVolumeClaim != nil)
		if volume.VolumeSource.PersistentVolumeClaim != nil {
			log.Print("claimName : ", volume.PersistentVolumeClaim.ClaimName)
			claimNames = append(claimNames, volume.PersistentVolumeClaim.ClaimName)
		}
	}
	log.Print("claimNames : ", claimNames)
	return claimNames
}

func getTotalRequestedVolumeByPVCs(pod v1.Pod) int64 {
	var totalRequestVolume int64 = 0
	claimNames := getPVCNamesFromPodSpec(pod)
	for _, claimName := range claimNames {
		pvc, err := GetPVC(claimName, pod.Namespace)
		if err != nil {
			log.Print("could not find pvc for name : ", claimName, " and namespace : ", pod.Namespace, " because : ", err)
		} else {
			storage, _ := pvc.Spec.Resources.Requests.Storage().AsInt64()
			log.Print("requests : ", pvc.Spec.Resources.Requests)
			log.Print("requests_storage : ", storage)
			totalRequestVolume = totalRequestVolume + storage
		}
	}
	return totalRequestVolume
}

// GetPVC returns a PersistentVolumeClaim object using the pvc name passed.
func GetPVC(name string, namespace string) (*corev1.PersistentVolumeClaim, error) {
	k8sClient, err := getK8sClient()
	if err != nil {
		log.Print("could not initialize k8s client : ", err)
		return nil, errors.Wrap(err, "error while getting k8s client")
	}
	pvc, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while getting persistent volume claim")
	}
	return pvc, nil
}
