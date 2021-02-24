package v1beta1

import (
	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	sdkapi "kubevirt.io/controller-lifecycle-operator-sdk/pkg/sdk/api"
	"reflect"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file

// HyperConvergedName is the name of the HyperConverged resource that will be reconciled
const HyperConvergedName = "kubevirt-hyperconverged"

// HyperConvergedSpec defines the desired state of HyperConverged
// +k8s:openapi-gen=true
type HyperConvergedSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// LocalStorageClassName the name of the local storage class.
	LocalStorageClassName string `json:"localStorageClassName,omitempty"`

	// infra HyperConvergedConfig influences the pod configuration (currently only placement)
	// for all the infra components needed on the virtualization enabled cluster
	// but not necessarely directly on each node running VMs/VMIs.
	// +optional
	Infra HyperConvergedConfig `json:"infra,omitempty"`

	// workloads HyperConvergedConfig influences the pod configuration (currently only placement) of components
	// which need to be running on a node where virtualization workloads should be able to run.
	// Changes to Workloads HyperConvergedConfig can be applied only without existing workload.
	// +optional
	Workloads HyperConvergedConfig `json:"workloads,omitempty"`

	// featureGates is a map of feature gate flags. Setting a flag to `true` will enable
	// the feature. Setting `false` or removing the feature gate, disables the feature.
	// +optional
	// +TODO: Always keep the default FeatureGates in sync with the default field values in HyperConvergedFeatureGates //NOSONAR
	// +kubebuilder:default={dataVolumes: true, sriov: true, liveMigration: true, cpuManager: true, cpuNodeDiscovery: true, sidecar: true, snapshot: true, sriovLiveMigration: false, hotplugVolumes: false, gpu: false, hostDevices: false, withHostPassthroughCPU: false, withHostModelCPU: true, hypervStrictCheck: true}
	FeatureGates *HyperConvergedFeatureGates `json:"featureGates,omitempty"`

	// operator version
	Version string `json:"version,omitempty"`
}

// HyperConvergedConfig defines a set of configurations to pass to components
type HyperConvergedConfig struct {
	// NodePlacement describes node scheduling configuration.
	// +optional
	NodePlacement *sdkapi.NodePlacement `json:"nodePlacement,omitempty"`
}

type FeatureGate *bool

// HyperConvergedFeatureGates is a set of optional feature gates to enable or disable new features that are not enabled
// by default yet.
// +optional
// +k8s:openapi-gen=true
// +kubebuilder:default={}
type HyperConvergedFeatureGates struct {
	// indicates that we support turning on DataVolume workflows. This means using DataVolumes in the VM and VMI
	// definitions. There was a period of time where this was in alpha and needed to be explicility enabled.
	// It also means that someone is using KubeVirt with CDI. So by not enabling this feature gate, someone can safely
	// use kubevirt without CDI and know that users of kubevirt will not be able to post VM/VMIs that use CDI workflows
	// that aren't available to them
	// +optional
	// +kubebuilder:default=true
	DataVolumes FeatureGate `json:"dataVolumes,omitempty"`

	// TODO add description
	// +optional
	// +kubebuilder:default=true
	SRIOV FeatureGate `json:"sriov,omitempty"`

	// Enables VMIs to be live migrated. Without this, migrations are not possible and will be blocked
	// +optional
	// +kubebuilder:default=true
	LiveMigration FeatureGate `json:"liveMigration,omitempty"`

	// TODO add description
	// +optional
	// +kubebuilder:default=true
	CPUManager FeatureGate `json:"cpuManager,omitempty"`

	// TODO add description
	// +optional
	// +kubebuilder:default=true
	CPUNodeDiscovery FeatureGate `json:"cpuNodeDiscovery,omitempty"`

	// Enables using our sidecar hooks for injecting custom logic into the VMI startup flow. This is a very advanced
	// feature that has security implications, which is why it is opt-in only
	// +optional
	// +kubebuilder:default=true
	Sidecar FeatureGate `json:"sidecar,omitempty"`

	// Enables the alpha offline snapshot functionality
	// +optional
	// +kubebuilder:default=true
	Snapshot FeatureGate `json:"snapshot,omitempty"`

	// Allow migrating a virtual machine with SRIOV interfaces.
	// When enabled virt-launcher pods of virtual machines with SRIOV
	// interfaces run with CAP_SYS_RESOURCE capability.
	// This may degrade virt-launcher security.
	// +optional
	// +kubebuilder:default=false
	SRIOVLiveMigration FeatureGate `json:"sriovLiveMigration,omitempty"`

	// Allow attaching a data volume to a running VMI
	// +optional
	// +kubebuilder:default=false
	HotplugVolumes FeatureGate `json:"hotplugVolumes,omitempty"`

	// Allow assigning GPU and vGPU devices to virtual machines
	// +optional
	// +kubebuilder:default=false
	GPU FeatureGate `json:"gpu,omitempty"`

	// Allow assigning host devices to virtual machines
	// +optional
	// +kubebuilder:default=false
	HostDevices FeatureGate `json:"hostDevices,omitempty"`

	// Allow migrating a virtual machine with CPU host-passthrough mode. This should be
	// enabled only when the Cluster is homogeneous from CPU HW perspective doc here
	// +optional
	// +kubebuilder:default=false
	WithHostPassthroughCPU FeatureGate `json:"withHostPassthroughCPU,omitempty"`

	// Support migration for VMs with host-model CPU mode
	// +optional
	// +kubebuilder:default=true
	WithHostModelCPU FeatureGate `json:"withHostModelCPU,omitempty"`

	// Enable HyperV strict host checking for HyperV enlightenments
	// Defaults to true, even when HyperConvergedFeatureGates is empty
	// +optional
	// +kubebuilder:default=true
	HypervStrictCheck FeatureGate `json:"hypervStrictCheck,omitempty"`

	// +optional
	// +k8s:openapi-gen=false
	enabledGateList map[string]bool `json:"-"`
}

