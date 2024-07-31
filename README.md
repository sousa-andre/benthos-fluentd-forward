# benthos-fluentd-forward
## Installation
```shell
go mod get github.com/sousa-andre/benthos-fluentd-forward`
```
## Components
### Output

#### Configuration
```yaml
output:
  fluentd_forward:
    tag: "benthos.ouput"
```


## Build
```go
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
```