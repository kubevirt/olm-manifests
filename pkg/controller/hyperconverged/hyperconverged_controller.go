package hyperconverged

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/kubevirt/hyperconverged-cluster-operator/pkg/controller/common"
	"github.com/kubevirt/hyperconverged-cluster-operator/pkg/controller/operands"
	
	"github.com/operator-framework/operator-sdk/pkg/ready"
	schedulingv1 "k8s.io/api/scheduling/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	networkaddonsv1 "github.com/kubevirt/cluster-network-addons-operator/pkg/apis/networkaddonsoperator/v1"
	hcov1beta1 "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1beta1"
	hcoutil "github.com/kubevirt/hyperconverged-cluster-operator/pkg/util"
	version "github.com/kubevirt/hyperconverged-cluster-operator/version"
	sspv1 "github.com/kubevirt/kubevirt-ssp-operator/pkg/apis/kubevirt/v1"
	vmimportv1beta1 "github.com/kubevirt/vm-import-operator/pkg/apis/v2v/v1beta1"
	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	objectreferencesv1 "github.com/openshift/custom-resource-status/objectreferences/v1"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kubevirtv1 "kubevirt.io/client-go/api/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1beta1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	log = logf.Log.WithName("controller_hyperconverged")
)

const (
	// We cannot set owner reference of cluster-wide resources to namespaced HyperConverged object. Therefore,
	// use finalizers to manage the cleanup.
	FinalizerName = "hyperconvergeds.hco.kubevirt.io"

	// OpenshiftNamespace is for resources that belong in the openshift namespace

	reconcileInit               = "Init"
	reconcileInitMessage        = "Initializing HyperConverged cluster"
	reconcileFailed             = "ReconcileFailed"
	reconcileCompleted          = "ReconcileCompleted"
	reconcileCompletedMessage   = "Reconcile completed successfully"
	invalidRequestReason        = "InvalidRequest"
	invalidRequestMessageFormat = "Request does not match expected name (%v) and namespace (%v)"
	commonDegradedReason        = "HCODegraded"
	commonProgressingReason     = "HCOProgressing"

	ErrCDIUninstall       = "ErrCDIUninstall"
	uninstallCDIErrorMsg  = "The uninstall request failed on CDI component: "
	ErrVirtUninstall      = "ErrVirtUninstall"
	uninstallVirtErrorMsg = "The uninstall request failed on virt component: "
	ErrHCOUninstall       = "ErrHCOUninstall"
	uninstallHCOErrorMsg  = "The uninstall request failed on dependent components, please check their logs."

	hcoVersionName = "operator"

	commonTemplatesBundleOldCrdName = "kubevirtcommontemplatesbundles.kubevirt.io"
	metricsAggregationOldCrdName    = "kubevirtmetricsaggregations.kubevirt.io"
	nodeLabellerBundlesOldCrdName   = "kubevirtnodelabellerbundles.kubevirt.io"
	templateValidatorsOldCrdName    = "kubevirttemplatevalidators.kubevirt.io"
)

// Add creates a new HyperConverged Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, ci hcoutil.ClusterInfo) error {
	return add(mgr, newReconciler(mgr, ci))
}

// temp map, until we move all the operands code
var operandMap = map[string]operands.Operand{}
func prepareHandlerMap(clt client.Client, scheme *runtime.Scheme) {
	operandMap["kvc"] = &operands.KvConfigHandler{Client: clt, Scheme: scheme}
	operandMap["kvpc"] = &operands.KvPriorityClassHandler{Client: clt, Scheme: scheme}
	operandMap["kv"] = &operands.KubevirtHandler{Client: clt, Scheme: scheme}
	operandMap["cdi"] = &operands.CdiHandler{Client: clt, Scheme: scheme}
	operandMap["cna"] = &operands.CnaHandler{Client: clt, Scheme: scheme}
	operandMap["vmimport"] = &operands.VmImportHandler{Client: clt, Scheme: scheme}
}


// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, ci hcoutil.ClusterInfo) reconcile.Reconciler {

	ownVersion := os.Getenv(hcoutil.HcoKvIoVersionName)
	if ownVersion == "" {
		ownVersion = version.Version
	}

	prepareHandlerMap(mgr.GetClient(), mgr.GetScheme())

	return &ReconcileHyperConverged{
		client:      mgr.GetClient(),
		scheme:      mgr.GetScheme(),
		recorder:    mgr.GetEventRecorderFor(hcoutil.HyperConvergedName),
		upgradeMode: false,
		ownVersion:  ownVersion,
		clusterInfo: ci,
		shouldRemoveOldCrd: map[string]bool{
			commonTemplatesBundleOldCrdName: true,
			metricsAggregationOldCrdName:    true,
			nodeLabellerBundlesOldCrdName:   true,
			templateValidatorsOldCrdName:    true,
		},
		eventEmitter: hcoutil.GetEventEmitter(),
		firstLoop:    true,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("hyperconverged-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource HyperConverged
	err = c.Watch(&source.Kind{Type: &hcov1beta1.HyperConverged{}}, &handler.EnqueueRequestForObject{}, predicate.GenerationChangedPredicate{})
	if err != nil {
		return err
	}

	hco, err := getHyperconverged()
	if err != nil {
		return err
	}

	// Watch secondary resources
	for _, resource := range []runtime.Object{
		&kubevirtv1.KubeVirt{},
		&cdiv1beta1.CDI{},
		&networkaddonsv1.NetworkAddonsConfig{},
		&sspv1.KubevirtCommonTemplatesBundle{},
		&sspv1.KubevirtNodeLabellerBundle{},
		&sspv1.KubevirtTemplateValidator{},
		&sspv1.KubevirtMetricsAggregation{},
		&schedulingv1.PriorityClass{},
		&vmimportv1beta1.VMImportConfig{},
	} {
		err = c.Watch(&source.Kind{Type: resource}, &handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(
				// always enqueue the same HyperConverged object, since there should be only one
				func(a handler.MapObject) []reconcile.Request {
					return []reconcile.Request{
						{NamespacedName: hco},
					}
				}),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileHyperConverged{}

// ReconcileHyperConverged reconciles a HyperConverged object
type ReconcileHyperConverged struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client             client.Client
	scheme             *runtime.Scheme
	recorder           record.EventRecorder
	upgradeMode        bool
	ownVersion         string
	clusterInfo        hcoutil.ClusterInfo
	shouldRemoveOldCrd map[string]bool
	eventEmitter       hcoutil.EventEmitter
	firstLoop          bool
}

// Reconcile reads that state of the cluster for a HyperConverged object and makes changes based on the state read
// and what is in the HyperConverged.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileHyperConverged) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	req := common.NewHcoRequest(request, log, r.upgradeMode)
	req.Logger.Info("Reconciling HyperConverged operator")

	// Fetch the HyperConverged instance
	instance, err := r.getHcoInstanceFromK8s(req)
	if instance == nil {
		return reconcile.Result{}, err
	}
	req.Instance = instance

	if r.firstLoop {
		// reload eventEmitter. The client should now find all the required resources
		r.eventEmitter.UpdateClient(req.Ctx, r.client, req.Logger)
	}

	res, err := r.doReconcile(req)
	if r.firstLoop {
		r.firstLoop = false
	}

	if err != nil {
		r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeWarning, "ReconcileError", err.Error())
	}

	/*
		From K8s API reference: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/
		============================================================================================================
		Replace: Replacing a resource object will update the resource by replacing the existing spec with the
		provided one. For read-then-write operations this is safe because an optimistic lock failure will occur if
		the resource was modified between the read and write.

		**Note: The ResourceStatus will be ignored by the system and will not be updated. To update the status, one
		must invoke the specific status update operation.**
		============================================================================================================

		In addition, updating the status should not update the metadata, so we need to update both the CR and the
		CR Status, and we need to update the status first, in order to prevent a conflict.
	*/

	if req.StatusDirty {
		updateErr := r.client.Status().Update(req.Ctx, req.Instance)
		if updateErr != nil {
			updateErrorMsg := "Failed to update HCO Status"
			r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeWarning, "HcoUpdateError", updateErrorMsg)
			req.Logger.Error(updateErr, updateErrorMsg)
			err = updateErr
		}
	}

	// recover Spec.Version if upgrade missed when upgrade completed
	// Doing it here because status.update overrides spec for some reason
	knownHcoVersion, versionFound := req.Instance.Status.GetVersion(hcoVersionName)
	if (!r.upgradeMode) && versionFound && (knownHcoVersion == r.ownVersion) && (req.Instance.Spec.Version != r.ownVersion) {
		req.Instance.Spec.Version = r.ownVersion
		req.Dirty = true
	}

	if req.Dirty {
		updateErr := r.client.Update(req.Ctx, req.Instance)
		if updateErr != nil {
			updateErrorMsg := "Failed to update HCO CR"
			r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeWarning, "HcoUpdateError", updateErrorMsg)
			req.Logger.Error(updateErr, updateErrorMsg)
			err = updateErr
		}
	}

	if apierrors.IsConflict(err) {
		res.Requeue = true
	}

	return res, err
}

