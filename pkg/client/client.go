package client

import (
	"context"

	"github.com/xiang90/kprober/pkg/spec"
	"github.com/xiang90/kprober/pkg/util/k8sutil"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

type ProbersCR interface {
	Get(ctx context.Context, namespace, name string) (*spec.Prober, error)
	RESTClient() *rest.RESTClient
}

type probersClient struct {
	restCli    *rest.RESTClient
	crScheme   *runtime.Scheme
	paramCodec runtime.ParameterCodec
}

func (pc *probersClient) RESTClient() *rest.RESTClient {
	return pc.restCli
}

func (pc *probersClient) Get(ctx context.Context, ns, name string) (*spec.Prober, error) {
	res := &spec.Prober{}
	err := pc.restCli.Get().
		Resource(spec.ProberResourcePlural).
		Namespace(ns).
		Name(name).
		Do().
		Into(res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func MustNewInCluster() ProbersCR {
	cfg, err := k8sutil.InClusterConfig()
	if err != nil {
		panic(err)
	}
	cli, crScheme, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return &probersClient{
		restCli:    cli,
		crScheme:   crScheme,
		paramCodec: runtime.NewParameterCodec(crScheme),
	}
}

func New(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	crScheme := runtime.NewScheme()
	if err := spec.AddToScheme(crScheme); err != nil {
		return nil, nil, err
	}

	config := *cfg
	config.GroupVersion = &spec.SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(crScheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, nil, err
	}

	return client, crScheme, nil
}
