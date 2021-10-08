package main

import (
	"fiber-graphql/internal/core/config"
	"fiber-graphql/internal/core/sql"
	"fiber-graphql/internal/core/utils"
	"fiber-graphql/internal/handlers/routes"
	"flag"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/sirupsen/logrus"
)

func main() {
	configs := flag.String("config", "configs", "set configs path, default as: 'configs'")
	flag.Parse()

	// Init configuration
	err := config.InitConfig(*configs)
	if err != nil {
		panic(err)
	}
	//=======================================================

	// set logrus
	logrus.SetReportCaller(true)
	if config.CF.App.Release {
		logrus.SetFormatter(stackdriver.NewFormatter(
			stackdriver.WithService("api"),
			stackdriver.WithVersion("v1.0.0")))
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	logrus.Infof("Initial 'Configuration'. %+v", config.CF)
	//=======================================================

	// Init return result
	err = config.InitReturnResult("configs")
	if err != nil {
		panic(err)
	}
	//=======================================================

	// Get public key && private key (JWT)
	err = utils.ReadECDSAKey(config.CF.JWT.PrivateKeyPath, config.CF.JWT.PublicKeyPath)
	if err != nil {
		panic(err)
	}
	// ======================================================

	// Init connection sql
	err = sql.InitConnectionDatabase(config.CF.SQL)
	if err != nil {
		panic(err)
	}

	if !config.CF.App.Release {
		sql.Debug()
	}
	// ======================================================

	routes.New()
}
