package hyperconverged

import (
	hcov1alpha1 "github.com/kubevirt/hyperconverged-cluster-operator/pkg/apis/hco/v1alpha1"
	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var (
	hcoConditionTypes = []conditionsv1.ConditionType{
		hcov1alpha1.ConditionReconcileComplete,
		conditionsv1.ConditionAvailable,
		conditionsv1.ConditionProgressing,
		conditionsv1.ConditionDegraded,
		conditionsv1.ConditionUpgradeable,
	}
)

type hcoConditions map[conditionsv1.ConditionType]conditionsv1.Condition

func newHcoConditions() hcoConditions {
	return hcoConditions{}
}

func (hc hcoConditions) setStatusCondition(newCondition conditionsv1.Condition) {
	existingCondition, ok := hc[newCondition.Type]

	if !ok {
		newCondition.LastTransitionTime = metav1.NewTime(time.Now())
		newCondition.LastHeartbeatTime = metav1.NewTime(time.Now())
		hc[newCondition.Type] = newCondition
		return
	}

	if existingCondition.Status != newCondition.Status {
		existingCondition.Status = newCondition.Status
		existingCondition.LastTransitionTime = metav1.NewTime(time.Now())
	}

	existingCondition.Reason = newCondition.Reason
	existingCondition.Message = newCondition.Message
	existingCondition.LastHeartbeatTime = metav1.NewTime(time.Now())
}

func (hc hcoConditions) empty() bool {
	return len(hc) == 0
}
