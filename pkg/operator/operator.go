package operator

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/util/intstr"
)

type operator struct {
	client    *unversioned.Client
	namespace string
}

func New(client *unversioned.Client, namespace string) *operator {

	return &operator{
		client:    client,
		namespace: namespace,
	}

}

func (o *operator) ProvisionInstance() error {

	instance := "default"

	_, err := o.client.Services(o.namespace).Create(newFEService(instance))

	if err != nil {
		return err
	}
	return nil
}

func newFEService(instance string) *api.Service {
	labels := map[string]string{
		"app":      "fe",
		"instance": instance,
	}
	svc := &api.Service{
		ObjectMeta: api.ObjectMeta{
			Name:   fmt.Sprintf("front-end-%s", instance),
			Labels: labels,
		},
		Spec: api.ServiceSpec{
			Ports: []api.ServicePort{
				{
					Name:       "web",
					Port:       8080,
					TargetPort: intstr.FromInt(3000),
					Protocol:   api.ProtocolTCP,
				},
			},
			Selector: labels,
		},
	}
	return svc
}
