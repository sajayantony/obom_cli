package cmd

import (
	"fmt"
	"os"

	obom "github.com/sajayantony/obom/internal"
	"github.com/spf13/cobra"
)

type showOptions struct {
	filename string
}

func showCmd() *cobra.Command {
	var opts showOptions
	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show the SBOM",
		Long:  `Show the SBOM`,
		Run: func(cmd *cobra.Command, args []string) {
			sbom, err := obom.LoadSBOM(opts.filename)
			if err != nil {
				fmt.Println("Error loading SBOM:", err)
				os.Exit(1)
			}

			obom.PrintSBOMSummary(sbom)
		},
	}

	showCmd.Flags().StringVarP(&opts.filename, "file", "f", "", "Path to the SPDX SBOM file")
	showCmd.MarkFlagRequired("file")

	return showCmd
}
