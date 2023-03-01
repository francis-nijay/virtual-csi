package main

import (
	"context"

	"github.com/dell/gocsi"

	"github.com/dell/virtual-csi/provider"
	"github.com/dell/virtual-csi/service"
)

// main is ignored when this package is built as a go plug-in.
func main() {
	driverName := flag.String("driver-name", "vcsi-driver.dellemc.com", "driver name")
	driverSecret := flag.String("driver-secret", "/vcsi-driver-config/config", "driver secret file")
	driverConfig := flag.String("driver-config", "/vcsi-driver-config-params/driver-config-params.yaml", "array config params")
	enableLeaderElection := flag.Bool("leader-election", false, "boolean to enable leader election")

	flag.Parse()

	if *driverName == "" {
		fmt.Fprintf(os.Stderr, "driver-name argument is mandatory")
		os.Exit(1)
	}
	service.Name = *driverName

	if *driverConfig == "" {
		fmt.Fprintf(os.Stderr, "driver-config argument is mandatory")
		os.Exit(1)
	}
	service.DriverConfig = *driverConfig

	if *driverSecret == "" {
		fmt.Fprintf(os.Stderr, "driver-secret argument is mandatory")
		os.Exit(1)
	}
	service.DriverSecret = *driverSecret

	if !*enableLeaderElection {
		fmt.Println("leader election")
	}

	// Always set X_CSI_DEBUG to false irrespective of what user has specified
	_ = os.Setenv(gocsi.EnvVarDebug, "false")
	// We always want to enable Request and Response logging (no reason for users to control this)
	_ = os.Setenv(gocsi.EnvVarReqLogging, "true")
	_ = os.Setenv(gocsi.EnvVarRepLogging, "true")

	service.DriverConfig = *driverConfig
	service.DriverSecret = *driverSecret
	
	gocsi.Run(
		context.Background(),
		service.Name,
		"Dell vCSI Driver",
		"",
		provider.New())
}
