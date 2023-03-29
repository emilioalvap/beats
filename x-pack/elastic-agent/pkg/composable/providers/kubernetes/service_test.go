// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package kubernetes

import (
	"testing"

	"github.com/elastic/beats/v7/libbeat/common/kubernetes/metadata"

	"github.com/elastic/beats/v7/libbeat/common"

	"github.com/stretchr/testify/assert"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/elastic/beats/v7/libbeat/common/kubernetes"
)

func TestGenerateServiceData(t *testing.T) {
	uid := "005f3b90-4b9d-12f8-acf0-31020a840133"
	service := &kubernetes.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testsvc",
			UID:       types.UID(uid),
			Namespace: "testns",
			Labels: map[string]string{
				"foo":        "bar",
				"with-dash":  "dash-value",
				"with/slash": "some/path",
			},
			Annotations: map[string]string{
				"baz": "ban",
			},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		Spec: v1.ServiceSpec{
			ClusterIP: "1.2.3.4",
			Selector: map[string]string{
				"app":   "istiod",
				"istio": "pilot",
			},
		},
	}

	data := generateServiceData(
		service,
		&Config{},
		&svcMeta{},
		common.MapStr{
			"nsa": "nsb",
		})

	mapping := map[string]interface{}{
		"service": common.MapStr{
			"uid":  string(service.GetUID()),
			"name": service.GetName(),
			"ip":   service.Spec.ClusterIP,
		},
		"namespace_annotations": common.MapStr{
			"nsa": "nsb",
		},
		"annotations": common.MapStr{
			"baz": "ban",
		},
		"labels": common.MapStr{
			"foo":        "bar",
			"with-dash":  "dash-value",
			"with/slash": "some/path",
		},
	}

	processors := map[string]interface{}{
		"orchestrator": common.MapStr{
			"cluster": common.MapStr{
				"name": "devcluster",
				"url":  "8.8.8.8:9090"},
		}, "kubernetes": common.MapStr{
			"service": common.MapStr{
				"uid":  string(service.GetUID()),
				"name": service.GetName(),
				"ip":   "1.2.3.4",
			},
			"labels": common.MapStr{
				"foo":        "bar",
				"with-dash":  "dash-value",
				"with/slash": "some/path",
			},
			"annotations": common.MapStr{
				"baz": "ban",
			},
		},
	}

	assert.Equal(t, service, data.service)
	assert.Equal(t, mapping, data.mapping)
	for _, v := range data.processors {
		k := v["add_fields"].(map[string]interface{})
		target := k["target"].(string)
		fields := k["fields"]
		assert.Equal(t, processors[target], fields)
	}
}

type svcMeta struct{}

// Generate generates svc metadata from a resource object
// Metadata map is in the following form:
//
//	{
//		  "kubernetes": {},
//	   "some.ecs.field": "asdf"
//	}
//
// All Kubernetes fields that need to be stored under kubernetes. prefix are populated by
// GenerateK8s method while fields that are part of ECS are generated by GenerateECS method
func (s *svcMeta) Generate(obj kubernetes.Resource, opts ...metadata.FieldOptions) common.MapStr {
	ecsFields := s.GenerateECS(obj)
	meta := common.MapStr{
		"kubernetes": s.GenerateK8s(obj, opts...),
	}
	meta.DeepUpdate(ecsFields)
	return meta
}

// GenerateECS generates svc ECS metadata from a resource object
func (s *svcMeta) GenerateECS(obj kubernetes.Resource) common.MapStr {
	return common.MapStr{
		"orchestrator": common.MapStr{
			"cluster": common.MapStr{
				"name": "devcluster",
				"url":  "8.8.8.8:9090",
			},
		},
	}
}

// GenerateK8s generates svc metadata from a resource object
func (s *svcMeta) GenerateK8s(obj kubernetes.Resource, opts ...metadata.FieldOptions) common.MapStr {
	k8sNode := obj.(*kubernetes.Service)
	return common.MapStr{
		"service": common.MapStr{
			"uid":  string(k8sNode.GetUID()),
			"name": k8sNode.GetName(),
			"ip":   "1.2.3.4",
		},
		"labels": common.MapStr{
			"foo":        "bar",
			"with-dash":  "dash-value",
			"with/slash": "some/path",
		},
		"annotations": common.MapStr{
			"baz": "ban",
		},
	}
}

// GenerateFromName generates svc metadata from a node name
func (s *svcMeta) GenerateFromName(name string, opts ...metadata.FieldOptions) common.MapStr {
	return nil
}
