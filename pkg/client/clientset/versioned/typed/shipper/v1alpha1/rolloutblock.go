/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha1

import (
	v1alpha1 "github.com/bookingcom/shipper/pkg/apis/shipper/v1alpha1"
	scheme "github.com/bookingcom/shipper/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// RolloutBlocksGetter has a method to return a RolloutBlockInterface.
// A group's client should implement this interface.
type RolloutBlocksGetter interface {
	RolloutBlocks(namespace string) RolloutBlockInterface
}

// RolloutBlockInterface has methods to work with RolloutBlock resources.
type RolloutBlockInterface interface {
	Create(*v1alpha1.RolloutBlock) (*v1alpha1.RolloutBlock, error)
	Update(*v1alpha1.RolloutBlock) (*v1alpha1.RolloutBlock, error)
	UpdateStatus(*v1alpha1.RolloutBlock) (*v1alpha1.RolloutBlock, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.RolloutBlock, error)
	List(opts v1.ListOptions) (*v1alpha1.RolloutBlockList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RolloutBlock, err error)
	RolloutBlockExpansion
}

// rolloutBlocks implements RolloutBlockInterface
type rolloutBlocks struct {
	client rest.Interface
	ns     string
}

// newRolloutBlocks returns a RolloutBlocks
func newRolloutBlocks(c *ShipperV1alpha1Client, namespace string) *rolloutBlocks {
	return &rolloutBlocks{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the rolloutBlock, and returns the corresponding rolloutBlock object, and an error if there is any.
func (c *rolloutBlocks) Get(name string, options v1.GetOptions) (result *v1alpha1.RolloutBlock, err error) {
	result = &v1alpha1.RolloutBlock{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rolloutblocks").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of RolloutBlocks that match those selectors.
func (c *rolloutBlocks) List(opts v1.ListOptions) (result *v1alpha1.RolloutBlockList, err error) {
	result = &v1alpha1.RolloutBlockList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rolloutblocks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested rolloutBlocks.
func (c *rolloutBlocks) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("rolloutblocks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a rolloutBlock and creates it.  Returns the server's representation of the rolloutBlock, and an error, if there is any.
func (c *rolloutBlocks) Create(rolloutBlock *v1alpha1.RolloutBlock) (result *v1alpha1.RolloutBlock, err error) {
	result = &v1alpha1.RolloutBlock{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("rolloutblocks").
		Body(rolloutBlock).
		Do().
		Into(result)
	return
}

// Update takes the representation of a rolloutBlock and updates it. Returns the server's representation of the rolloutBlock, and an error, if there is any.
func (c *rolloutBlocks) Update(rolloutBlock *v1alpha1.RolloutBlock) (result *v1alpha1.RolloutBlock, err error) {
	result = &v1alpha1.RolloutBlock{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("rolloutblocks").
		Name(rolloutBlock.Name).
		Body(rolloutBlock).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *rolloutBlocks) UpdateStatus(rolloutBlock *v1alpha1.RolloutBlock) (result *v1alpha1.RolloutBlock, err error) {
	result = &v1alpha1.RolloutBlock{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("rolloutblocks").
		Name(rolloutBlock.Name).
		SubResource("status").
		Body(rolloutBlock).
		Do().
		Into(result)
	return
}

// Delete takes name of the rolloutBlock and deletes it. Returns an error if one occurs.
func (c *rolloutBlocks) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("rolloutblocks").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *rolloutBlocks) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("rolloutblocks").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched rolloutBlock.
func (c *rolloutBlocks) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RolloutBlock, err error) {
	result = &v1alpha1.RolloutBlock{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("rolloutblocks").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
