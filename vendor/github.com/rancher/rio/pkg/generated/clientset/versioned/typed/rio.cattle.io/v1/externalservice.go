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
	"time"

	v1 "github.com/rancher/rio/pkg/apis/rio.cattle.io/v1"
	scheme "github.com/rancher/rio/pkg/generated/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ExternalServicesGetter has a method to return a ExternalServiceInterface.
// A group's client should implement this interface.
type ExternalServicesGetter interface {
	ExternalServices(namespace string) ExternalServiceInterface
}

// ExternalServiceInterface has methods to work with ExternalService resources.
type ExternalServiceInterface interface {
	Create(*v1.ExternalService) (*v1.ExternalService, error)
	Update(*v1.ExternalService) (*v1.ExternalService, error)
	UpdateStatus(*v1.ExternalService) (*v1.ExternalService, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.ExternalService, error)
	List(opts metav1.ListOptions) (*v1.ExternalServiceList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ExternalService, err error)
	ExternalServiceExpansion
}

// externalServices implements ExternalServiceInterface
type externalServices struct {
	client rest.Interface
	ns     string
}

// newExternalServices returns a ExternalServices
func newExternalServices(c *RioV1Client, namespace string) *externalServices {
	return &externalServices{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the externalService, and returns the corresponding externalService object, and an error if there is any.
func (c *externalServices) Get(name string, options metav1.GetOptions) (result *v1.ExternalService, err error) {
	result = &v1.ExternalService{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("externalservices").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ExternalServices that match those selectors.
func (c *externalServices) List(opts metav1.ListOptions) (result *v1.ExternalServiceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ExternalServiceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("externalservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested externalServices.
func (c *externalServices) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("externalservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a externalService and creates it.  Returns the server's representation of the externalService, and an error, if there is any.
func (c *externalServices) Create(externalService *v1.ExternalService) (result *v1.ExternalService, err error) {
	result = &v1.ExternalService{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("externalservices").
		Body(externalService).
		Do().
		Into(result)
	return
}

// Update takes the representation of a externalService and updates it. Returns the server's representation of the externalService, and an error, if there is any.
func (c *externalServices) Update(externalService *v1.ExternalService) (result *v1.ExternalService, err error) {
	result = &v1.ExternalService{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("externalservices").
		Name(externalService.Name).
		Body(externalService).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *externalServices) UpdateStatus(externalService *v1.ExternalService) (result *v1.ExternalService, err error) {
	result = &v1.ExternalService{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("externalservices").
		Name(externalService.Name).
		SubResource("status").
		Body(externalService).
		Do().
		Into(result)
	return
}

// Delete takes name of the externalService and deletes it. Returns an error if one occurs.
func (c *externalServices) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("externalservices").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *externalServices) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("externalservices").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched externalService.
func (c *externalServices) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ExternalService, err error) {
	result = &v1.ExternalService{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("externalservices").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
