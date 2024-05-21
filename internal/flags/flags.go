package flags

import "github.com/spf13/cobra"

type Flags struct {
	CommonAncestor *string
	CurrentVersion *string
	OtherVersion   *string
}

func AddFlags(cmd *cobra.Command) Flags {
	flags := cmd.Flags()

	return Flags{
		CommonAncestor: flags.StringP("common-ancestor", "O", "", "Common ancestor file"),
		CurrentVersion: flags.StringP("current-version", "A", "", "Current version file"),
		OtherVersion:   flags.StringP("other-version", "B", "", "Other version file"),
	}
}
