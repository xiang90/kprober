package operator

import (
	"context"

	"github.com/xiang90/kprober/pkg/spec"
	"github.com/xiang90/kprober/pkg/util/k8sutil"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Probers struct {
	kubecli    kubernetes.Interface
	kubeExtCli apiextensionsclient.Interface
}

func New() *Probers {
	return &Probers{
		kubecli:    k8sutil.MustNewKubeClient(),
		kubeExtCli: k8sutil.MustNewKubeExtClient(),
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

func (p *Probers) init(ctx context.Context) error {
	crd := &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: spec.ProberResourcePlural + "." + spec.GroupName,
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   spec.GroupName,
			Version: spec.SchemeGroupVersion.Version,
			Scope:   apiextensionsv1beta1.NamespaceScoped,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural: spec.ProberResourcePlural,
				Kind:   spec.ProberResourceKind,
			},
		},
	}
	_, err := p.kubeExtCli.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	return err
}