func (r *ReconcileHyperConverged) doReconcile(req *common.HcoRequest) (reconcile.Result, error) {

	valid, err := r.validateNamespace(req)
	if !valid {
		return reconcile.Result{}, err
	}
	// Add conditions if there are none
	init := req.Instance.Status.Conditions == nil
	if init {
		r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeNormal, "InitHCO", "Initiating the HyperConverged")
		err = r.setInitialConditions(req)
		if err != nil {
			req.Logger.Error(err, "Failed to add conditions to status")
			return reconcile.Result{}, err
		}
	}

	r.setLabels(req)

	// in-memory conditions should start off empty. It will only ever hold
	// negative conditions (!Available, Degraded, Progressing)
	req.Conditions = common.NewHcoConditions()

	// Handle finalizers
	if req.Instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// Add the finalizer if it's not there
		if !contains(req.Instance.ObjectMeta.Finalizers, FinalizerName) {
			req.Instance.ObjectMeta.Finalizers = append(req.Instance.ObjectMeta.Finalizers, FinalizerName)
			req.Dirty = true
		}
	} else {
		if contains(req.Instance.ObjectMeta.Finalizers, FinalizerName) {
			return r.ensureHcoDeleted(req)
		}
	}

	// If the current version is not updated in CR ,then we're updating. This is also works when updating from
	// an old version, since Status.Versions will be empty.
	knownHcoVersion, _ := req.Instance.Status.GetVersion(hcoVersionName)

	if !r.upgradeMode && !init && knownHcoVersion != r.ownVersion {
		r.upgradeMode = true
		r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeNormal, "UpgradeHCO", "Upgrading the HyperConverged to version "+r.ownVersion)
		req.Logger.Info(fmt.Sprintf("Start upgrating from version %s to version %s", knownHcoVersion, r.ownVersion))
	}

	req.SetUpgradeMode(r.upgradeMode)

	r.ensureConsoleCLIDownload(req)

	err = r.ensureHco(req)
	if err != nil {
		return reconcile.Result{}, r.updateConditions(req)
	}

	req.Logger.Info("Reconcile complete")

	// Requeue if we just created everything
	if init {
		return reconcile.Result{Requeue: true}, err
	}

	err = r.completeReconciliation(req)

	return reconcile.Result{}, err
}

func (r *ReconcileHyperConverged) getHcoInstanceFromK8s(req *common.HcoRequest) (*hcov1beta1.HyperConverged, error) {
	instance := &hcov1beta1.HyperConverged{}
	err := r.client.Get(req.Ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			req.Logger.Info("No HyperConverged resource")
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return nil, nil
		}
		// Error reading the object - requeue the request.
		return nil, err
	}
	return instance, nil
}

func (r *ReconcileHyperConverged) validateNamespace(req *common.HcoRequest) (bool, error) {
	hco, err := getHyperconverged()
	if err != nil {
		req.Logger.Error(err, "Failed to get HyperConverged namespaced name")
		return false, err
	}

	// Ignore invalid requests
	if req.NamespacedName != hco {
		req.Logger.Info("Invalid request", "HyperConverged.Namespace", hco.Namespace, "HyperConverged.Name", hco.Name)
		req.Conditions.SetStatusCondition(conditionsv1.Condition{
			Type:    hcov1beta1.ConditionReconcileComplete,
			Status:  corev1.ConditionFalse,
			Reason:  invalidRequestReason,
			Message: fmt.Sprintf(invalidRequestMessageFormat, hco.Name, hco.Namespace),
		})
		err := r.updateConditions(req)
		return false, err
	}
	return true, nil
}

func (r *ReconcileHyperConverged) setInitialConditions(req *common.HcoRequest) error {
	req.Instance.Status.UpdateVersion(hcoVersionName, r.ownVersion)
	req.Instance.Spec.Version = r.ownVersion
	req.Dirty = true

	req.Conditions.SetStatusCondition(conditionsv1.Condition{
		Type:    hcov1beta1.ConditionReconcileComplete,
		Status:  corev1.ConditionUnknown, // we just started trying to reconcile
		Reason:  reconcileInit,
		Message: reconcileInitMessage,
	})
	req.Conditions.SetStatusCondition(conditionsv1.Condition{
		Type:    conditionsv1.ConditionAvailable,
		Status:  corev1.ConditionFalse,
		Reason:  reconcileInit,
		Message: reconcileInitMessage,
	})
	req.Conditions.SetStatusCondition(conditionsv1.Condition{
		Type:    conditionsv1.ConditionProgressing,
		Status:  corev1.ConditionTrue,
		Reason:  reconcileInit,
		Message: reconcileInitMessage,
	})
	req.Conditions.SetStatusCondition(conditionsv1.Condition{
		Type:    conditionsv1.ConditionDegraded,
		Status:  corev1.ConditionFalse,
		Reason:  reconcileInit,
		Message: reconcileInitMessage,
	})
	req.Conditions.SetStatusCondition(conditionsv1.Condition{
		Type:    conditionsv1.ConditionUpgradeable,
		Status:  corev1.ConditionUnknown,
		Reason:  reconcileInit,
		Message: reconcileInitMessage,
	})

	return r.updateConditions(req)
}

