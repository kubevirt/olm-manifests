package util

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("OperatorCondition", func() {
	DescribeTable("should return no error when setting the condition, in not-supported environments", func(ci ClusterInfo) {
		oc, err := NewOperatorCondition(ci, nil, UpgradeableCondition)
		Expect(err).To(BeNil())

		ctx := context.Background()
		err = oc.Set(ctx, metav1.ConditionTrue, "Reason", "message")
		Expect(err).To(BeNil())
	},
		Entry("should no-op when not managed by OLM", &ClusterInfoImp{
			managedByOLM:   false,
			runningLocally: false,
		}),
		Entry("should no-op when running locally", &ClusterInfoImp{
			managedByOLM:   true,
			runningLocally: true,
		}),
		Entry("should no-op when running locally and not managed by OLM", &ClusterInfoImp{
			managedByOLM:   false,
			runningLocally: true,
		}),
	)
})

func TestOperatorCondition(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OperatorCondition Suite")
}
