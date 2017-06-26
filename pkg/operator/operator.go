package operator

import (
	"context"

	"github.com/xiang90/kprober/pkg/util/k8sutil"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

const (
	crdGroupName      = "monitoring.coreos.com"
	crdResourcePlural = "probers"
	crdResourceKind   = "Prober"
)

var crdGroupVersion = schema.GroupVersion{Group: crdGroupName, Version: "v1alpha1"}

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
			Name: crdResourcePlural + "." + crdGroupName,
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   crdGroupName,
			Version: crdGroupVersion.Version,
			Scope:   apiextensionsv1beta1.NamespaceScoped,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural: crdResourcePlural,
				Kind:   crdResourceKind,
			},
		},
	}
	_, err := p.kubeExtCli.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	return err
}
