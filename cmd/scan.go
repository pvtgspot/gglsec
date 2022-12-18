package cmd

import (
	"github.com/pvtgspot/gglsec/internal/gglsec"
	"github.com/pvtgspot/gglsec/internal/gglsec/rules"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

const (
	FLAG_NAME_TOKEN    = "token"
	FLAG_NAME_ENDPOINT = "endpoint"
	FLAG_NAME_GROUPS   = "groups"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "GitLab configuration security scan made with Go",
	Long: `gglsec - GitLab configuration security scanner.

It scans the configuration of projects and groups of a GitLab instance for compliance with security rules
and outputs a summary of the scan results to the console`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := cmd.Flags().GetString(FLAG_NAME_TOKEN)
		if err != nil {
			panic(err)
		}
		endpoint, err := cmd.Flags().GetString(FLAG_NAME_ENDPOINT)
		if err != nil {
			panic(err)
		}

		groups, err := cmd.Flags().GetStringSlice(FLAG_NAME_GROUPS)
		if err != nil {
			panic(err)
		}

		gitlabClient, err := gitlab.NewClient(
			token,
			gitlab.WithBaseURL(endpoint),
		)
		if err != nil {
			panic(err)
		}

		rl := gglsec.NewRuleList()
		for _, gid := range groups {
			rl.Append(
				rules.NewGroupBranchProtectionRule(gid, gitlabClient),
				rules.NewVisibilityLevelRule(gid, gitlabClient),
			)
		}

		ss := gglsec.NewGitlabConfigScanner(rl)
		scanResult := ss.Run()

		scanResult.PrintReport()
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringP(FLAG_NAME_TOKEN, "t", "", "GitLab API token (required)")
	scanCmd.MarkFlagRequired(FLAG_NAME_TOKEN)
	scanCmd.Flags().StringP(FLAG_NAME_ENDPOINT, "e", "https://gitlab.com/api/v4/", "GitLab API endpoint")
	scanCmd.Flags().StringSliceP(FLAG_NAME_GROUPS, "g", make([]string, 0), "List of the groups id to scan")
}
