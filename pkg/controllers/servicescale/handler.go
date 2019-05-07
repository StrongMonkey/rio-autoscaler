package servicescale

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/knative/serving/pkg/autoscaler"
	autoscalev1 "github.com/rancher/rio/pkg/apis/autoscale.rio.cattle.io/v1"
	corev1controller "github.com/rancher/rio/pkg/generated/controllers/core/v1"
	riov1controller "github.com/rancher/rio/pkg/generated/controllers/rio.cattle.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var SyncMap sync.Map

type ssrHandler struct {
	ctx         context.Context
	metrics     *autoscaler.MultiScaler
	pollers     map[string]*poller
	pollerLock  sync.Mutex
	rioServices riov1controller.ServiceController
	services    corev1controller.ServiceCache
	pods        corev1controller.PodCache
}

func NewHandler(ctx context.Context, metrics *autoscaler.MultiScaler,
	rioServiceClient riov1controller.ServiceController,
	serviceClientCache corev1controller.ServiceCache,
	podClientCache corev1controller.PodCache) *ssrHandler {

	return &ssrHandler{
		ctx:         ctx,
		metrics:     metrics,
		pollers:     map[string]*poller{},
		rioServices: rioServiceClient,
		services:    serviceClientCache,
		pods:        podClientCache,
	}
}

func (s *ssrHandler) OnChange(key string, obj runtime.Object) (runtime.Object, error) {
	ssr := obj.(*autoscalev1.ServiceScaleRecommendation)
	m, err := s.createMetric(ssr)
	if err != nil {
		return ssr, err
	}

	s.monitor(ssr)

	ssr.Status.DesiredScale = bounded(m.Status.DesiredScale, ssr.Spec.MinScale, ssr.Spec.MaxScale)
	return ssr, SetDeploymentScale(s.rioServices, ssr)
}

func bounded(value, lower, upper int32) *int32 {
	if value < lower {
		return &lower
	}
	if upper > 0 && value > upper {
		return &upper
	}
	return &value
}

func (s *ssrHandler) OnRemove(key string, ssr *autoscalev1.ServiceScaleRecommendation) (*autoscalev1.ServiceScaleRecommendation, error) {
	s.pollerLock.Lock()
	defer s.pollerLock.Unlock()

	p := s.pollers[key]
	if p != nil {
		p.Stop()
	}

	delete(s.pollers, key)
	return ssr, s.deleteMetric(ssr)
}

func (s *ssrHandler) deleteMetric(ssr *autoscalev1.ServiceScaleRecommendation) error {
	return s.metrics.Delete(s.ctx, ssr.Namespace, ssr.Name)
}

func (s *ssrHandler) createMetric(ssr *autoscalev1.ServiceScaleRecommendation) (*autoscaler.Decider, error) {
	metric, err := s.metrics.Get(s.ctx, ssr.Namespace, ssr.Name)
	if err != nil && errors.IsNotFound(err) {
		return s.metrics.Create(s.ctx, &autoscaler.Decider{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ssr.Name,
				Namespace: ssr.Namespace,
			},
			Spec: autoscaler.DeciderSpec{
				TargetConcurrency: float64(ssr.Spec.Concurrency),
			},
		})
	} else if err != nil {
		return nil, err
	}
	return metric, nil
}

func SetDeploymentScale(rioServices riov1controller.ServiceController, ssr *autoscalev1.ServiceScaleRecommendation) error {
	svc, err := rioServices.Cache().Get(ssr.Namespace, ssr.Name)
	if err != nil {
		return err
	}
	// wait for a minute after scale from zero
	if svc.Status.ScaleFromZeroTimestamp != nil && svc.Status.ScaleFromZeroTimestamp.Add(time.Minute).After(time.Now()) {
		logrus.Infof("skipping setting scale because service  %s/%s is scaled from zero within a minute", svc.Namespace, svc.Name)
		return nil
	}
	logrus.Infof("Setting desired scale %v for %v/%v", *ssr.Status.DesiredScale, svc.Namespace, svc.Name)
	observedScale := int(*ssr.Status.DesiredScale)
	svc.Status.ObservedScale = &observedScale
	if _, err := rioServices.Update(svc); err != nil {
		return err
	}
	return nil
}

func (s *ssrHandler) monitor(ssr *autoscalev1.ServiceScaleRecommendation) {
	s.pollerLock.Lock()
	defer s.pollerLock.Unlock()

	key := key(ssr)
	p, ok := s.pollers[key]
	if ok {
		p.update(ssr)
		return
	}

	p = newPoller(s.ctx, ssr, s.pods, func(stat autoscaler.Stat) {
		s.metrics.RecordStat(key, stat)
	})

	s.pollers[key] = p
}

func key(ssr *autoscalev1.ServiceScaleRecommendation) string {
	return fmt.Sprintf("%s/%s", ssr.Namespace, ssr.Name)
}
