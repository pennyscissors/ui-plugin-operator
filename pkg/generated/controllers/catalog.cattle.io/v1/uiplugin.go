/*
Copyright 2024 Rancher Labs, Inc.

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
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	v1 "github.com/rancher/ui-plugin-operator/pkg/apis/catalog.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type UIPluginHandler func(string, *v1.UIPlugin) (*v1.UIPlugin, error)

type UIPluginController interface {
	generic.ControllerMeta
	UIPluginClient

	OnChange(ctx context.Context, name string, sync UIPluginHandler)
	OnRemove(ctx context.Context, name string, sync UIPluginHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() UIPluginCache
}

type UIPluginClient interface {
	Create(*v1.UIPlugin) (*v1.UIPlugin, error)
	Update(*v1.UIPlugin) (*v1.UIPlugin, error)
	UpdateStatus(*v1.UIPlugin) (*v1.UIPlugin, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.UIPlugin, error)
	List(namespace string, opts metav1.ListOptions) (*v1.UIPluginList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.UIPlugin, err error)
}

type UIPluginCache interface {
	Get(namespace, name string) (*v1.UIPlugin, error)
	List(namespace string, selector labels.Selector) ([]*v1.UIPlugin, error)

	AddIndexer(indexName string, indexer UIPluginIndexer)
	GetByIndex(indexName, key string) ([]*v1.UIPlugin, error)
}

type UIPluginIndexer func(obj *v1.UIPlugin) ([]string, error)

type uIPluginController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewUIPluginController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) UIPluginController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &uIPluginController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromUIPluginHandlerToHandler(sync UIPluginHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.UIPlugin
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.UIPlugin))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *uIPluginController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.UIPlugin))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateUIPluginDeepCopyOnChange(client UIPluginClient, obj *v1.UIPlugin, handler func(obj *v1.UIPlugin) (*v1.UIPlugin, error)) (*v1.UIPlugin, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *uIPluginController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *uIPluginController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *uIPluginController) OnChange(ctx context.Context, name string, sync UIPluginHandler) {
	c.AddGenericHandler(ctx, name, FromUIPluginHandlerToHandler(sync))
}

func (c *uIPluginController) OnRemove(ctx context.Context, name string, sync UIPluginHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromUIPluginHandlerToHandler(sync)))
}

func (c *uIPluginController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *uIPluginController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *uIPluginController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *uIPluginController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *uIPluginController) Cache() UIPluginCache {
	return &uIPluginCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *uIPluginController) Create(obj *v1.UIPlugin) (*v1.UIPlugin, error) {
	result := &v1.UIPlugin{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *uIPluginController) Update(obj *v1.UIPlugin) (*v1.UIPlugin, error) {
	result := &v1.UIPlugin{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *uIPluginController) UpdateStatus(obj *v1.UIPlugin) (*v1.UIPlugin, error) {
	result := &v1.UIPlugin{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *uIPluginController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *uIPluginController) Get(namespace, name string, options metav1.GetOptions) (*v1.UIPlugin, error) {
	result := &v1.UIPlugin{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *uIPluginController) List(namespace string, opts metav1.ListOptions) (*v1.UIPluginList, error) {
	result := &v1.UIPluginList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *uIPluginController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *uIPluginController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.UIPlugin, error) {
	result := &v1.UIPlugin{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type uIPluginCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *uIPluginCache) Get(namespace, name string) (*v1.UIPlugin, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.UIPlugin), nil
}

func (c *uIPluginCache) List(namespace string, selector labels.Selector) (ret []*v1.UIPlugin, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.UIPlugin))
	})

	return ret, err
}

func (c *uIPluginCache) AddIndexer(indexName string, indexer UIPluginIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.UIPlugin))
		},
	}))
}

func (c *uIPluginCache) GetByIndex(indexName, key string) (result []*v1.UIPlugin, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.UIPlugin, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.UIPlugin))
	}
	return result, nil
}

type UIPluginStatusHandler func(obj *v1.UIPlugin, status v1.UIPluginStatus) (v1.UIPluginStatus, error)

type UIPluginGeneratingHandler func(obj *v1.UIPlugin, status v1.UIPluginStatus) ([]runtime.Object, v1.UIPluginStatus, error)

func RegisterUIPluginStatusHandler(ctx context.Context, controller UIPluginController, condition condition.Cond, name string, handler UIPluginStatusHandler) {
	statusHandler := &uIPluginStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromUIPluginHandlerToHandler(statusHandler.sync))
}

func RegisterUIPluginGeneratingHandler(ctx context.Context, controller UIPluginController, apply apply.Apply,
	condition condition.Cond, name string, handler UIPluginGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &uIPluginGeneratingHandler{
		UIPluginGeneratingHandler: handler,
		apply:                     apply,
		name:                      name,
		gvk:                       controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterUIPluginStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type uIPluginStatusHandler struct {
	client    UIPluginClient
	condition condition.Cond
	handler   UIPluginStatusHandler
}

func (a *uIPluginStatusHandler) sync(key string, obj *v1.UIPlugin) (*v1.UIPlugin, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type uIPluginGeneratingHandler struct {
	UIPluginGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *uIPluginGeneratingHandler) Remove(key string, obj *v1.UIPlugin) (*v1.UIPlugin, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.UIPlugin{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *uIPluginGeneratingHandler) Handle(obj *v1.UIPlugin, status v1.UIPluginStatus) (v1.UIPluginStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.UIPluginGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
