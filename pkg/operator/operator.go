package operator

import "context"
import "k8s.io/client-go/kubernetes"

type Probers struct {
	kubecli kubernetes.Interface
}

func New(kubecli kubernetes.Interface) *Probers {
	return &Probers{
		kubecli: kubecli,
	}
}

func (p *Probers) Start(ctx context.Context) {
	p.init(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (p *Probers) init(ctx context.Context) {
	// TODO: Create CRD
}
