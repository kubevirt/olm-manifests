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
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.CertRotateConfigCA":          schema_pkg_apis_hco_v1beta1_CertRotateConfigCA(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.CertRotateConfigServer":      schema_pkg_apis_hco_v1beta1_CertRotateConfigServer(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConverged":              schema_pkg_apis_hco_v1beta1_HyperConverged(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedCertConfig":    schema_pkg_apis_hco_v1beta1_HyperConvergedCertConfig(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedFeatureGates":  schema_pkg_apis_hco_v1beta1_HyperConvergedFeatureGates(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedObsoleteCPUs":  schema_pkg_apis_hco_v1beta1_HyperConvergedObsoleteCPUs(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedSpec":          schema_pkg_apis_hco_v1beta1_HyperConvergedSpec(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedStatus":        schema_pkg_apis_hco_v1beta1_HyperConvergedStatus(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.LiveMigrationConfigurations": schema_pkg_apis_hco_v1beta1_LiveMigrationConfigurations(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.MediatedHostDevice":          schema_pkg_apis_hco_v1beta1_MediatedHostDevice(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.OperandResourceRequirements": schema_pkg_apis_hco_v1beta1_OperandResourceRequirements(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.PciHostDevice":               schema_pkg_apis_hco_v1beta1_PciHostDevice(ref),
		"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.PermittedHostDevices":        schema_pkg_apis_hco_v1beta1_PermittedHostDevices(ref),
	}
}

func schema_pkg_apis_hco_v1beta1_CertRotateConfigCA(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "CertRotateConfigCA contains the tunables for TLS certificates.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"duration": {
						SchemaProps: spec.SchemaProps{
							Description: "The requested 'duration' (i.e. lifetime) of the Certificate. This should comply with golang's ParseDuration format (https://golang.org/pkg/time/#ParseDuration)",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.Duration"),
						},
					},
					"renewBefore": {
						SchemaProps: spec.SchemaProps{
							Description: "The amount of time before the currently issued certificate's `notAfter` time that we will begin to attempt to renew the certificate. This should comply with golang's ParseDuration format (https://golang.org/pkg/time/#ParseDuration)",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.Duration"),
						},
					},
				},
				Required: []string{"duration", "renewBefore"},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1.Duration"},
	}
}

func schema_pkg_apis_hco_v1beta1_CertRotateConfigServer(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "CertRotateConfigServer contains the tunables for TLS certificates.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"duration": {
						SchemaProps: spec.SchemaProps{
							Description: "The requested 'duration' (i.e. lifetime) of the Certificate. This should comply with golang's ParseDuration format (https://golang.org/pkg/time/#ParseDuration)",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.Duration"),
						},
					},
					"renewBefore": {
						SchemaProps: spec.SchemaProps{
							Description: "The amount of time before the currently issued certificate's `notAfter` time that we will begin to attempt to renew the certificate. This should comply with golang's ParseDuration format (https://golang.org/pkg/time/#ParseDuration)",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.Duration"),
						},
					},
				},
				Required: []string{"duration", "renewBefore"},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1.Duration"},
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
				Required: []string{"spec"},
			},
		},
		Dependencies: []string{
			"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedSpec", "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_hco_v1beta1_HyperConvergedCertConfig(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "HyperConvergedCertConfig holds the CertConfig entries for the HCO operands",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"ca": {
						SchemaProps: spec.SchemaProps{
							Description: "CA configuration - CA certs are kept in the CA bundle as long as they are valid",
							Ref:         ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.CertRotateConfigCA"),
						},
					},
					"server": {
						SchemaProps: spec.SchemaProps{
							Description: "Server configuration - Certs are rotated and discarded",
							Ref:         ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.CertRotateConfigServer"),
						},
					},
				},
				Required: []string{"ca", "server"},
			},
		},
		Dependencies: []string{
			"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.CertRotateConfigCA", "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.CertRotateConfigServer"},
	}
}

func schema_pkg_apis_hco_v1beta1_HyperConvergedFeatureGates(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "HyperConvergedFeatureGates is a set of optional feature gates to enable or disable new features that are not enabled by default yet.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"withHostPassthroughCPU": {
						SchemaProps: spec.SchemaProps{
							Description: "Allow migrating a virtual machine with CPU host-passthrough mode. This should be enabled only when the Cluster is homogeneous from CPU HW perspective doc here",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"sriovLiveMigration": {
						SchemaProps: spec.SchemaProps{
							Description: "Allow migrating a virtual machine with SRIOV interfaces. When enabled virt-launcher pods of virtual machines with SRIOV interfaces run with CAP_SYS_RESOURCE capability. This may degrade virt-launcher security.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
				},
				Required: []string{"withHostPassthroughCPU", "sriovLiveMigration"},
			},
		},
	}
}

