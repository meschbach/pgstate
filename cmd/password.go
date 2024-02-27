package main

import (
	"fmt"
	"github.com/meschbach/pgstate"
	"github.com/spf13/cobra"
)

func generatePasswordCommand() *cobra.Command {
	type flags struct {
		disableSpecial bool
	}
	passed := flags{}
	generatePassword := &cobra.Command{
		Use:  "generate-password",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			passwordCfg := pgstate.GenPasswordConfig{AllowSpecial: !passed.disableSpecial}
			password := pgstate.GeneratePasswordWithConfig(passwordCfg)
			fmt.Printf("Len %d -- %q\n", len(password), password)
		},
	}
	cmdFlags := generatePassword.Flags()
	cmdFlags.BoolVarP(&passed.disableSpecial, "allow-special", "s", false, "Disable using special characters")
	return generatePassword
}