func (fgs *HyperConvergedFeatureGates) RebuildEnabledGateMap() map[string]bool {
	if fgs == nil {
		return nil
	}

	e := reflect.ValueOf(fgs).Elem()
	t := e.Type()

	fgs.enabledGateList = make(map[string]bool)

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := e.Field(i)

		if fieldType.Type.Name() == "FeatureGate" && !fieldValue.IsNil() && fieldValue.Elem().Bool() {
			fgs.enabledGateList[fieldType.Name] = true
		}
	}

	return fgs.enabledGateList
}

func (fgs *HyperConvergedFeatureGates) GetEnabledGateMap() map[string]bool {
	if fgs == nil {
		return nil
	}

	return fgs.enabledGateList
}

// HyperConvergedStatus defines the observed state of HyperConverged
// +k8s:openapi-gen=true
type HyperConvergedStatus struct {
	// Conditions describes the state of the HyperConverged resource.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []conditionsv1.Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`

	// RelatedObjects is a list of objects created and maintained by this
	// operator. Object references will be added to this list after they have
	// been created AND found in the cluster.
	// +optional
	RelatedObjects []corev1.ObjectReference `json:"relatedObjects,omitempty"`

	// Versions is a list of HCO component versions, as name/version pairs. The version with a name of "operator"
	// is the HCO version itself, as described here:
	// https://github.com/openshift/cluster-version-operator/blob/master/docs/dev/clusteroperator.md#version
	// +optional
	Versions Versions `json:"versions,omitempty"`
}

func (hcs *HyperConvergedStatus) UpdateVersion(name, version string) {
	if hcs.Versions == nil {
		hcs.Versions = Versions{}
	}
	hcs.Versions.updateVersion(name, version)
}

func (hcs *HyperConvergedStatus) GetVersion(name string) (string, bool) {
	return hcs.Versions.getVersion(name)
}

type Version struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

func newVersion(name, version string) Version {
	return Version{Name: name, Version: version}
}

type Versions []Version

func (vs *Versions) updateVersion(name, version string) {
	for i, v := range *vs {
		if v.Name == name {
			(*vs)[i].Version = version
			return
		}
	}
	*vs = append(*vs, newVersion(name, version))
}

func (vs *Versions) getVersion(name string) (string, bool) {
	for _, v := range *vs {
		if v.Name == name {
			return v.Version, true
		}
	}
	return "", false
}

const (

	// ConditionReconcileComplete communicates the status of the HyperConverged resource's
	// reconcile functionality. Basically, is the Reconcile function running to completion.
	ConditionReconcileComplete conditionsv1.ConditionType = "ReconcileComplete"

	// ConditionTaintedConfiguration indicates that a hidden/debug configuration
	// has been applied to the HyperConverged resource via a specialized annotation.
	// This condition is exposed only when its value is True, and is otherwise hidden.
	ConditionTaintedConfiguration conditionsv1.ConditionType = "TaintedConfiguration"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HyperConverged is the Schema for the hyperconvergeds API
// +k8s:openapi-gen=true
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:resource:scope=Namespaced,categories={all},shortName={hco,hcos}
// +kubebuilder:subresource:status
type HyperConverged struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HyperConvergedSpec   `json:"spec,omitempty"`
	Status HyperConvergedStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HyperConvergedList contains a list of HyperConverged
type HyperConvergedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HyperConverged `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HyperConverged{}, &HyperConvergedList{})
}