func (r *ReconcileHyperConverged) ensureHcoDeleted(req *common.HcoRequest) (reconcile.Result, error) {
	for _, obj := range []runtime.Object{
		req.Instance.NewKubeVirt(),
		req.Instance.NewCDI(),
		req.Instance.NewNetworkAddons(),
		req.Instance.NewKubeVirtCommonTemplateBundle(),
		req.Instance.NewConsoleCLIDownload(),
		operands.NewVMImportForCR(req.Instance),
	} {
		err := hcoutil.EnsureDeleted(req.Ctx, r.client, obj, req.Instance.Name, req.Logger, false)
		if err != nil {
			req.Logger.Error(err, "Failed to manually delete objects")

			// TODO: ask to other components to expose something like
			// func IsDeleteRefused(err error) bool
			// to be able to clearly distinguish between an explicit
			// refuse from other operator and any other kind of error that
			// could potentially happen in the process

			errT := ErrHCOUninstall
			errMsg := uninstallHCOErrorMsg
			switch obj.(type) {
			case *kubevirtv1.KubeVirt:
				errT = ErrVirtUninstall
				errMsg = uninstallVirtErrorMsg + err.Error()
			case *cdiv1beta1.CDI:
				errT = ErrCDIUninstall
				errMsg = uninstallCDIErrorMsg + err.Error()
			}

			r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeWarning, errT, errMsg)

			return reconcile.Result{}, err
		}

		if key, err := client.ObjectKeyFromObject(obj); err == nil {
			r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeNormal, "Killing", fmt.Sprintf("Removed %s %s", obj.GetObjectKind().GroupVersionKind().Kind, key.Name))
		}
	}

	// Remove the finalizer
	req.Instance.ObjectMeta.Finalizers = drop(req.Instance.ObjectMeta.Finalizers, FinalizerName)
	req.Dirty = true

	// Need to requeue because finalizer update does not change metadata.generation
	return reconcile.Result{Requeue: true}, nil
}

func (r *ReconcileHyperConverged) ensureHco(req *common.HcoRequest) error {
	for _, f := range []func(*common.HcoRequest) *operands.EnsureResult{
		r.ensureKubeVirtPriorityClass,
		r.ensureKubeVirtConfig,
		r.ensureKubeVirt,
		r.ensureCDI,
		r.ensureNetworkAddons,
		r.ensureKubeVirtCommonTemplateBundle,
		r.ensureKubeVirtNodeLabellerBundle,
		r.ensureKubeVirtTemplateValidator,
		r.ensureKubeVirtMetricsAggregation,
		r.ensureIMSConfig,
		r.ensureVMImport,
	} {
		res := f(req)
		if res.Err != nil {
			req.ComponentUpgradeInProgress = false
			req.Conditions.SetStatusCondition(conditionsv1.Condition{
				Type:    hcov1beta1.ConditionReconcileComplete,
				Status:  corev1.ConditionFalse,
				Reason:  reconcileFailed,
				Message: fmt.Sprintf("Error while reconciling: %v", res.Err),
			})
			return res.Err
		}

		if res.Created {
			r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeNormal, "Created", fmt.Sprintf("Created %s %s", res.Type, res.Name))
		} else if res.Updated {
			r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeNormal, "Updated", fmt.Sprintf("Updated %s %s", res.Type, res.Name))
		}

		req.ComponentUpgradeInProgress = req.ComponentUpgradeInProgress && res.UpgradeDone
	}
	return nil
}

