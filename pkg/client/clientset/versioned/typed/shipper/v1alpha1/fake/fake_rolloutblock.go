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

package fake

import (
	v1alpha1 "github.com/bookingcom/shipper/pkg/apis/shipper/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRolloutBlocks implements RolloutBlockInterface
type FakeRolloutBlocks struct {
	Fake *FakeShipperV1alpha1
	ns   string
}

var rolloutblocksResource = schema.GroupVersionResource{Group: "shipper.booking.com", Version: "v1alpha1", Resource: "rolloutblocks"}

var rolloutblocksKind = schema.GroupVersionKind{Group: "shipper.booking.com", Version: "v1alpha1", Kind: "RolloutBlock"}

// Get takes name of the rolloutBlock, and returns the corresponding rolloutBlock object, and an error if there is any.
func (c *FakeRolloutBlocks) Get(name string, options v1.GetOptions) (result *v1alpha1.RolloutBlock, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(rolloutblocksResource, c.ns, name), &v1alpha1.RolloutBlock{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RolloutBlock), err
}

// List takes label and field selectors, and returns the list of RolloutBlocks that match those selectors.
func (c *FakeRolloutBlocks) List(opts v1.ListOptions) (result *v1alpha1.RolloutBlockList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(rolloutblocksResource, rolloutblocksKind, c.ns, opts), &v1alpha1.RolloutBlockList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.RolloutBlockList{}
	for _, item := range obj.(*v1alpha1.RolloutBlockList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested rolloutBlocks.
func (c *FakeRolloutBlocks) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(rolloutblocksResource, c.ns, opts))

}

// Create takes the representation of a rolloutBlock and creates it.  Returns the server's representation of the rolloutBlock, and an error, if there is any.
func (c *FakeRolloutBlocks) Create(rolloutBlock *v1alpha1.RolloutBlock) (result *v1alpha1.RolloutBlock, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(rolloutblocksResource, c.ns, rolloutBlock), &v1alpha1.RolloutBlock{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RolloutBlock), err
}

// Update takes the representation of a rolloutBlock and updates it. Returns the server's representation of the rolloutBlock, and an error, if there is any.
func (c *FakeRolloutBlocks) Update(rolloutBlock *v1alpha1.RolloutBlock) (result *v1alpha1.RolloutBlock, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(rolloutblocksResource, c.ns, rolloutBlock), &v1alpha1.RolloutBlock{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RolloutBlock), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRolloutBlocks) UpdateStatus(rolloutBlock *v1alpha1.RolloutBlock) (*v1alpha1.RolloutBlock, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(rolloutblocksResource, "status", c.ns, rolloutBlock), &v1alpha1.RolloutBlock{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RolloutBlock), err
}

// Delete takes name of the rolloutBlock and deletes it. Returns an error if one occurs.
func (c *FakeRolloutBlocks) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(rolloutblocksResource, c.ns, name), &v1alpha1.RolloutBlock{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRolloutBlocks) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(rolloutblocksResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.RolloutBlockList{})
	return err
}

// Patch applies the patch and returns the patched rolloutBlock.
func (c *FakeRolloutBlocks) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RolloutBlock, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(rolloutblocksResource, c.ns, name, data, subresources...), &v1alpha1.RolloutBlock{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RolloutBlock), err
}
