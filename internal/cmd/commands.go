package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yjhatfdu/aip/internal/i18n"
)

func newStubCommand(lang i18n.Lang, name, shortKey string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: i18n.T(lang, shortKey),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("%s: %s", name, i18n.T(lang, "err.not_implemented"))
		},
	}
}
