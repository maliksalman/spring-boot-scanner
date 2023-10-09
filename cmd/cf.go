package cmd

import (
	"github.com/maliksalman/spring-boot-scanner/cf"
	"github.com/spf13/cobra"
)

func NewCmdCf() *cobra.Command {
	return &cobra.Command{
		Use:   "cf",
		Short: "Scans a CloudFoundry foundation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {

			apps := cf.ScanForApps()
			scannerOutput := cf.ScanAppContents(apps)
			scannerOutput.WriteAsJSON("output.json")

			return nil
		},
	}
}
