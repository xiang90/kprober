package main

import (
	"context"

	"github.com/xiang90/kprober/pkg/operator"
	"github.com/xiang90/kprober/pkg/util/k8sutil"
)

func main() {
	po := operator.New(k8sutil.MustNewKubeClient())
	po.Start(context.TODO())
}
