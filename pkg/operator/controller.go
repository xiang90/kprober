package operator

import (
	"context"
	"log"

	"github.com/xiang90/kprober/pkg/spec"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

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
	prober := obj.(*spec.Prober)
	log.Printf("Add: %v", prober)
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
