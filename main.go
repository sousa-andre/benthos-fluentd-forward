package main

import (
	"context"
	"github.com/redpanda-data/benthos/v4/public/service"
	_ "github.com/redpanda-data/connect/v4/public/components/all"
	_ "github.com/sousa-andre/benthos-fluentd-forward/output"
)

func main() {
	service.RunCLI(context.Background())
}
