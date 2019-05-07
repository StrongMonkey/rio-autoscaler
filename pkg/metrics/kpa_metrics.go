package metrics

import (
	"context"
	"time"

	"github.com/knative/serving/pkg/autoscaler"
	"github.com/rancher/rio-autoscaler/pkg/logger"
	"github.com/rancher/rio-autoscaler/types"
	kubeinformers "k8s.io/client-go/informers"
)

var (
	DefaultConfig = autoscaler.Config{
		ContainerConcurrencyTargetDefault:    1,
		ContainerConcurrencyTargetPercentage: 1.0,
		EnableScaleToZero:                    true,
		MaxScaleUpRate:                       10,
		PanicWindow:                          6 * time.Second,
		ScaleToZeroGracePeriod:               2 * time.Minute,
		StableWindow:                         60 * time.Second,
		TickInterval:                         2 * time.Second,
	}
)

func New(ctx context.Context, rContext *types.Context) *autoscaler.MultiScaler {
	dynConfig := autoscaler.NewDynamicConfig(&DefaultConfig, logger.SugaredLogger)
	return autoscaler.NewMultiScaler(dynConfig, ctx.Done(), nil, wrap(rContext), nil, logger.SugaredLogger)
}

func wrap(rContext *types.Context) func(decider *autoscaler.Decider, config *autoscaler.DynamicConfig) (autoscaler.UniScaler, error) {
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(rContext.K8s, time.Second*30)
	return func(decider *autoscaler.Decider, config *autoscaler.DynamicConfig) (scaler autoscaler.UniScaler, e error) {
		return autoscaler.New(config, decider.Namespace, decider.Name, kubeInformerFactory.Core().V1().Endpoints(), decider.Spec.TargetConcurrency, (*nilReporter)(nil))
	}
}

type nilReporter struct {
}

func (n *nilReporter) ReportDesiredPodCount(v int64) error {
	return nil
}
func (n *nilReporter) ReportRequestedPodCount(v int64) error {
	return nil
}
func (n *nilReporter) ReportActualPodCount(v int64) error {
	return nil
}
func (n *nilReporter) ReportObservedPodCount(v float64) error {
	return nil
}
func (n *nilReporter) ReportStableRequestConcurrency(v float64) error {
	return nil
}
func (n *nilReporter) ReportPanicRequestConcurrency(v float64) error {
	return nil
}
func (n *nilReporter) ReportTargetRequestConcurrency(v float64) error {
	return nil
}
func (n *nilReporter) ReportPanic(v int64) error {
	return nil
}
