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