func (r *ReconcileHyperConverged) aggregateComponentConditions(req *common.HcoRequest) bool {
	/*
		See the chart at design/aggregateComponentConditions.svg; The numbers below follows the numbers in the chart
		Here is the PlantUML code for the chart that describes the aggregation of the sub-components conditions.
		Find the PlantURL syntax here: https://plantuml.com/activity-diagram-beta

		@startuml ../../../design/aggregateComponentConditions.svg
		title Aggregate Component Conditions

		start
		  #springgreen:Set **ReconcileComplete = True**]
		  !x=1
		if ((x) [Degraded = True] Exists) then
		  !x=x+1
		  #orangered:<<implicit>>\n**Degraded = True** /
		  -[#orangered]-> yes;
		  if ((x) [Progressing = True] Exists) then
			!x=x+1
			-[#springgreen]-> no;
			#springgreen:(x) Set **Progressing = False**]
			!x=x+1
		  else
			-[#orangered]-> yes;
			#orangered:<<implicit>>\n**Progressing = True** /
		  endif
		  if ((x) [Upgradable = False] Exists) then
			!x=x+1
			-[#springgreen]-> no;
			#orangered:(x) Set **Upgradable = False**]
			!x=x+1
		  else
			-[#orangered]-> yes;
			#orangered:<<implicit>>\n**Upgradable = False** /
		  endif
		  if ((x) [Available = False] Exists) then
			!x=x+1
			-[#springgreen]-> no;
			#orangered:(x) Set **Available = False**]
			!x=x+1
		  else
			-[#orangered]-> yes;
			#orangered:<<implicit>>\n**Available = False** /
		  endif
		else
		  -[#springgreen]-> no;
		  #springgreen:(x) Set **Degraded = False**]
		  !x=x+1
		  if ((x) [Progressing = True] Exists) then
			!x=x+1
			-[#orangered]-> yes;
			#orangered:<<implicit>>\n**Progressing = True** /
			if ((x) [Upgradable = False] Exists) then
			  !x=x+1
			  -[#springgreen]-> no;
			  #orangered:(x) Set **Upgradable = False**]
			  !x=x+1
			else
			  -[#orangered]-> yes;
			  #orangered:<<implicit>>\n**Upgradable = False** /
			endif
			if ((x) [Available = False] Exists) then
			  !x=x+1
			  -[#springgreen]-> no;
			  #springgreen:(x) Set **Available = True**]
			  !x=x+1
			else
			  #orangered:<<implicit>>\n**Available = False** /
			  -[#orangered]-> yes;
			endif
		  else
			-[#springgreen]-> no;
			#springgreen:(x) Set **Progressing = False**]
			!x=x+1
			if ((x) [Upgradable = False] Exists) then
			  !x=x+1
			  -[#springgreen]-> no;
			  #springgreen:(x) Set **Upgradable = True**]
			  !x=x+1
			else
			#orangered:<<implicit>>\n**Upgradable = False** /
			  -[#orangered]-> yes;
			endif
			if ((x) [Available = False] Exists) then
			  !x=x+1
			  -[#springgreen]-> no;
			  #springgreen:(x) Set **Available = True**]
			  !x=x+1
			else
			  -[#orangered]-> yes;
			  #orangered:<<implicit>>\n**Available = False** /
			endif
		  endif
		endif
		end
		@enduml
	*/

	/*
		    If any component operator reports negatively we want to write that to
			the instance while preserving it's lastTransitionTime.
			For example, consider the KubeVirt resource has the Available condition
			type with type "False". When reconciling KubeVirt's resource we would
			add it to the in-memory representation of HCO's conditions (r.conditions)
			and here we are simply writing it back to the server.
			One shortcoming is that only one failure of a particular condition can be
			captured at one time (ie. if KubeVirt and CDI are both reporting !Available,
		    you will only see CDI as it updates last).
	*/
	allComponentsAreUp := req.Conditions.Empty()
	req.Conditions.SetStatusCondition(conditionsv1.Condition{
		Type:    hcov1beta1.ConditionReconcileComplete,
		Status:  corev1.ConditionTrue,
		Reason:  reconcileCompleted,
		Message: reconcileCompletedMessage,
	})

	if _, conditionFound := req.Conditions[conditionsv1.ConditionDegraded]; conditionFound { // (#chart 1)
		if _, conditionFound = req.Conditions[conditionsv1.ConditionProgressing]; !conditionFound { // (#chart 2)
			req.Conditions.SetStatusCondition(conditionsv1.Condition{ // (#chart 3)
				Type:    conditionsv1.ConditionProgressing,
				Status:  corev1.ConditionFalse,
				Reason:  reconcileCompleted,
				Message: reconcileCompletedMessage,
			})
		} // else - Progressing is already exists

		if _, conditionFound = req.Conditions[conditionsv1.ConditionUpgradeable]; !conditionFound { // (#chart 4)
			req.Conditions.SetStatusCondition(conditionsv1.Condition{ // (#chart 5)
				Type:    conditionsv1.ConditionUpgradeable,
				Status:  corev1.ConditionFalse,
				Reason:  commonDegradedReason,
				Message: "HCO is not Upgradeable due to degraded components",
			})
		} // else - Upgradeable is already exists
		if _, conditionFound = req.Conditions[conditionsv1.ConditionAvailable]; !conditionFound { // (#chart 6)
			req.Conditions.SetStatusCondition(conditionsv1.Condition{ // (#chart 7)
				Type:    conditionsv1.ConditionAvailable,
				Status:  corev1.ConditionFalse,
				Reason:  commonDegradedReason,
				Message: "HCO is not available due to degraded components",
			})
		} // else - Available is already exists
	} else {
		// Degraded is not found. add it.
		req.Conditions.SetStatusCondition(conditionsv1.Condition{ // (#chart 8)
			Type:    conditionsv1.ConditionDegraded,
			Status:  corev1.ConditionFalse,
			Reason:  reconcileCompleted,
			Message: reconcileCompletedMessage,
		})

		if _, conditionFound = req.Conditions[conditionsv1.ConditionProgressing]; conditionFound { // (#chart 9)

			if _, conditionFound = req.Conditions[conditionsv1.ConditionUpgradeable]; !conditionFound { // (#chart 10)
				req.Conditions.SetStatusCondition(conditionsv1.Condition{ // (#chart 11)
					Type:    conditionsv1.ConditionUpgradeable,
					Status:  corev1.ConditionFalse,
					Reason:  commonProgressingReason,
					Message: "HCO is not Upgradeable due to progressing components",
				})
			} // else - Upgradeable is already exists

			if _, conditionFound = req.Conditions[conditionsv1.ConditionAvailable]; !conditionFound { // (#chart 12)
				req.Conditions.SetStatusCondition(conditionsv1.Condition{ // (#chart 13)
					Type:    conditionsv1.ConditionAvailable,
					Status:  corev1.ConditionTrue,
					Reason:  reconcileCompleted,
					Message: reconcileCompletedMessage,
				})
			} // else - Available is already exists
		} else {
			req.Conditions.SetStatusCondition(conditionsv1.Condition{ // (#chart 14)
				Type:    conditionsv1.ConditionProgressing,
				Status:  corev1.ConditionFalse,
				Reason:  reconcileCompleted,
				Message: reconcileCompletedMessage,
			})

			if _, conditionFound = req.Conditions[conditionsv1.ConditionUpgradeable]; !conditionFound { // (#chart 15)
				req.Conditions.SetStatusCondition(conditionsv1.Condition{ // (#chart 16)
					Type:    conditionsv1.ConditionUpgradeable,
					Status:  corev1.ConditionTrue,
					Reason:  reconcileCompleted,
					Message: reconcileCompletedMessage,
				})
			}

			if _, conditionFound = req.Conditions[conditionsv1.ConditionAvailable]; !conditionFound { // (#chart 17) {
				req.Conditions.SetStatusCondition(conditionsv1.Condition{ // (#chart 18)
					Type:    conditionsv1.ConditionAvailable,
					Status:  corev1.ConditionTrue,
					Reason:  reconcileCompleted,
					Message: reconcileCompletedMessage,
				})
			}
		}
	}
	return allComponentsAreUp
}

func (r *ReconcileHyperConverged) completeReconciliation(req *common.HcoRequest) error {
	allComponentsAreUp := r.aggregateComponentConditions(req)

	hcoReady := false

	if allComponentsAreUp {
		req.Logger.Info("No component operator reported negatively")

		// if in upgrade mode, and all the components are upgraded - upgrade is completed
		if r.upgradeMode && req.ComponentUpgradeInProgress {
			// update the new version only when upgrade is completed
			req.Instance.Status.UpdateVersion(hcoVersionName, r.ownVersion)
			req.StatusDirty = true

			req.Instance.Spec.Version = r.ownVersion
			req.Dirty = true

			r.upgradeMode = false
			req.ComponentUpgradeInProgress = false
			req.Logger.Info(fmt.Sprintf("Successfuly upgraded to version %s", r.ownVersion))
			r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeNormal, "UpgradeHCO", fmt.Sprintf("Successfuly upgraded to version %s", r.ownVersion))
		}

		// If not in upgrade mode, then we're ready, because all the operators reported positive conditions.
		// if upgrade was done successfully, r.upgradeMode is already false here.
		hcoReady = !r.upgradeMode
	}

	if r.upgradeMode {
		// override the Progressing condition during upgrade
		req.Conditions.SetStatusCondition(conditionsv1.Condition{
			Type:    conditionsv1.ConditionProgressing,
			Status:  corev1.ConditionTrue,
			Reason:  "HCOUpgrading",
			Message: "HCO is now upgrading to version " + r.ownVersion,
		})
	}

	fr := ready.NewFileReady()
	if hcoReady {
		// If no operator whose conditions we are watching reports an error, then it is safe
		// to set readiness.
		r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeNormal, "ReconcileHCO", "HCO Reconcile completed successfully")
		err := fr.Set()
		if err != nil {
			req.Logger.Error(err, "Failed to mark operator ready")
			return err
		}
	} else {
		// If for any reason we marked ourselves !upgradeable...then unset readiness
		if r.upgradeMode {
			r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeNormal, "ReconcileHCO", "HCO Upgrade in progress")
		} else {
			r.eventEmitter.EmitEvent(req.Instance, corev1.EventTypeWarning, "ReconcileHCO", "Not all the operators are ready")
		}
		err := fr.Unset()
		if err != nil {
			req.Logger.Error(err, "Failed to mark operator unready")
			return err
		}
	}
	return r.updateConditions(req)
}

