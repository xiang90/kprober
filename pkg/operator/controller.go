package operator

import (
	"context"
	"fmt"

	"github.com/xiang90/kprober/pkg/spec"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

func (p *Probers) run(ctx context.Context) {
	source := cache.NewListWatchFromClient(
		p.probersClient,
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
}

func (p *Probers) onAdd(obj interface{}) {
	prober := obj.(*spec.Prober)
	fmt.Printf("Add: %v\n", prober)
}

func (p *Probers) onUpdate(oldObj, newObj interface{}) {
	oldProber := oldObj.(*spec.Prober)
	newProber := newObj.(*spec.Prober)
	fmt.Printf("Update: old: %v, new: %v", oldProber, newProber)
}

func (p *Probers) onDelete(obj interface{}) {
	prober := obj.(*spec.Prober)
	fmt.Printf("Delete: %v\n", prober)
}
