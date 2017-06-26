package main

import (
	"context"

	"github.com/xiang90/kprober/pkg/operator"
)

func main() {
	po := operator.New()
	po.Start(context.TODO())
}