func (r *ReconcileHyperConverged) checkComponentVersion(versionEnvName, actualVersion string) bool {
	expectedVersion := os.Getenv(versionEnvName)
	return expectedVersion != "" && expectedVersion == actualVersion
}

func (r *ReconcileHyperConverged) ensureKubeVirtPriorityClass(req *common.HcoRequest) *operands.EnsureResult {
	return operandMap["kvpc"].Ensure(req)
}

func (r *ReconcileHyperConverged) ensureKubeVirtConfig(req *common.HcoRequest) *operands.EnsureResult {
	return operandMap["kvc"].Ensure(req)
}

func (r *ReconcileHyperConverged) ensureKubeVirt(req *common.HcoRequest) *operands.EnsureResult {
	return operandMap["kv"].Ensure(req)
}

func (r *ReconcileHyperConverged) ensureCDI(req *common.HcoRequest) *operands.EnsureResult {
	return operandMap["cdi"].Ensure(req)
}

func (r *ReconcileHyperConverged) ensureNetworkAddons(req *common.HcoRequest) *operands.EnsureResult {
	return operandMap["cna"].Ensure(req)
}

// handleComponentConditions - read and process a sub-component conditions.
// returns true if the the conditions indicates "ready" state and false if not.
func handleComponentConditions(r *ReconcileHyperConverged, req *common.HcoRequest, component string, componentConds []conditionsv1.Condition) (isReady bool) {
	isReady = true
	if len(componentConds) == 0 {
		isReady = false
		reason := fmt.Sprintf("%sConditions", component)
		message := fmt.Sprintf("%s resource has no conditions", component)
		req.Logger.Info(fmt.Sprintf("%s's resource is not reporting Conditions on it's Status", component))
		req.Conditions.SetStatusCondition(conditionsv1.Condition{
			Type:    conditionsv1.ConditionAvailable,
			Status:  corev1.ConditionFalse,
			Reason:  reason,
			Message: message,
		})
		req.Conditions.SetStatusCondition(conditionsv1.Condition{
			Type:    conditionsv1.ConditionProgressing,
			Status:  corev1.ConditionTrue,
			Reason:  reason,
			Message: message,
		})
		req.Conditions.SetStatusCondition(conditionsv1.Condition{
			Type:    conditionsv1.ConditionUpgradeable,
			Status:  corev1.ConditionFalse,
			Reason:  reason,
			Message: message,
		})
	} else {
		foundAvailableCond := false
		foundProgressingCond := false
		foundDegradedCond := false
		for _, condition := range componentConds {
			switch condition.Type {
			case conditionsv1.ConditionAvailable:
				foundAvailableCond = true
				if condition.Status == corev1.ConditionFalse {
					isReady = false
					msg := fmt.Sprintf("%s is not available: %v", component, string(condition.Message))
					r.componentNotAvailable(req, component, msg)
				}
			case conditionsv1.ConditionProgressing:
				foundProgressingCond = true
				if condition.Status == corev1.ConditionTrue {
					isReady = false
					req.Logger.Info(fmt.Sprintf("%s is 'Progressing'", component))
					req.Conditions.SetStatusCondition(conditionsv1.Condition{
						Type:    conditionsv1.ConditionProgressing,
						Status:  corev1.ConditionTrue,
						Reason:  fmt.Sprintf("%sProgressing", component),
						Message: fmt.Sprintf("%s is progressing: %v", component, string(condition.Message)),
					})
					req.Conditions.SetStatusCondition(conditionsv1.Condition{
						Type:    conditionsv1.ConditionUpgradeable,
						Status:  corev1.ConditionFalse,
						Reason:  fmt.Sprintf("%sProgressing", component),
						Message: fmt.Sprintf("%s is progressing: %v", component, string(condition.Message)),
					})
				}
			case conditionsv1.ConditionDegraded:
				foundDegradedCond = true
				if condition.Status == corev1.ConditionTrue {
					isReady = false
					req.Logger.Info(fmt.Sprintf("%s is 'Degraded'", component))
					req.Conditions.SetStatusCondition(conditionsv1.Condition{
						Type:    conditionsv1.ConditionDegraded,
						Status:  corev1.ConditionTrue,
						Reason:  fmt.Sprintf("%sDegraded", component),
						Message: fmt.Sprintf("%s is degraded: %v", component, string(condition.Message)),
					})
				}
			}
		}

		if !foundAvailableCond {
			r.componentNotAvailable(req, component, `missing "Available" condition`)
		}

		isReady = isReady && foundAvailableCond && foundProgressingCond && foundDegradedCond
	}

	return isReady
}

func (r *ReconcileHyperConverged) componentNotAvailable(req *common.HcoRequest, component string, msg string) {
	req.Logger.Info(fmt.Sprintf("%s is not 'Available'", component))
	req.Conditions.SetStatusCondition(conditionsv1.Condition{
		Type:    conditionsv1.ConditionAvailable,
		Status:  corev1.ConditionFalse,
		Reason:  fmt.Sprintf("%sNotAvailable", component),
		Message: msg,
	})
}

