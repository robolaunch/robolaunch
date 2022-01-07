package main

import (
	"log"
	"os"

	launchflow "github.com/robolaunch/robolaunch/launch/pkg/workflow"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.NewClient(client.Options{
		HostPort: os.Getenv("TEMPORAL_SERVER_IP"),
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	w := worker.New(c, launchflow.LaunchQueue, worker.Options{})
	w.RegisterWorkflow(launchflow.LaunchWorkflow)
	w.RegisterActivity(launchflow.CreateLaunch)
	w.RegisterActivity(launchflow.DeleteLaunch)
	w.RegisterActivity(launchflow.ScaleOut)
	w.RegisterActivity(launchflow.ScaleUp)
	w.RegisterActivity(launchflow.CreateUserSpace)
	w.RegisterActivity(launchflow.CreateLaunchHelm)
	w.RegisterActivity(launchflow.DeleteLaunchHelm)
	w.RegisterActivity(launchflow.ScaleDownHelm)
	w.RegisterActivity(launchflow.ScaleUpHelm)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Worker cannot start")
	}
}
