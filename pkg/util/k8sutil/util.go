package k8sutil

import (
	"net"
	"os"

	"github.com/xiang90/kprober/pkg/spec"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	appsv1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"k8s.io/client-go/rest"
)

const (
	ContainerProbeDirPath        = "/var/tmp/containerprobe"
	ContainerProbeOutputFilePath = ContainerProbeDirPath + "/result"

	proberImage = "gcr.io/coreos-k8s-scale-testing/kprober"
)

func IPFromPod(ns, podname string) (string, error) {
	return "", nil
}

func IPFromService(kubecli kubernetes.Interface, ns, svcName string) (string, error) {
	svc, err := kubecli.CoreV1().Services(ns).Get(svcName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return svc.Spec.ClusterIP, nil
}

func IPsFromReplicaSet() []string {
	return nil

}

func IPsFromDeployments() []string {
	return nil
}

func DeployProber(kubecli kubernetes.Interface, pr *spec.Prober) error {
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

	if c := pr.Spec.Probe.Container; c != nil {
		vn := "sharedResult"
		shared := v1.Volume{
			Name: vn,
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		}
		vm := v1.VolumeMount{
			Name:      vn,
			MountPath: ContainerProbeDirPath,
			ReadOnly:  false,
		}
		podTempl.Spec.Volumes = append(podTempl.Spec.Volumes, shared)

		cp := v1.Container{
			Name:         "container-probe",
			Image:        c.Image,
			Command:      []string{"probe > " + ContainerProbeOutputFilePath},
			VolumeMounts: []v1.VolumeMount{vm},
			// TODO: add IP and Target Env.
		}
		podTempl.Spec.Containers[0].VolumeMounts = append(podTempl.Spec.Containers[0].VolumeMounts, vm)
		podTempl.Spec.Containers = append(podTempl.Spec.Containers, cp)
	}

	d := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   pr.Name,
			Labels: selector,
		},
		Spec: appsv1beta1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: selector},
			Template: podTempl,
			Strategy: appsv1beta1.DeploymentStrategy{
				Type: appsv1beta1.RecreateDeploymentStrategyType,
			},
		},
	}
	_, err := kubecli.AppsV1beta1().Deployments(pr.Namespace).Create(d)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
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

	_, err = kubecli.CoreV1().Services(pr.Namespace).Create(svc)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func UpdateProber() error {
	return nil
}

func DeleteProber(kubecli kubernetes.Interface, ns, name string) error {
	err := kubecli.AppsV1beta1().Deployments(ns).Delete(name, cascadeDeleteOptions(0))
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	err = kubecli.CoreV1().Services(ns).Delete(name, nil)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}

func cascadeDeleteOptions(gracePeriodSeconds int64) *metav1.DeleteOptions {
	return &metav1.DeleteOptions{
		GracePeriodSeconds: func(t int64) *int64 { return &t }(gracePeriodSeconds),
		PropagationPolicy: func() *metav1.DeletionPropagation {
			foreground := metav1.DeletePropagationForeground
			return &foreground
		}(),
	}
}

func MustNewKubeClient() kubernetes.Interface {
	cfg, err := InClusterConfig()
	if err != nil {
		panic(err)
	}
	return kubernetes.NewForConfigOrDie(cfg)
}

func MustNewKubeExtClient() apiextensionsclient.Interface {
	cfg, err := InClusterConfig()
	if err != nil {
		panic(err)
	}
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
