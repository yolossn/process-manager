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
	status := man.Status()
	for k, v := range status {
		fmt.Println("-----")
		fmt.Println("command:", k)
		fmt.Println("isSuccessful:", v.IsSuccessful)
		fmt.Println("output:", v.Output)
		fmt.Println("error:", v.Err)
		fmt.Println("-----")
	}
	// Done
	fmt.Println("done", time.Since(now))
}
