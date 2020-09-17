<h1 align="center">Process Manager</h1>

## Design Decision

1. Manager provides the process with command, backoff strategy and maxRetries. The process takes care of retries internally.
   Why ?
   In this case the manager doesnt dynamically change the retry strategy or retry count based on the process outcome so the sole responsibility of completing the command with retries belongs to the process.
2. Both processes which have reached maxRetries and the ones which have run succesfully are considered complete.
3. If an Interupt or Kill is made when the processes are executing the manager kills the processes which have not completed yet and exits.

## Improvements

1. Improve tests (Add table tests, increase coverage)
2. Make backoff strategy configurable. Now it is a constant backoff of 1 second
3. The use of GOTO statement in the process.Run is not the best version of readable code, can be changed to a loop with checks.

## Example

### Config

```yaml
commands:
  - command: bash
    args:
      - -c
      - echo $PWD
    envs:
      - key: PWD
        value: /Users/santhoshnagarajs/git/
    maxRetries: 2
  - command: ld
    args:
      - $PWD
    envs:
      - key: PWD
        value: /Users/santhoshnagarajs/git/yolossn/proccess-manager
    maxRetries: 100
```

### Usage

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

## Test

> go test pkg/\* -v -cover -race

```
=== RUN   TestNewStaticBackoff
--- PASS: TestNewStaticBackoff (0.00s)
PASS
coverage: 100.0% of statements
ok      github.com/yolossn/process-manager/pkg/backoff  1.619s  coverage: 100.0% of statements
=== RUN   TestFromFile
=== PAUSE TestFromFile
=== CONT  TestFromFile
--- PASS: TestFromFile (0.00s)
PASS
coverage: 91.7% of statements
ok      github.com/yolossn/process-manager/pkg/config   1.410s  coverage: 91.7% of statements
=== RUN   TestNewManager
=== PAUSE TestNewManager
=== RUN   TestRun
=== PAUSE TestRun
=== CONT  TestRun
=== CONT  TestNewManager
--- PASS: TestNewManager (0.00s)
--- PASS: TestRun (6.08s)
    manager_test.go:31: 1
    manager_test.go:32: 3
PASS
coverage: 89.4% of statements
ok      github.com/yolossn/process-manager/pkg/manager  7.992s  coverage: 89.4% of statements
=== RUN   TestNew
=== PAUSE TestNew
=== RUN   TestRun
=== PAUSE TestRun
=== CONT  TestRun
=== CONT  TestNew
--- PASS: TestNew (0.00s)
--- PASS: TestRun (5.12s)
PASS
coverage: 96.4% of statements
ok      github.com/yolossn/process-manager/pkg/process  7.326s  coverage: 96.4% of statements
```
