package cmd

import (
	"github.com/edgelesssys/constellation/internal/constants"
	"github.com/spf13/cobra"
)

// NewVerifyCmd returns a new cobra.Command for the verify command.
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version of this CLI",
		Long:  "Display version of this CLI.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("CLI Version: v%s \n", constants.CliVersion)
		},
	}
	return cmd
}