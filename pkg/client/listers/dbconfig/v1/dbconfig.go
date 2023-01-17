/*
Copyright The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DbConfigLister helps list DbConfigs.
// All objects returned here must be treated as read-only.
type DbConfigLister interface {
	// List lists all DbConfigs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.DbConfig, err error)
	// DbConfigs returns an object that can list and get DbConfigs.
	DbConfigs(namespace string) DbConfigNamespaceLister
	DbConfigListerExpansion
}

// dbConfigLister implements the DbConfigLister interface.
type dbConfigLister struct {
	indexer cache.Indexer
}

// NewDbConfigLister returns a new DbConfigLister.
func NewDbConfigLister(indexer cache.Indexer) DbConfigLister {
	return &dbConfigLister{indexer: indexer}
}

// List lists all DbConfigs in the indexer.
func (s *dbConfigLister) List(selector labels.Selector) (ret []*v1.DbConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.DbConfig))
	})
	return ret, err
}

// DbConfigs returns an object that can list and get DbConfigs.
func (s *dbConfigLister) DbConfigs(namespace string) DbConfigNamespaceLister {
	return dbConfigNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DbConfigNamespaceLister helps list and get DbConfigs.
// All objects returned here must be treated as read-only.
type DbConfigNamespaceLister interface {
	// List lists all DbConfigs in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.DbConfig, err error)
	// Get retrieves the DbConfig from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.DbConfig, error)
	DbConfigNamespaceListerExpansion
}

// dbConfigNamespaceLister implements the DbConfigNamespaceLister
// interface.
type dbConfigNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all DbConfigs in the indexer for a given namespace.
func (s dbConfigNamespaceLister) List(selector labels.Selector) (ret []*v1.DbConfig, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.DbConfig))
	})
	return ret, err
}

// Get retrieves the DbConfig from the indexer for a given namespace and name.
func (s dbConfigNamespaceLister) Get(name string) (*v1.DbConfig, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("dbconfig"), name)
	}
	return obj.(*v1.DbConfig), nil
}