func (r *ReconcileHyperConverged) ensureKubeVirtCommonTemplateBundle(req *common.HcoRequest) *operands.EnsureResult {

	kvCTB := req.Instance.NewKubeVirtCommonTemplateBundle()
	res := operands.NewEnsureResult(kvCTB)
	if !r.clusterInfo.IsOpenshift() { // SSP operators Only supported in OpenShift. Ignore in K8s.
		return res.SetUpgradeDone(true)
	}

	key, err := client.ObjectKeyFromObject(kvCTB)
	if err != nil {
		req.Logger.Error(err, "Failed to get object key for KubeVirt Common Templates Bundle")
	}

	res.SetName(key.Name)
	found := &sspv1.KubevirtCommonTemplatesBundle{}

	err = r.client.Get(req.Ctx, key, found)
	if err != nil {
		if apierrors.IsNotFound(err) {
			req.Logger.Info("Creating KubeVirt Common Templates Bundle")
			err = r.client.Create(req.Ctx, kvCTB)
			if err == nil {
				return res.SetCreated()
			}
		}
		return res.Error(err)
	}

	existingOwners := found.GetOwnerReferences()

	// Previous versions used to have HCO-operator (namespace: kubevirt-hyperconverged)
	// as the owner of kvCTB (namespace: OpenshiftNamespace).
	// It's not legal, so remove that.
	if len(existingOwners) > 0 {
		req.Logger.Info("kvCTB has owners, removing...")
		found.SetOwnerReferences([]metav1.OwnerReference{})
		err = r.client.Update(req.Ctx, found)
		if err != nil {
			req.Logger.Error(err, "Failed to remove kvCTB's previous owners")
		}
	}

	req.Logger.Info("KubeVirt Common Templates Bundle already exists", "bundle.Namespace", found.Namespace, "bundle.Name", found.Name)

	if !reflect.DeepEqual(kvCTB.Spec, found.Spec) {
		req.Logger.Info("Updating existing KubeVirt Common Templates Bundle")
		kvCTB.Spec.DeepCopyInto(&found.Spec)
		err = r.client.Update(req.Ctx, found)
		if err != nil {
			return res.Error(err)
		}
		return res.SetUpdated()
	}

	// Add it to the list of RelatedObjects if found
	objectRef, err := reference.GetReference(r.scheme, found)
	if err != nil {
		return res.Error(err)
	}
	objectreferencesv1.SetObjectReference(&req.Instance.Status.RelatedObjects, *objectRef)

	isReady := handleComponentConditions(r, req, "KubevirtCommonTemplatesBundle", found.Status.Conditions)

	upgradeInProgress := false
	if isReady {
		upgradeInProgress = r.upgradeMode && r.checkComponentVersion(hcoutil.SspVersionEnvV, found.Status.ObservedVersion)
		if (upgradeInProgress || !r.upgradeMode) && r.shouldRemoveOldCrd[commonTemplatesBundleOldCrdName] {
			if r.removeCrd(req, commonTemplatesBundleOldCrdName) {
				r.shouldRemoveOldCrd[commonTemplatesBundleOldCrdName] = false
			}
		}
	}

	return res.SetUpgradeDone(req.ComponentUpgradeInProgress && upgradeInProgress)
}

func newKubeVirtNodeLabellerBundleForCR(cr *hcov1beta1.HyperConverged, namespace string) *sspv1.KubevirtNodeLabellerBundle {
	labels := map[string]string{
		hcoutil.AppLabel: cr.Name,
	}

	spec := sspv1.ComponentSpec{
		// UseKVM: isKVMAvailable(),
	}

	if cr.Spec.Workloads.NodePlacement != nil {
		if cr.Spec.Workloads.NodePlacement.Affinity != nil {
			cr.Spec.Workloads.NodePlacement.Affinity.DeepCopyInto(&spec.Affinity)
		}

		if cr.Spec.Workloads.NodePlacement.NodeSelector != nil {
			spec.NodeSelector = make(map[string]string)
			for k, v := range cr.Spec.Workloads.NodePlacement.NodeSelector {
				spec.NodeSelector[k] = v
			}
		}

		for _, hcoTolr := range cr.Spec.Workloads.NodePlacement.Tolerations {
			nlbTolr := corev1.Toleration{}
			hcoTolr.DeepCopyInto(&nlbTolr)
			spec.Tolerations = append(spec.Tolerations, nlbTolr)
		}
	}

	return &sspv1.KubevirtNodeLabellerBundle{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node-labeller-" + cr.Name,
			Labels:    labels,
			Namespace: namespace,
		},
		Spec: spec,
	}
}

func (r *ReconcileHyperConverged) ensureKubeVirtNodeLabellerBundle(req *common.HcoRequest) *operands.EnsureResult {
	kvNLB := newKubeVirtNodeLabellerBundleForCR(req.Instance, req.Namespace)
	res := operands.NewEnsureResult(kvNLB)
	if !r.clusterInfo.IsOpenshift() { // SSP operators Only supported in OpenShift. Ignore in K8s.
		return res.SetUpgradeDone(true)
	}

	if err := controllerutil.SetControllerReference(req.Instance, kvNLB, r.scheme); err != nil {
		return res.Error(err)
	}

	key, err := client.ObjectKeyFromObject(kvNLB)
	if err != nil {
		req.Logger.Error(err, "Failed to get object key for KubeVirt Node Labeller Bundle")
	}

	res.SetName(key.Name)
	found := &sspv1.KubevirtNodeLabellerBundle{}

	err = r.client.Get(req.Ctx, key, found)
	if err != nil {
		if apierrors.IsNotFound(err) {
			req.Logger.Info("Creating KubeVirt Node Labeller Bundle")
			err = r.client.Create(req.Ctx, kvNLB)
			if err == nil {
				return res.SetCreated()
			}
		}
		return res.Error(err)
	}

	req.Logger.Info("KubeVirt Node Labeller Bundle already exists", "bundle.Namespace", found.Namespace, "bundle.Name", found.Name)

	if !reflect.DeepEqual(kvNLB.Spec, found.Spec) {
		req.Logger.Info("Updating existing KubeVirt Node Labeller Bundle")
		kvNLB.Spec.DeepCopyInto(&found.Spec)
		err = r.client.Update(req.Ctx, found)
		if err != nil {
			return res.Error(err)
		}
		return res.SetUpdated()
	}

	// Add it to the list of RelatedObjects if found
	objectRef, err := reference.GetReference(r.scheme, found)
	if err != nil {
		return res.Error(err)
	}
	objectreferencesv1.SetObjectReference(&req.Instance.Status.RelatedObjects, *objectRef)

	isReady := handleComponentConditions(r, req, "KubevirtNodeLabellerBundle", found.Status.Conditions)

	upgradeInProgress := false
	if isReady {
		upgradeInProgress = r.upgradeMode && r.checkComponentVersion(hcoutil.SspVersionEnvV, found.Status.ObservedVersion)
		if (upgradeInProgress || !r.upgradeMode) && r.shouldRemoveOldCrd[nodeLabellerBundlesOldCrdName] {
			if r.removeCrd(req, nodeLabellerBundlesOldCrdName) {
				r.shouldRemoveOldCrd[nodeLabellerBundlesOldCrdName] = false
			}
		}
	}

	return res.SetUpgradeDone(req.ComponentUpgradeInProgress && upgradeInProgress)
}

func newIMSConfigForCR(cr *hcov1beta1.HyperConverged, namespace string) *corev1.ConfigMap {
	labels := map[string]string{
		hcoutil.AppLabel: cr.Name,
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "v2v-vmware",
			Labels:    labels,
			Namespace: namespace,
		},
		Data: map[string]string{
			"v2v-conversion-image":              os.Getenv("CONVERSION_CONTAINER"),
			"kubevirt-vmware-image":             os.Getenv("VMWARE_CONTAINER"),
			"kubevirt-vmware-image-pull-policy": "IfNotPresent",
		},
	}
}

