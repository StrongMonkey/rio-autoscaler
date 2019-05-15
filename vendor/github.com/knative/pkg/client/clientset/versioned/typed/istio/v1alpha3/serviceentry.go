/*
Copyright 2018 The Knative Authors

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha3

import (
	v1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
	scheme "github.com/knative/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ServiceEntriesGetter has a method to return a ServiceEntryInterface.
// A group's client should implement this interface.
type ServiceEntriesGetter interface {
	ServiceEntries(namespace string) ServiceEntryInterface
}

// ServiceEntryInterface has methods to work with ServiceEntry resources.
type ServiceEntryInterface interface {
	Create(*v1alpha3.ServiceEntry) (*v1alpha3.ServiceEntry, error)
	Update(*v1alpha3.ServiceEntry) (*v1alpha3.ServiceEntry, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha3.ServiceEntry, error)
	List(opts v1.ListOptions) (*v1alpha3.ServiceEntryList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha3.ServiceEntry, err error)
	ServiceEntryExpansion
}

// serviceEntries implements ServiceEntryInterface
type serviceEntries struct {
	client rest.Interface
	ns     string
}

// newServiceEntries returns a ServiceEntries
func newServiceEntries(c *NetworkingV1alpha3Client, namespace string) *serviceEntries {
	return &serviceEntries{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the serviceEntry, and returns the corresponding serviceEntry object, and an error if there is any.
func (c *serviceEntries) Get(name string, options v1.GetOptions) (result *v1alpha3.ServiceEntry, err error) {
	result = &v1alpha3.ServiceEntry{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("serviceentries").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ServiceEntries that match those selectors.
func (c *serviceEntries) List(opts v1.ListOptions) (result *v1alpha3.ServiceEntryList, err error) {
	result = &v1alpha3.ServiceEntryList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("serviceentries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested serviceEntries.
func (c *serviceEntries) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("serviceentries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a serviceEntry and creates it.  Returns the server's representation of the serviceEntry, and an error, if there is any.
func (c *serviceEntries) Create(serviceEntry *v1alpha3.ServiceEntry) (result *v1alpha3.ServiceEntry, err error) {
	result = &v1alpha3.ServiceEntry{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("serviceentries").
		Body(serviceEntry).
		Do().
		Into(result)
	return
}

// Update takes the representation of a serviceEntry and updates it. Returns the server's representation of the serviceEntry, and an error, if there is any.
func (c *serviceEntries) Update(serviceEntry *v1alpha3.ServiceEntry) (result *v1alpha3.ServiceEntry, err error) {
	result = &v1alpha3.ServiceEntry{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("serviceentries").
		Name(serviceEntry.Name).
		Body(serviceEntry).
		Do().
		Into(result)
	return
}

// Delete takes name of the serviceEntry and deletes it. Returns an error if one occurs.
func (c *serviceEntries) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("serviceentries").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *serviceEntries) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("serviceentries").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched serviceEntry.
func (c *serviceEntries) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha3.ServiceEntry, err error) {
	result = &v1alpha3.ServiceEntry{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("serviceentries").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
