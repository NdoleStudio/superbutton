package main

import (
	"os"

	"github.com/NdoleStudio/superbutton/pkg/di"
)

// Version is the version of the API
var Version string

func main() {
	if len(os.Args) == 1 {
		di.LoadEnv()
	}

	container := di.NewContainer(Version, os.Getenv("GCP_PROJECT_ID"))
	container.Logger().Info(container.App().Listen(":8000").Error())
}