func (r *ReconcileHyperConverged) ensureIMSConfig(req *common.HcoRequest) *operands.EnsureResult {
	imsConfig := newIMSConfigForCR(req.Instance, req.Namespace)
	res := operands.NewEnsureResult(imsConfig)
	if os.Getenv("CONVERSION_CONTAINER") == "" {
		return res.Error(errors.New("ims-conversion-container not specified"))
	}

	if os.Getenv("VMWARE_CONTAINER") == "" {
		return res.Error(errors.New("ims-vmware-container not specified"))
	}

	err := controllerutil.SetControllerReference(req.Instance, imsConfig, r.scheme)
	if err != nil {
		return res.Error(err)
	}

	key, err := client.ObjectKeyFromObject(imsConfig)
	if err != nil {
		req.Logger.Error(err, "Failed to get object key for IMS Configmap")
	}

	res.SetName(key.Name)
	found := &corev1.ConfigMap{}

	err = r.client.Get(req.Ctx, key, found)
	if err != nil {
		if apierrors.IsNotFound(err) {
			req.Logger.Info("Creating IMS Configmap")
			err = r.client.Create(req.Ctx, imsConfig)
			if err == nil {
				return res.SetCreated()
			}
		}
		return res.Error(err)
	}

	req.Logger.Info("IMS Configmap already exists", "imsConfigMap.Namespace", found.Namespace, "imsConfigMap.Name", found.Name)

	// Add it to the list of RelatedObjects if found
	objectRef, err := reference.GetReference(r.scheme, found)
	if err != nil {
		return res.Error(err)
	}
	objectreferencesv1.SetObjectReference(&req.Instance.Status.RelatedObjects, *objectRef)

	// in an ideal world HCO should be managing the whole config map,
	// now due to a bad design only a few values of this config map are
	// really managed by HCO while others are managed by other entities
	// TODO: fix this bad design splitting the config map into two distinct objects and reconcile the whole object here
	needsUpdate := false
	for key, value := range imsConfig.Data {
		if found.Data[key] != value {
			found.Data[key] = value
			needsUpdate = true
		}
	}
	if needsUpdate {
		req.Logger.Info("Updating existing IMS Configmap to its default values")
		err = r.client.Update(req.Ctx, found)
		if err != nil {
			return res.Error(err)
		}
		return res.SetUpdated()
	}

	return res.SetUpgradeDone(req.ComponentUpgradeInProgress)
}

func (r *ReconcileHyperConverged) ensureVMImport(req *common.HcoRequest) *operands.EnsureResult {
	return operandMap["vmimport"].Ensure(req)
}

func (r *ReconcileHyperConverged) ensureConsoleCLIDownload(req *common.HcoRequest) error {
	ccd := req.Instance.NewConsoleCLIDownload()

	found := req.Instance.NewConsoleCLIDownload()
	err := hcoutil.EnsureCreated(req.Ctx, r.client, found, req.Logger)
	if err != nil {
		if meta.IsNoMatchError(err) {
			req.Logger.Info("ConsoleCLIDownload was not found, skipping")
		}
		return err
	}

	// Make sure we hold the right link spec
	if reflect.DeepEqual(found.Spec, ccd.Spec) {
		objectRef, err := reference.GetReference(r.scheme, found)
		if err != nil {
			req.Logger.Error(err, "failed getting object reference for ConsoleCLIDownload")
			return err
		}
		objectreferencesv1.SetObjectReference(&req.Instance.Status.RelatedObjects, *objectRef)
		return nil
	}

	ccd.Spec.DeepCopyInto(&found.Spec)

	err = r.client.Update(req.Ctx, found)
	if err != nil {
		return err
	}

	return nil
}

func newKubeVirtTemplateValidatorForCR(cr *hcov1beta1.HyperConverged, namespace string) *sspv1.KubevirtTemplateValidator {
	labels := map[string]string{
		hcoutil.AppLabel: cr.Name,
	}

	spec := sspv1.TemplateValidatorSpec{}
	if cr.Spec.Infra.NodePlacement != nil {
		if cr.Spec.Infra.NodePlacement.Affinity != nil {
			cr.Spec.Infra.NodePlacement.Affinity.DeepCopyInto(&spec.Affinity)
		}

		if cr.Spec.Infra.NodePlacement.NodeSelector != nil {
			spec.NodeSelector = make(map[string]string)
			for k, v := range cr.Spec.Infra.NodePlacement.NodeSelector {
				spec.NodeSelector[k] = v
			}
		}

		for _, hcoTolr := range cr.Spec.Infra.NodePlacement.Tolerations {
			tvTolr := corev1.Toleration{}
			hcoTolr.DeepCopyInto(&tvTolr)
			spec.Tolerations = append(spec.Tolerations, tvTolr)
		}
	}

	return &sspv1.KubevirtTemplateValidator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "template-validator-" + cr.Name,
			Labels:    labels,
			Namespace: namespace,
		},
		Spec: spec,
	}
}

func (r *ReconcileHyperConverged) ensureKubeVirtTemplateValidator(req *common.HcoRequest) *operands.EnsureResult {
	kvTV := newKubeVirtTemplateValidatorForCR(req.Instance, req.Namespace)
	res := operands.NewEnsureResult(kvTV)
	if !r.clusterInfo.IsOpenshift() { // SSP operators Only supported in OpenShift. Ignore in K8s.
		return res.SetUpgradeDone(true)
	}

	if err := controllerutil.SetControllerReference(req.Instance, kvTV, r.scheme); err != nil {
		return res.Error(err)
	}

	key, err := client.ObjectKeyFromObject(kvTV)
	if err != nil {
		req.Logger.Error(err, "Failed to get object key for KubeVirt Template Validator")
	}
	res.SetName(key.Name)

	found := &sspv1.KubevirtTemplateValidator{}
	err = r.client.Get(req.Ctx, key, found)
	if err != nil {
		if apierrors.IsNotFound(err) {
			req.Logger.Info("Creating KubeVirt Template Validator")
			err = r.client.Create(req.Ctx, kvTV)
			if err == nil {
				return res.SetCreated()
			}
		}
		return res.Error(err)
	}

	req.Logger.Info("KubeVirt Template Validator already exists", "validator.Namespace", found.Namespace, "validator.Name", found.Name)

	if !reflect.DeepEqual(kvTV.Spec, found.Spec) {
		req.Logger.Info("Updating existing KubeVirt Template Validator")
		kvTV.Spec.DeepCopyInto(&found.Spec)
		err = r.client.Update(req.Ctx, found)
		if err != nil {
			return res.Error(err)
		}
		return res.SetUpdated()
	}
	// Add it to the list of RelatedObjects if found
	objectRef, err := reference.GetReference(r.scheme, found)
	if err != nil {
		return res.Error(err)
	}
	objectreferencesv1.SetObjectReference(&req.Instance.Status.RelatedObjects, *objectRef)

	isReady := handleComponentConditions(r, req, "KubevirtTemplateValidator", found.Status.Conditions)

	upgradeInProgress := false
	if isReady {
		upgradeInProgress = r.upgradeMode && r.checkComponentVersion(hcoutil.SspVersionEnvV, found.Status.ObservedVersion)
		if (upgradeInProgress || !r.upgradeMode) && r.shouldRemoveOldCrd[templateValidatorsOldCrdName] {
			if r.removeCrd(req, templateValidatorsOldCrdName) {
				r.shouldRemoveOldCrd[templateValidatorsOldCrdName] = false
			}
		}
	}

	return res.SetUpgradeDone(req.ComponentUpgradeInProgress && upgradeInProgress)
}

