package cmd

import (
	"fmt"
	"runtime"

	"aip/internal/i18n"
	"aip/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCommand(lang i18n.Lang) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: i18n.T(lang, "cmd.version.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "aip %s\n", version.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "commit: %s\n", version.Commit)
			fmt.Fprintf(cmd.OutOrStdout(), "built: %s\n", version.BuildTime)
			fmt.Fprintf(cmd.OutOrStdout(), "go: %s\n", runtime.Version())
			fmt.Fprintf(cmd.OutOrStdout(), "default_model: %s\n", version.DefaultModel)
			return nil
		},
	}
}
