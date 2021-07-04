package util

import (
	"context"
	"os"

	"github.com/go-logr/logr"
	openshiftconfigv1 "github.com/openshift/api/config/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterInfo interface {
	Init(ctx context.Context, creader client.Reader, logger logr.Logger, runningLocally bool) error
	IsOpenshift() bool
	IsRunningLocally() bool
	IsManagedByOLM() bool
}

var _ ClusterInfo = (*ClusterInfoImp)(nil)

type ClusterInfoImp struct {
	runningInOpenshift bool
	managedByOLM       bool
	runningLocally     bool
}

var clusterInfo ClusterInfo

var GetClusterInfo = func() ClusterInfo {
	return clusterInfo
}

// OperatorConditionNameEnvVar - this Env var is set by OLM, so the Operator can discover it's OperatorCondition.
const OperatorConditionNameEnvVar = "OPERATOR_CONDITION_NAME"

func (c *ClusterInfoImp) Init(ctx context.Context, creader client.Reader, logger logr.Logger, runningLocally bool) error {
	c.checkManagedByOLM()
	return c.checkRunningInOpenshift(ctx, creader, logger, runningLocally)
}

func (c *ClusterInfoImp) checkManagedByOLM() {
	// We assume that this Operator is managed by OLM when this variable is present.
	_, c.managedByOLM = os.LookupEnv(OperatorConditionNameEnvVar)
}

func (c *ClusterInfoImp) checkRunningInOpenshift(ctx context.Context, creader client.Reader, logger logr.Logger, runningLocally bool) error {
	c.runningLocally = runningLocally
	isOpenShift := false
	version := ""

	clusterVersion := &openshiftconfigv1.ClusterVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "version",
		},
	}
	if err := creader.Get(ctx, client.ObjectKeyFromObject(clusterVersion), clusterVersion); err != nil {
		if meta.IsNoMatchError(err) || apierrors.IsNotFound(err) {
			// Not on OpenShift
			isOpenShift = false
		} else {
			logger.Error(err, "Failed to get ClusterVersion")
			return err
		}
	} else {
		isOpenShift = true
		version = clusterVersion.Status.Desired.Version
	}

	c.runningInOpenshift = isOpenShift
	if isOpenShift {
		logger.Info("Cluster type = openshift", "version", version)
	} else {
		logger.Info("Cluster type = kubernetes")
	}

	return nil
}

func (c ClusterInfoImp) IsOpenshift() bool {
	return c.runningInOpenshift
}

func (c ClusterInfoImp) IsRunningLocally() bool {
	return c.runningLocally
}

func (c ClusterInfoImp) IsManagedByOLM() bool {
	return c.managedByOLM
}

func init() {
	clusterInfo = &ClusterInfoImp{
		runningInOpenshift: false,
	}
}
