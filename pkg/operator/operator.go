package operator

import (
	"context"
	"os"

	"github.com/xiang90/kprober/pkg/client"
	"github.com/xiang90/kprober/pkg/util/k8sutil"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Probers struct {
	kubecli       kubernetes.Interface
	kubeExtCli    apiextensionsclient.Interface
	probersClient *rest.RESTClient
	namespace     string
}

func New() *Probers {
	probersCli, _ := client.MustNewInCluster()
	return &Probers{
		kubecli:       k8sutil.MustNewKubeClient(),
		kubeExtCli:    k8sutil.MustNewKubeExtClient(),
		probersClient: probersCli,
		namespace:     os.Getenv("MY_POD_NAMESPACE"),
	}
}

func (p *Probers) Start(ctx context.Context) {
	p.init(ctx)
	p.run(ctx)
	<-ctx.Done()
}

func (p *Probers) init(ctx context.Context) error {
	err := k8sutil.CreateProberCRD(p.kubeExtCli)
	if err != nil {
		return err
	}
	return k8sutil.WaitProberCRDCreated(p.kubeExtCli)
}
