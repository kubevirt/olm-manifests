// +build !ignore_autogenerated

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

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1beta1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConverged":             schema_pkg_apis_hco_v1beta1_HyperConverged(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedFeatureGates": schema_pkg_apis_hco_v1beta1_HyperConvergedFeatureGates(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedSpec":         schema_pkg_apis_hco_v1beta1_HyperConvergedSpec(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedStatus":       schema_pkg_apis_hco_v1beta1_HyperConvergedStatus(ref),
	}
}

func schema_pkg_apis_hco_v1beta1_HyperConverged(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "HyperConverged is the Schema for the hyperconvergeds API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedSpec", "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_hco_v1beta1_HyperConvergedFeatureGates(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "HyperConvergedFeatureGates is a set of optional feature gates to enable or disable new features that are not enabled by default yet.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"HotplugVolumes": {
						SchemaProps: spec.SchemaProps{
							Description: "Allow attaching a data volume to a running VMI",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_hco_v1beta1_HyperConvergedSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "HyperConvergedSpec defines the desired state of HyperConverged",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"bareMetalPlatform": {
						SchemaProps: spec.SchemaProps{
							Description: "BareMetalPlatform indicates whether the infrastructure is baremetal.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"localStorageClassName": {
						SchemaProps: spec.SchemaProps{
							Description: "LocalStorageClassName the name of the local storage class.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"infra": {
						SchemaProps: spec.SchemaProps{
							Description: "infra HyperConvergedConfig influences the pod configuration (currently only placement) for all the infra components needed on the virtualization enabled cluster but not necessarely directly on each node running VMs/VMIs.",
							Ref:         ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedConfig"),
						},
					},
					"workloads": {
						SchemaProps: spec.SchemaProps{
							Description: "workloads HyperConvergedConfig influences the pod configuration (currently only placement) of components which need to be running on a node where virtualization workloads should be able to run. Changes to Workloads HyperConvergedConfig can be applied only without existing workload.",
							Ref:         ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedConfig"),
						},
					},
					"featureGates": {
						SchemaProps: spec.SchemaProps{
							Description: "featureGates is a map of feature gate flags. Setting a flag to `true` will enable the feature. Setting `false` or removing the feature gate, disables the feature.",
							Ref:         ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedFeatureGates"),
						},
					},
					"version": {
						SchemaProps: spec.SchemaProps{
							Description: "operator version",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedConfig", "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedFeatureGates"},
	}
}

func schema_pkg_apis_hco_v1beta1_HyperConvergedStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "HyperConvergedStatus defines the observed state of HyperConverged",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-patch-merge-key": "type",
								"x-kubernetes-patch-strategy":  "merge",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Conditions describes the state of the HyperConverged resource.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/openshift/custom-resource-status/conditions/v1.Condition"),
									},
								},
							},
						},
					},
					"relatedObjects": {
						SchemaProps: spec.SchemaProps{
							Description: "RelatedObjects is a list of objects created and maintained by this operator. Object references will be added to this list after they have been created AND found in the cluster.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("k8s.io/api/core/v1.ObjectReference"),
									},
								},
							},
						},
					},
					"versions": {
						SchemaProps: spec.SchemaProps{
							Description: "Versions is a list of HCO component versions, as name/version pairs. The version with a name of \"operator\" is the HCO version itself, as described here: https://github.com/openshift/cluster-version-operator/blob/master/docs/dev/clusteroperator.md#version",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.Version"),
									},
								},
							},
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.Version", "github.com/openshift/custom-resource-status/conditions/v1.Condition", "k8s.io/api/core/v1.ObjectReference"},
	}
}
