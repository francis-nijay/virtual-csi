package main

import (
	"context"

	"github.com/dell/gocsi"

	"github.com/dell/virtual-csi/provider"
	"github.com/dell/virtual-csi/service"
)

// main is ignored when this package is built as a go plug-in.
func main() {
	gocsi.Run(
		context.Background(),
		service.Name,
		"A description of the SP",
		"",
		provider.New())
}
