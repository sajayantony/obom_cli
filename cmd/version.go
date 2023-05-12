package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/sajayantony/obom/internal/version"
	"github.com/spf13/cobra"
)

// -ldflags="-X 'github.com/sajayantony/obom.cmd.Version=$TAG'"

func init() {
	if version.Tag == "" {
		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {
			// https://github.com/golang/go/issues/51831#issuecomment-1074188363
			return
		}
		version.Tag = buildInfo.Main.Version
	}
}

func versionCmd() *cobra.Command {
	// Create version command for cobra root command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  `Show the version of obom`,
		Run: func(cmd *cobra.Command, args []string) {
			if version.Tag == "" {
				fmt.Println("Error getting version")
				return
			} else {
				fmt.Printf("Version:	%s\n", version.Tag)
				fmt.Printf("Git Commit:	%s\n", version.GitCommit)
				fmt.Printf("Git Tree State:	%s\n", version.GitTreeState)
			}
		},
	}

	return versionCmd
}
