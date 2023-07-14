package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	obom "github.com/Azure/obom/internal"
	"github.com/spf13/cobra"

	"oras.land/oras-go/v2/registry"
)

type pushOpts struct {
	filename            string
	reference           string
	username            string
	password            string
	ManifestAnnotations []string
}

var (
	errAnnotationFormat      = errors.New("missing key in `--annotation` flag")
	errAnnotationDuplication = errors.New("duplicate annotation key")
)

func pushCmd() *cobra.Command {
	var opts pushOpts
	var pushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push the SPDX SBOM to the registry",
		Long: `Push the SDPX with the annotations to an OCI registry

Example - Push an SPDX SBOM to a registry
	obom push -f spdx.json localhost:5000/spdx:latest 

Example - Push an SPDX SBOM to a registry with annotations
	obom push -f spdx.json localhost:5000/spdx:latest --annotation key1=value1 --annotation key2=value2

Example - Push an SPDX SBOM to a registry with annotations and credentials
	obom push -f spdx.json localhost:5000/spdx:latest --annotation key1=value1 --annotation key2=value2 --username user --password pass
`,
		Run: func(cmd *cobra.Command, args []string) {

			// get the reference as the first argument
			opts.reference = args[0]

			// validate if reference is valid
			_, err := registry.ParseReference(opts.reference)
			if err != nil {
				fmt.Println("Error parsing reference:", err)
				os.Exit(1)
			}

			sbom, desc, err := obom.LoadSBOM(opts.filename)
			if err != nil {
				fmt.Println("Error loading SBOM:", err)
				os.Exit(1)
			}

			obom.PrintSBOMSummary(sbom, desc)

			annotations, err := obom.GetAnnotations(sbom)
			if err != nil {
				fmt.Println("Error getting annotations:", err)
				os.Exit(1)
			}

			// parse the annotations from the flags and merge with annotations from the SBOM
			inputAnnotations, err := parseAnnotationFlags(opts.ManifestAnnotations)
			if err != nil {
				fmt.Println("Error parsing annotations:", err)
				os.Exit(1)
			}
			for k, v := range inputAnnotations {
				annotations[k] = v
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

	pushCmd.Flags().StringArrayVarP(&opts.ManifestAnnotations, "annotation", "a", nil, "manifest annotations")

	pushCmd.Flags().StringVarP(&opts.username, "username", "u", "", "Username for the registry")
	pushCmd.Flags().StringVarP(&opts.password, "password", "p", "", "Password for the registry")

	// Add positional argument called reference to pushCmd
	pushCmd.Args = cobra.ExactArgs(1)

	return pushCmd
}

func parseAnnotationFlags(flags []string) (map[string]string, error) {
	manifestAnnotations := make(map[string]string)
	for _, anno := range flags {
		key, val, success := strings.Cut(anno, "=")
		if !success {
			return nil, fmt.Errorf("%w: %s", errAnnotationFormat, anno)
		}
		if _, ok := manifestAnnotations[key]; ok {
			return nil, fmt.Errorf("%w: %v, ", errAnnotationDuplication, key)
		}
		manifestAnnotations[key] = val
	}
	return manifestAnnotations, nil
}
