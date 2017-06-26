package k8sutil

import (
	"net"
	"os"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func IPFromPod(ns, podname string) (string, error) {
	return "", nil
}

func IPFromService() string {
	return ""
}

func IPsFromReplicaSet() []string {
	return nil

}

func IPsFromDeployments() []string {
	return nil
}

func DeployProber() error {
	return nil
}

func UpdateProber() error {
	return nil
}

func DeleteProber() error {
	return nil
}

func MustNewKubeClient() kubernetes.Interface {
	cfg, err := InClusterConfig()
	if err != nil {
		panic(err)
	}
	return kubernetes.NewForConfigOrDie(cfg)
}

func MustNewKubeExtClient() kubernetes.Interface {
	cfg, err := InClusterConfig()
	return apiextensionsclient.NewForConfigOrDie(cfg)
}

func InClusterConfig() (*rest.Config, error) {
	// Work around https://github.com/kubernetes/kubernetes/issues/40973
	// See https://github.com/coreos/etcd-operator/issues/731#issuecomment-283804819
	if len(os.Getenv("KUBERNETES_SERVICE_HOST")) == 0 {
		addrs, err := net.LookupHost("kubernetes.default.svc")
		if err != nil {
			panic(err)
		}
		os.Setenv("KUBERNETES_SERVICE_HOST", addrs[0])
	}
	if len(os.Getenv("KUBERNETES_SERVICE_PORT")) == 0 {
		os.Setenv("KUBERNETES_SERVICE_PORT", "443")
	}
	return rest.InClusterConfig()
}