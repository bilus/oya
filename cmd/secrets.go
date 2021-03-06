package cmd

import (
	"os"

	"github.com/bilus/oya/cmd/internal"
	"github.com/spf13/cobra"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage secrets in Oyafile.secrets",
}

var secretsViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		return internal.SecretsView(cwd, cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

var secretsEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit secrets file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		return internal.SecretsEdit(cwd, cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

var secretsEncryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt secrets file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		return internal.SecretsEncrypt(cwd, cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

func init() {
	secretsCmd.AddCommand(secretsViewCmd)
	secretsCmd.AddCommand(secretsEditCmd)
	secretsCmd.AddCommand(secretsEncryptCmd)
	rootCmd.AddCommand(secretsCmd)
}
