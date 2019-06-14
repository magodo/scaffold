package main

//go:generate go run pkg/gen-template.go template

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var (
		templateValue struct {
			// Modify following members as you wish
			GreetMsg string
		}

		outputDir string
	)
	var rootCmd = &cobra.Command{
		Use:   "scaffold",
		Short: "Create greeting scaffold",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := GenScaffold(outputDir, templateValue); err != nil {
				log.Fatal(err)
			}
		},
	}
	rootCmd.Flags().StringVarP(&outputDir, "output_dir", "o", "scaffold", "outputdir of project scaffold")

	// Modify below lings also
	rootCmd.Flags().StringVarP(&templateValue.GreetMsg, "msg", "m", "hello...", "greeting message")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