func schema_pkg_apis_hco_v1beta1_HyperConvergedObsoleteCPUs(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "HyperConvergedObsoleteCPUs allows avoiding scheduling of VMs for obsolete CPU models",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"minCPUModel": {
						SchemaProps: spec.SchemaProps{
							Description: "MinCPUModel is the Minimum CPU model that is used for basic CPU features; e.g. Penryn or Haswell",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"cpuModels": {
						SchemaProps: spec.SchemaProps{
							Description: "CPUModels is a list of obsolete CPU models. When the node-labeller obtains the list of obsolete CPU models, it eliminates those CPU models and creates labels for valid CPU models.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
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
					"liveMigrationConfig": {
						SchemaProps: spec.SchemaProps{
							Description: "Live migration limits and timeouts are applied so that migration processes do not overwhelm the cluster.",
							Ref:         ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.LiveMigrationConfigurations"),
						},
					},
					"permittedHostDevices": {
						SchemaProps: spec.SchemaProps{
							Description: "PermittedHostDevices holds inforamtion about devices allowed for passthrough",
							Ref:         ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.PermittedHostDevices"),
						},
					},
					"certConfig": {
						SchemaProps: spec.SchemaProps{
							Description: "certConfig holds the rotation policy for internal, self-signed certificates",
							Ref:         ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedCertConfig"),
						},
					},
					"resourceRequirements": {
						SchemaProps: spec.SchemaProps{
							Description: "ResourceRequirements describes the resource requirements for the operand workloads.",
							Ref:         ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.OperandResourceRequirements"),
						},
					},
					"scratchSpaceStorageClass": {
						SchemaProps: spec.SchemaProps{
							Description: "Override the storage class used for scratch space during transfer operations. The scratch space storage class is determined in the following order: value of scratchSpaceStorageClass, if that doesn't exist, use the default storage class, if there is no default storage class, use the storage class of the DataVolume, if no storage class specified, use no storage class for scratch space",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"vddkInitImage": {
						SchemaProps: spec.SchemaProps{
							Description: "VDDK Init Image eventually used to import VMs from external providers",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"obsoleteCPUs": {
						SchemaProps: spec.SchemaProps{
							Description: "ObsoleteCPUs allows avoiding scheduling of VMs for obsolete CPU models",
							Ref:         ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedObsoleteCPUs"),
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
				Required: []string{"featureGates", "liveMigrationConfig", "certConfig"},
			},
		},
		Dependencies: []string{
			"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedCertConfig", "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedConfig", "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedFeatureGates", "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.HyperConvergedObsoleteCPUs", "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.LiveMigrationConfigurations", "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.OperandResourceRequirements", "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.PermittedHostDevices"},
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

func schema_pkg_apis_hco_v1beta1_LiveMigrationConfigurations(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "LiveMigrationConfigurations - Live migration limits and timeouts are applied so that migration processes do not overwhelm the cluster.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"parallelMigrationsPerCluster": {
						SchemaProps: spec.SchemaProps{
							Description: "Number of migrations running in parallel in the cluster.",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"parallelOutboundMigrationsPerNode": {
						SchemaProps: spec.SchemaProps{
							Description: "Maximum number of outbound migrations per node.",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"bandwidthPerMigration": {
						SchemaProps: spec.SchemaProps{
							Description: "Bandwidth limit of each migration, in MiB/s.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"completionTimeoutPerGiB": {
						SchemaProps: spec.SchemaProps{
							Description: "The migration will be canceled if it has not completed in this time, in seconds per GiB of memory. For example, a virtual machine instance with 6GiB memory will timeout if it has not completed migration in 4800 seconds. If the Migration Method is BlockMigration, the size of the migrating disks is included in the calculation.",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"progressTimeout": {
						SchemaProps: spec.SchemaProps{
							Description: "The migration will be canceled if memory copy fails to make progress in this time, in seconds.",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_hco_v1beta1_MediatedHostDevice(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "MediatedHostDevice represents a host mediated device allowed for passthrough",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"mdevNameSelector": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"resourceName": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"externalResourceProvider": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"boolean"},
							Format: "",
						},
					},
				},
				Required: []string{"mdevNameSelector", "resourceName"},
			},
		},
	}
}

func schema_pkg_apis_hco_v1beta1_OperandResourceRequirements(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OperandResourceRequirements is a list of resource requirements for the operand workloads pods",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"storageWorkloads": {
						SchemaProps: spec.SchemaProps{
							Description: "StorageWorkloads defines the resources requirements for storage workloads. It will propagate to the CDI custom resource",
							Ref:         ref("k8s.io/api/core/v1.ResourceRequirements"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/api/core/v1.ResourceRequirements"},
	}
}

func schema_pkg_apis_hco_v1beta1_PciHostDevice(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PciHostDevice represents a host PCI device allowed for passthrough",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"pciVendorSelector": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"resourceName": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"externalResourceProvider": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"boolean"},
							Format: "",
						},
					},
				},
				Required: []string{"pciVendorSelector", "resourceName"},
			},
		},
	}
}

func schema_pkg_apis_hco_v1beta1_PermittedHostDevices(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PermittedHostDevices holds information about devices allowed for passthrough",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"pciHostDevices": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "atomic",
							},
						},
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.PciHostDevice"),
									},
								},
							},
						},
					},
					"mediatedDevices": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "atomic",
							},
						},
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.MediatedHostDevice"),
									},
								},
							},
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.MediatedHostDevice", "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1.PciHostDevice"},
	}
}
