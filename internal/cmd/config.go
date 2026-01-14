package cmd

import (
	"errors"
	"fmt"
	"os"

	"aip/internal/config"
	"aip/internal/i18n"
	"github.com/spf13/cobra"
)

func newConfigCommand(lang i18n.Lang) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: i18n.T(lang, "cmd.config.short"),
	}

	cmd.AddCommand(newConfigShowCommand(lang))
	cmd.AddCommand(newConfigPathCommand(lang))
	cmd.AddCommand(newConfigGetCommand(lang))
	cmd.AddCommand(newConfigSetCommand(lang))
	cmd.AddCommand(newConfigWizardCommand(lang))
	return cmd
}

func newConfigShowCommand(lang i18n.Lang) *cobra.Command {
	var merged bool
	cmd := &cobra.Command{
		Use:   "show",
		Short: i18n.T(lang, "cmd.config.show.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			var cfg config.Config
			if merged {
				cfg, err = config.LoadMerged(path)
				if err != nil {
					return err
				}
			} else {
				cfg, err = config.Load(path)
				if err != nil {
					if errors.Is(err, os.ErrNotExist) {
						return fmt.Errorf("%s: %s", path, i18n.T(lang, "err.config_missing"))
					}
					return err
				}
			}
			fmt.Fprint(cmd.OutOrStdout(), renderConfig(cfg))
			return nil
		},
	}
	cmd.Flags().BoolVar(&merged, "merged", false, i18n.T(lang, "cmd.config.show.merged"))
	return cmd
}

func newConfigWizardCommand(lang i18n.Lang) *cobra.Command {
	return &cobra.Command{
		Use:   "wizard",
		Short: i18n.T(lang, "cmd.config.wizard.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			var current config.Config
			if cfg, err := config.Load(path); err == nil {
				current = cfg
			}
			cfg, err := config.Wizard(cmd.InOrStdin(), cmd.ErrOrStderr(), current)
			if err != nil {
				return err
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "%s\n", i18n.T(lang, "msg.config_saved"))
			return nil
		},
	}
}

func renderConfig(cfg config.Config) string {
	return fmt.Sprintf("base_url = %q\napi_key = %q\nmodel = %q\n", cfg.BaseURL, cfg.APIKey, cfg.Model)
}

func newConfigPathCommand(lang i18n.Lang) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: i18n.T(lang, "cmd.config.path.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), path)
			return nil
		},
	}
}

func newConfigGetCommand(lang i18n.Lang) *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: i18n.T(lang, "cmd.config.get.short"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			cfg, err := config.Load(path)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("%s: %s", path, i18n.T(lang, "err.config_missing"))
				}
				return err
			}
			val, ok := config.GetByKey(cfg, args[0])
			if !ok {
				return fmt.Errorf("%s: %s", args[0], i18n.T(lang, "err.config_key"))
			}
			fmt.Fprintln(cmd.OutOrStdout(), val)
			return nil
		},
	}
}

func newConfigSetCommand(lang i18n.Lang) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: i18n.T(lang, "cmd.config.set.short"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.DefaultPath()
			if err != nil {
				return err
			}
			var cfg config.Config
			if loaded, err := config.Load(path); err == nil {
				cfg = loaded
			}
			if !config.SetByKey(&cfg, args[0], args[1]) {
				return fmt.Errorf("%s: %s", args[0], i18n.T(lang, "err.config_key"))
			}
			if err := config.Save(path, cfg); err != nil {
				return err
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "%s\n", i18n.T(lang, "msg.config_saved"))
			return nil
		},
	}
}
