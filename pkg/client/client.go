package client

import (
	"github.com/xiang90/kprober/pkg/spec"
	"github.com/xiang90/kprober/pkg/util/k8sutil"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

func MustNewInCluster() (*rest.RESTClient, *runtime.Scheme) {
	cfg, err := k8sutil.InClusterConfig()
	if err != nil {
		panic(err)
	}
	cli, scheme, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return cli, scheme
}

func New(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := spec.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	config := *cfg
	config.GroupVersion = &spec.SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, nil, err
	}

	return client, scheme, nil
}