func newKubeVirtMetricsAggregationForCR(cr *hcov1beta1.HyperConverged, namespace string) *sspv1.KubevirtMetricsAggregation {
	labels := map[string]string{
		hcoutil.AppLabel: cr.Name,
	}
	return &sspv1.KubevirtMetricsAggregation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "metrics-aggregation-" + cr.Name,
			Labels:    labels,
			Namespace: namespace,
		},
	}
}

func (r *ReconcileHyperConverged) ensureKubeVirtMetricsAggregation(req *common.HcoRequest) *operands.EnsureResult {
	kubevirtMetricsAggregation := newKubeVirtMetricsAggregationForCR(req.Instance, req.Namespace)
	res := operands.NewEnsureResult(kubevirtMetricsAggregation)
	if !r.clusterInfo.IsOpenshift() { // SSP operators Only supported in OpenShift. Ignore in K8s.
		return res.SetUpgradeDone(true)
	}

	err := controllerutil.SetControllerReference(req.Instance, kubevirtMetricsAggregation, r.scheme)
	if err != nil {
		return res.Error(err)
	}

	key, err := client.ObjectKeyFromObject(kubevirtMetricsAggregation)
	if err != nil {
		req.Logger.Error(err, "Failed to get object key for KubeVirt Metrics Aggregation")
	}

	res.SetName(key.Name)
	found := &sspv1.KubevirtMetricsAggregation{}

	err = r.client.Get(req.Ctx, key, found)
	if err != nil {
		if apierrors.IsNotFound(err) {
			req.Logger.Info("Creating KubeVirt Metrics Aggregation")
			err = r.client.Create(req.Ctx, kubevirtMetricsAggregation)
			if err == nil {
				return res.SetCreated()
			}
		}
		return res.Error(err)
	}

	req.Logger.Info("KubeVirt Metrics Aggregation already exists", "metrics.Namespace", found.Namespace, "metrics.Name", found.Name)

	if !reflect.DeepEqual(kubevirtMetricsAggregation.Spec, found.Spec) {
		req.Logger.Info("Updating existing KubeVirt Metrics Aggregation")
		kubevirtMetricsAggregation.Spec.DeepCopyInto(&found.Spec)
		err = r.client.Update(req.Ctx, found)
		if err != nil {
			return res.Error(err)
		}
		return res.SetUpdated()
	}
	// Add it to the list of RelatedObjects if found
	objectRef, err := reference.GetReference(r.scheme, found)
	if err != nil {
		return res.Error(err)
	}
	objectreferencesv1.SetObjectReference(&req.Instance.Status.RelatedObjects, *objectRef)

	isReady := handleComponentConditions(r, req, "KubeVirtMetricsAggregation", found.Status.Conditions)

	upgradeInProgress := false
	if isReady {
		upgradeInProgress = r.upgradeMode && r.checkComponentVersion(hcoutil.SspVersionEnvV, found.Status.ObservedVersion)
		if (upgradeInProgress || !r.upgradeMode) && r.shouldRemoveOldCrd[metricsAggregationOldCrdName] {
			if r.removeCrd(req, metricsAggregationOldCrdName) {
				r.shouldRemoveOldCrd[metricsAggregationOldCrdName] = false
			}
		}
	}

	return res.SetUpgradeDone(req.ComponentUpgradeInProgress && upgradeInProgress)
}

// This function is used to exit from the reconcile function, updating the conditions and returns the reconcile result
func (r *ReconcileHyperConverged) updateConditions(req *common.HcoRequest) error {
	for _, condType := range common.HcoConditionTypes {
		cond, found := req.Conditions[condType]
		if !found {
			cond = conditionsv1.Condition{
				Type:    condType,
				Status:  corev1.ConditionUnknown,
				Message: "Unknown Status",
			}
		}
		conditionsv1.SetStatusCondition(&req.Instance.Status.Conditions, cond)
	}

	req.StatusDirty = true
	return nil
}

func (r *ReconcileHyperConverged) setLabels(req *common.HcoRequest) {
	if req.Instance.ObjectMeta.Labels == nil {
		req.Instance.ObjectMeta.Labels = map[string]string{}
	}
	if req.Instance.ObjectMeta.Labels[hcoutil.AppLabel] == "" {
		req.Instance.ObjectMeta.Labels[hcoutil.AppLabel] = req.Instance.Name
		req.Dirty = true
	}
}

// return true if not found or if deletion succeeded
func (r *ReconcileHyperConverged) removeCrd(req *common.HcoRequest, crdName string) bool {
	found := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "CustomResourceDefinition",
			"apiVersion": "apiextensions.k8s.io/v1",
		},
	}
	key := client.ObjectKey{Namespace: req.Namespace, Name: crdName}
	err := r.client.Get(req.Ctx, key, found)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			req.Logger.Error(err, fmt.Sprintf("failed to read the %s CRD; %s", crdName, err.Error()))
			return false
		}
	} else {
		err = r.client.Delete(req.Ctx, found)
		if err != nil {
			req.Logger.Error(err, fmt.Sprintf("failed to remove the %s CRD; %s", crdName, err.Error()))
			return false
		} else {
			req.Logger.Info("successfully removed CRD", "CRD Name", crdName)
		}
	}

	return true
}

// getHyperconverged returns the name/namespace of the HyperConverged resource
func getHyperconverged() (types.NamespacedName, error) {
	hco := types.NamespacedName{
		Name: hcoutil.HyperConvergedName,
	}

	namespace, err := hcoutil.GetOperatorNamespaceFromEnv()
	if err != nil {
		return hco, err
	}
	hco.Namespace = namespace

	return hco, nil
}

func contains(slice []string, s string) bool {
	for _, element := range slice {
		if element == s {
			return true
		}
	}
	return false
}

func drop(slice []string, s string) []string {
	newSlice := []string{}
	for _, element := range slice {
		if element != s {
			newSlice = append(newSlice, element)
		}
	}
	return newSlice
}

// translateKubeVirtConds translates list of KubeVirt conditions to a list of custom resource
// conditions.
func translateKubeVirtConds(orig []kubevirtv1.KubeVirtCondition) []conditionsv1.Condition {
	translated := make([]conditionsv1.Condition, len(orig))

	for i, origCond := range orig {
		translated[i] = conditionsv1.Condition{
			Type:    conditionsv1.ConditionType(origCond.Type),
			Status:  origCond.Status,
			Reason:  origCond.Reason,
			Message: origCond.Message,
		}
	}

	return translated
}
