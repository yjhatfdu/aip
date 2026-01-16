package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yjhatfdu/aip/internal/i18n"
)

func Execute() error {
	return newRoot().Execute()
}

func newRoot() *cobra.Command {
	lang := i18n.Detect()
	root := &cobra.Command{
		Use:           "aip",
		Short:         i18n.T(lang, "root.short"),
		Long:          i18n.T(lang, "root.long"),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(
		newSummaryCommand(lang),
		// newStubCommand(lang, "map", "cmd.map.short"),
		// newStubCommand(lang, "watch", "cmd.watch.short"),
		newNormCommand(lang),
		// newStubCommand(lang, "reduce", "cmd.reduce.short"),
		newClusterCommand(lang),
		// newStubCommand(lang, "sample", "cmd.sample.short"),
		// newStubCommand(lang, "diagnose", "cmd.diagnose.short"),
		// newStubCommand(lang, "cache", "cmd.cache.short"),
		newConfigCommand(lang),
		newVersionCommand(lang),
	)

	return root
}
