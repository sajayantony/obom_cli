package cmd

import (
	"fmt"
	"os"

	obom "github.com/sajayantony/obom/internal"
	"github.com/spf13/cobra"

	"oras.land/oras-go/v2/registry"
)

type pushOpts struct {
	filename  string
	reference string
	username  string
	password  string
}

func pushCmd() *cobra.Command {
	var opts pushOpts
	var pushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push the SBOM",
		Long:  `Push the SBOM with the annotations`,
		Run: func(cmd *cobra.Command, args []string) {

			// get the reference as the first argument
			opts.reference = args[0]

			// validate if reference is valid
			_, err := registry.ParseReference(opts.reference)
			if err != nil {
				fmt.Println("Error parsing reference:", err)
				os.Exit(1)
			}

			sbom, err := obom.LoadSBOM(opts.filename)
			if err != nil {
				fmt.Println("Error loading SBOM:", err)
				os.Exit(1)
			}

			annotations, err := obom.GetAnnotations(sbom)
			if err != nil {
				fmt.Println("Error getting annotations:", err)
				os.Exit(1)
			}

			err = obom.PushFiles(opts.filename, opts.reference, annotations, opts.username, opts.password)
			if err != nil {
				fmt.Println("Error pushing SBOM:", err)
				os.Exit(1)
			}

		},
	}

	pushCmd.Flags().StringVarP(&opts.filename, "file", "f", "", "Path to the SPDX SBOM file")
	pushCmd.MarkFlagRequired("file")

	pushCmd.Flags().StringVarP(&opts.username, "username", "u", "", "Username for the registry")
	pushCmd.Flags().StringVarP(&opts.password, "password", "p", "", "Password for the registry")

	// Add positional argument called reference to pushCmd
	pushCmd.Args = cobra.ExactArgs(1)

	return pushCmd
}
