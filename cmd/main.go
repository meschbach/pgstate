package main

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/meschbach/pgstate"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var pgxConnectionConfig string

	ensure := &cobra.Command{
		Use:   "ensure <database-name> <secret>",
		Short: "Creates a role and database by the same name with the given secret",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := pgx.ParseConfig(pgxConnectionConfig)
			if err != nil {
				return err
			}
			if err := pgstate.EnsureDatabase(cmd.Context(), config, args[0], args[1]); err != nil {
				fmt.Fprintf(os.Stderr, "Failed: %s\n", err.Error())
				os.Exit(-1)
			}
			return nil
		},
	}
	drop := &cobra.Command{
		Use:  "drop <database-name>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := pgx.ParseConfig(pgxConnectionConfig)
			if err != nil {
				return err
			}
			if err := pgstate.DestroyDatabase(cmd.Context(), config, args[0]); err != nil {
				fmt.Fprintf(os.Stderr, "Failed: %s\n", err.Error())
				os.Exit(-1)
			}
			if err := pgstate.DestroyRole(cmd.Context(), config, args[0]); err != nil {
				fmt.Fprintf(os.Stderr, "Failed: %s\n", err.Error())
				os.Exit(-2)
			}
			return nil
		},
	}

	root := &cobra.Command{
		Use:           "pgstate",
		SilenceErrors: true,
	}
	root.AddCommand(ensure)
	root.AddCommand(drop)
	root.AddCommand(generatePasswordCommand())
	rootFlags := root.PersistentFlags()
	rootFlags.StringVarP(&pgxConnectionConfig, "cluster", "c", "", "Connection string conforming to PGX DSN")

	if err := root.Execute(); err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "Failed %s", err.Error()); err != nil {
			panic(err)
		}
		os.Exit(-1)
	}
}
