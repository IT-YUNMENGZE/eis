package register

import (
	"Edge-Scheduler/pkg/plugins/computeintensive"
	"Edge-Scheduler/pkg/plugins/dataintensive"
	"Edge-Scheduler/pkg/plugins/delaysensitive"

	"github.com/spf13/cobra"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func Register() *cobra.Command {
	return app.NewSchedulerCommand(
		app.WithPlugin(computeintensive.Name, computeintensive.New),
		app.WithPlugin(dataintensive.Name, dataintensive.New),
		app.WithPlugin(delaysensitive.Name, delaysensitive.New),
	)
}
