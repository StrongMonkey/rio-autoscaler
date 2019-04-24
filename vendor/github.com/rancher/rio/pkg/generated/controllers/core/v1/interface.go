/*
Copyright 2019 Rancher Labs.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/api/core/v1"
	informers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Interface interface {
	ConfigMap() ConfigMapController
	Endpoints() EndpointsController
	Namespace() NamespaceController
	Node() NodeController
	PersistentVolumeClaim() PersistentVolumeClaimController
	Pod() PodController
	Secret() SecretController
	Service() ServiceController
	ServiceAccount() ServiceAccountController
}

func New(controllerManager *generic.ControllerManager, client clientset.CoreV1Interface,
	informers informers.Interface) Interface {
	return &version{
		controllerManager: controllerManager,
		client:            client,
		informers:         informers,
	}
}

type version struct {
	controllerManager *generic.ControllerManager
	informers         informers.Interface
	client            clientset.CoreV1Interface
}

func (c *version) ConfigMap() ConfigMapController {
	return NewConfigMapController(v1.SchemeGroupVersion.WithKind("ConfigMap"), c.controllerManager, c.client, c.informers.ConfigMaps())
}
func (c *version) Endpoints() EndpointsController {
	return NewEndpointsController(v1.SchemeGroupVersion.WithKind("Endpoints"), c.controllerManager, c.client, c.informers.Endpoints())
}
func (c *version) Namespace() NamespaceController {
	return NewNamespaceController(v1.SchemeGroupVersion.WithKind("Namespace"), c.controllerManager, c.client, c.informers.Namespaces())
}
func (c *version) Node() NodeController {
	return NewNodeController(v1.SchemeGroupVersion.WithKind("Node"), c.controllerManager, c.client, c.informers.Nodes())
}
func (c *version) PersistentVolumeClaim() PersistentVolumeClaimController {
	return NewPersistentVolumeClaimController(v1.SchemeGroupVersion.WithKind("PersistentVolumeClaim"), c.controllerManager, c.client, c.informers.PersistentVolumeClaims())
}
func (c *version) Pod() PodController {
	return NewPodController(v1.SchemeGroupVersion.WithKind("Pod"), c.controllerManager, c.client, c.informers.Pods())
}
func (c *version) Secret() SecretController {
	return NewSecretController(v1.SchemeGroupVersion.WithKind("Secret"), c.controllerManager, c.client, c.informers.Secrets())
}
func (c *version) Service() ServiceController {
	return NewServiceController(v1.SchemeGroupVersion.WithKind("Service"), c.controllerManager, c.client, c.informers.Services())
}
func (c *version) ServiceAccount() ServiceAccountController {
	return NewServiceAccountController(v1.SchemeGroupVersion.WithKind("ServiceAccount"), c.controllerManager, c.client, c.informers.ServiceAccounts())
}