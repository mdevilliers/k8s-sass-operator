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

	services := serviceDefinitions(instance)

	for _, service := range services {

		_, err := o.client.Services(o.namespace).Create(service)

		if err != nil {
			return err
		}
	}

	return nil
}

func serviceDefinitions(instance string) []*api.Service {
	return []*api.Service{
		frontEndService(instance),
		storeService(instance),
		userService(instance),
	}
}

func frontEndService(instance string) *api.Service {
	labels := map[string]string{
		"app":      "front-end",
		"instance": instance,
	}
	return &api.Service{
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
}

func storeService(instance string) *api.Service {
	labels := map[string]string{
		"app":      "store-service",
		"instance": instance,
	}
	return &api.Service{
		ObjectMeta: api.ObjectMeta{
			Name:   fmt.Sprintf("store-service-%s", instance),
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
}

func userService(instance string) *api.Service {
	labels := map[string]string{
		"app":      "user-service",
		"instance": instance,
	}
	return &api.Service{
		ObjectMeta: api.ObjectMeta{
			Name:   fmt.Sprintf("user-service-%s", instance),
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
}
