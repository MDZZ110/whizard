package receive_ingestor

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

func (r *ReceiveIngestor) service() (runtime.Object, resources.Operation, error) {
	var s = &corev1.Service{ObjectMeta: r.meta(r.name("operated"))}

	s.Spec = corev1.ServiceSpec{
		ClusterIP: corev1.ClusterIPNone,
		Selector:  r.labels(),
		Ports: []corev1.ServicePort{
			{
				Protocol: corev1.ProtocolTCP,
				Name:     "grpc",
				Port:     10901,
			},
			{
				Protocol: corev1.ProtocolTCP,
				Name:     "http",
				Port:     10902,
			},
			{
				Name:     "remote-write",
				Protocol: corev1.ProtocolTCP,
				Port:     19291,
			},
		},
	}
	return s, resources.OperationCreateOrUpdate, nil
}
