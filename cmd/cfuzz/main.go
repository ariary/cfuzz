package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cfuzz [flags] [command]",
	Short: "Fuzz any command line execution",
	Long: `cfuzz fuzzes any shell command by replacing a keyword with wordlist entries.

Put FUZZ in your command, provide a wordlist, and cfuzz runs every substitution concurrently.

Examples:
  cfuzz -w passwords.txt echo FUZZ
  cfuzz -w users.txt --success ssh FUZZ@host id
  cfuzz wordlist "SSH usernames" | cfuzz --stdin-wordlist "ssh FUZZ@host id"`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func main() {
	log.SetFlags(0)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
