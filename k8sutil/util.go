package k8sutil

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
