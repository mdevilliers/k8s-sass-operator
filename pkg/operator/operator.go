package operator

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"

	"k8s.io/kubernetes/pkg/api"
	apierrors "k8s.io/kubernetes/pkg/api/errors"
	unversionedAPI "k8s.io/kubernetes/pkg/api/unversioned"
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

	//TODO : is there an instance of this name already
	// if so patch or delete and create
	// at the moment just press on regardless

	instance := "default"
	services := serviceDefinitions(instance)

	logrus.Info("deploying services")
	for _, service := range services {

		logrus.Info("deploying : ", service.Name)

		_, err := o.client.Services(o.namespace).Create(&service)

		err = filterKubernetesResourceAlreadyExistError(err)

		if err != nil {
			return nil
		}
	}

	replicationControllers := replicationControllercDefinitions(instance)

	logrus.Info("deploying replication controllers")

	for _, rc := range replicationControllers {

		logrus.Info("deploying : ", rc.Name)

		_, err := o.client.ReplicationControllers(o.namespace).Create(&rc)

		err = filterKubernetesResourceAlreadyExistError(err)

		if err != nil {

			logrus.Info(err)

			return nil
		}
	}

	return nil
}

func serviceDefinitions(instance string) []api.Service {
	return []api.Service{
		frontEndService(instance),
		storeService(instance),
		userService(instance),
	}
}

func frontEndService(instance string) api.Service {
	labels := map[string]string{
		"app":      "front-end",
		"instance": instance,
	}
	return api.Service{
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
			Type:     api.ServiceTypeNodePort,
		},
	}
}

func storeService(instance string) api.Service {
	labels := map[string]string{
		"app":      "store-service",
		"instance": instance,
	}
	return api.Service{
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

func userService(instance string) api.Service {
	labels := map[string]string{
		"app":      "user-service",
		"instance": instance,
	}
	return api.Service{
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

func replicationControllercDefinitions(instance string) []api.ReplicationController {
	return []api.ReplicationController{
		frontEndServiceRC(instance),
		storeServiceRC(instance),
		userServiceRC(instance),
	}
}

func frontEndServiceRC(instance string) api.ReplicationController {
	labels := map[string]string{
		"app":      "front-end-service",
		"instance": instance,
		"version":  "1",
	}
	return api.ReplicationController{
		ObjectMeta: api.ObjectMeta{
			Name: "fe-service-rc",
		},
		Spec: api.ReplicationControllerSpec{
			Replicas: 3,
			Selector: labels,
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: labels,
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						api.Container{
							Name:            "web",
							Image:           "sass-infrastructure/fe",
							ImagePullPolicy: api.PullIfNotPresent,
							Ports: []api.ContainerPort{
								api.ContainerPort{
									ContainerPort: 3000,
									Protocol:      api.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}
}

func storeServiceRC(instance string) api.ReplicationController {
	labels := map[string]string{
		"app":      "store-service",
		"instance": instance,
		"version":  "1",
	}
	return api.ReplicationController{
		ObjectMeta: api.ObjectMeta{
			Name: "store-service-rc",
		},
		Spec: api.ReplicationControllerSpec{
			Replicas: 3,
			Selector: labels,
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: labels,
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						api.Container{
							Name:            "api",
							Image:           "sass-infrastructure/store-service",
							ImagePullPolicy: api.PullIfNotPresent,
							Ports: []api.ContainerPort{
								api.ContainerPort{
									ContainerPort: 3000,
									Protocol:      api.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}
}

func userServiceRC(instance string) api.ReplicationController {
	labels := map[string]string{
		"app":      "user-service",
		"instance": instance,
		"version":  "1",
	}
	return api.ReplicationController{
		ObjectMeta: api.ObjectMeta{
			Name: "user-service-rc",
		},
		Spec: api.ReplicationControllerSpec{
			Replicas: 3,
			Selector: labels,
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: labels,
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						api.Container{
							Name:            "api",
							Image:           "sass-infrastructure/user-service",
							ImagePullPolicy: api.PullIfNotPresent,
							Ports: []api.ContainerPort{
								api.ContainerPort{
									ContainerPort: 3000,
									Protocol:      api.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}
}

func filterKubernetesResourceAlreadyExistError(err error) error {
	se, ok := err.(*apierrors.StatusError)
	if !ok {
		return err
	}
	if se.Status().Code == http.StatusConflict && se.Status().Reason == unversionedAPI.StatusReasonAlreadyExists {
		return nil
	}
	return err
}
