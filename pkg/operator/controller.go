package operator

import (
	"context"
	"log"

	"github.com/xiang90/kprober/pkg/spec"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/pkg/api/v1"
	appsv1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"k8s.io/client-go/tools/cache"
)

const proberImage = "gcr.io/coreos-k8s-scale-testing/kprober"

func (p *Probers) run(ctx context.Context) {
	source := cache.NewListWatchFromClient(
		p.probersCli.RESTClient(),
		spec.ProberResourcePlural,
		p.namespace,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,
		// The object type.
		&spec.Prober{},
		// resyncPeriod
		// Every resyncPeriod, all resources in the cache will retrigger events.
		// Set to 0 to disable the resync.
		0,
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			AddFunc:    p.onAdd,
			UpdateFunc: p.onUpdate,
			DeleteFunc: p.onDelete,
		})

	go controller.Run(ctx.Done())
	log.Println("start processing probers changes...")
}

func (p *Probers) onAdd(obj interface{}) {
	pr := obj.(*spec.Prober)

	selector := map[string]string{"app": "prober", "prober": pr.Name}

	podTempl := v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:   pr.Name,
			Labels: selector,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{
				Name:  "prober",
				Image: proberImage,
				Command: []string{
					"prober",
					"-n=" + pr.Name,
					"-ns=" + pr.Namespace,
				},
			}},
		},
	}

	d := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   pr.Name,
			Labels: map[string]string{"prober": pr.Name},
		},
		Spec: appsv1beta1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: selector},
			Template: podTempl,
			Strategy: appsv1beta1.DeploymentStrategy{
				Type: appsv1beta1.RecreateDeploymentStrategyType,
			},
		},
	}
	_, err := p.kubecli.AppsV1beta1().Deployments(pr.Namespace).Create(d)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		// TODO: retry or report failure status in CR
		panic(err)
	}

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   pr.Name,
			Labels: selector,
		},
		Spec: v1.ServiceSpec{
			Selector: selector,
			Ports: []v1.ServicePort{{
				Name: "metrics",
				Port: 17783,
			}},
		},
	}

	_, err = p.kubecli.CoreV1().Services(pr.Namespace).Create(svc)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		// TODO: retry or report failure status in CR
		panic(err)
	}
}

func (p *Probers) onUpdate(oldObj, newObj interface{}) {
	oldProber := oldObj.(*spec.Prober)
	newProber := newObj.(*spec.Prober)
	log.Printf("Update: old: %v, new: %v", oldProber, newProber)
}

func (p *Probers) onDelete(obj interface{}) {
	prober := obj.(*spec.Prober)
	log.Printf("Delete: %v", prober)
}
