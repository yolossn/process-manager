<h1 align="center">Process Manager</h1>

### Design Decision

1. Manager provides the process with command, backoff strategy and maxRetries. The process takes care of retries internally.
   Why ?
   In this case the manager doesnt dynamically change the retry strategy or retry count based on the process outcome so the sole responsibility of completing the command with retries belongs to the process.
2. Both processes which have reached maxRetries and the ones which have run succesfully are considered complete.
3. If an Interupt or Kill is made when the processes are executing the manager kills the processes which have not completed yet and exits.

### Improvements

1. Improve tests (Add table tests, increase coverage)
2. Make backoff strategy configurable. Now it is a constant backoff of 1 second
3. The use of GOTO statement in the process.Run is not the best version of readable code, can be changed to a loop with checks.

### Example

```go

package main

import (
	"fmt"
	"time"

	"github.com/yolossn/process-manager/pkg/config"
	"github.com/yolossn/process-manager/pkg/manager"
)

func main() {

	conf, err := config.FromYaml("./config.yaml")
	if err != nil {
		panic(err)
	}

	man := manager.New(conf.Commands)

	now := time.Now()

	done := man.Run()

	time.Sleep(time.Second * 3)

	man.Stop()
	<-done

	fmt.Println("Successful process count: ", man.SuccessCount())
	fmt.Println("Failed process count: ", man.FailCount())

	// Done
	fmt.Println("done", time.Since(now))
}

```
