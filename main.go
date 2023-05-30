package main

import (
	"fmt"
	"os"
	"k8s.io/component-base/cli"
	//"context"
	//"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"github.com/mariopetricevic/plugins"
	//runtime "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
)

func main() {

	
	command := app.NewSchedulerCommand(
		app.WithPlugin(plugins.Name, plugins.New),
	)
	
	code := cli.Run(command)
	
	if(code == 1){
		fmt.Println("registered plugin successfully.");
	}
	
	os.Exit(code)
}
