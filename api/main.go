package main

import (
	"os"

	_ "github.com/NdoleStudio/superbutton/docs"
	"github.com/NdoleStudio/superbutton/pkg/di"
)

// Version is the version of the API
var Version string

// @title       Superbutton
// @version     1.0
// @description Backend API to construct floating buttons
//
// @contact.name  Acho Arnold
// @contact.email support@superbutton.app
//
// @license.name MIT
// @license.url  https://raw.githubusercontent.com/NdoleStudio/superbutton/main/LICENSE
//
// @host     api.superbutton.app
// @schemes  https
// @BasePath /v1
//
// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if len(os.Args) == 1 {
		di.LoadEnv()
	}

	container := di.NewContainer(Version, os.Getenv("GCP_PROJECT_ID"))
	container.Logger().Info(container.App().Listen(":8000").Error())
}
