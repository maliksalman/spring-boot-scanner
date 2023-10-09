package cmd

import "github.com/spf13/cobra"

func NewCmdRoot() *cobra.Command {

	root := &cobra.Command{
		Use:   "spring-boot-scanner",
		Short: "spring-boot-scanner is a tool to scan an app platform for spring-boot applications",
		RunE:  func(cmd *cobra.Command, _ []string) error { return cmd.Usage() },
	}

	root.AddCommand(
		NewCmdCf(),
		NewCmdK8s(),
	)
	return root
}
