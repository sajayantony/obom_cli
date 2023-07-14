package cmd

import (
	"fmt"
	"os"

	obom "github.com/Azure/obom/internal"
	"github.com/spf13/cobra"
)

type packagesOptions struct {
	filename string
}

func packagesCmd() *cobra.Command {
	var opts packagesOptions
	var packagesCmd = &cobra.Command{
		Use:   "packages",
		Short: "List packages the SBOM",
		Long:  `List packages the SBOM that have external refs`,
		Run: func(cmd *cobra.Command, args []string) {
			sbom, _, err := obom.LoadSBOM(opts.filename)
			if err != nil {
				fmt.Println("Error loading SBOM:", err)
				os.Exit(1)
			}

			packages, err := obom.GetPackages(sbom)
			if err != nil {
				fmt.Println("Error getting packages:", err)
				os.Exit(1)
			}

			for _, pkg := range packages {
				fmt.Println(pkg)
			}
		},
	}

	packagesCmd.Flags().StringVarP(&opts.filename, "file", "f", "", "Path to the SPDX SBOM file")
	packagesCmd.MarkFlagRequired("file")

	return packagesCmd
}
