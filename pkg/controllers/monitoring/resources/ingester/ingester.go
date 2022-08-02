package ingester

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

const (
	secretsDir = "/etc/thanos/secrets"
	storageDir = "/thanos"
)

type Ingester struct {
	resources.BaseReconciler
	ingester *monitoringv1alpha1.Ingester
}

func New(reconciler resources.BaseReconciler, ingester *monitoringv1alpha1.Ingester) *Ingester {
	return &Ingester{
		BaseReconciler: reconciler,
		ingester:       ingester,
	}
}

func (r *Ingester) labels() map[string]string {
	labels := r.BaseLabels()
	labels[resources.LabelNameAppName] = resources.AppNameIngester
	labels[resources.LabelNameAppManagedBy] = r.ingester.Name
	return labels
}

func (r *Ingester) name(nameSuffix ...string) string {
	return resources.QualifiedName(resources.AppNameIngester, r.ingester.Name, nameSuffix...)
}

func (r *Ingester) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:            name,
		Namespace:       r.ingester.Namespace,
		Labels:          r.labels(),
		OwnerReferences: r.OwnerReferences(),
	}
}

func (r *Ingester) OwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion: r.ingester.APIVersion,
			Kind:       r.ingester.Kind,
			Name:       r.ingester.Name,
			UID:        r.ingester.UID,
			Controller: pointer.BoolPtr(true),
		},
	}
}

func (r *Ingester) GrpcAddrs() []string {
	var addrs []string
	if r.ingester.Spec.Replicas == nil {
		addrs = make([]string, 1)
	} else {
		addrs = make([]string, *r.ingester.Spec.Replicas)
	}
	for i := range addrs {
		addrs[i] = fmt.Sprintf("%s-%d.%s.%s.svc:%d",
			r.name(), i, r.name(resources.ServiceNameSuffixOperated), r.ingester.Namespace, resources.ThanosGRPCPort)
	}
	return addrs
}

func (r *Ingester) Reconcile() error {
	return r.ReconcileResources([]resources.Resource{
		r.statefulSet,
		r.service,
	})
}
