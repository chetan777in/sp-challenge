/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v56/github"
	"github.com/spf13/cobra"
)

var (
	fromEmail string
	toEmail   string
	repoName  string
	repoOwner string
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Provides repo PR summary for last week",
	RunE:  repoSummary,
}

func init() {
	rootCmd.AddCommand(summaryCmd)
	summaryCmd.Flags().StringVarP(&fromEmail, "fromEmail", "f", "abc@gmail.com", "From email")
	summaryCmd.Flags().StringVarP(&toEmail, "toEmail", "t", "xyz@gmail.com", "To email")
	summaryCmd.Flags().StringVar(&repoName, "repoName", "rN", "Public repo Name")
	summaryCmd.Flags().StringVar(&repoOwner, "repoOwner", "rO", "Public repo Owner")
}

func repoSummary(cmd *cobra.Command, args []string) error {
	client := github.NewClient(nil)
	options := &github.PullRequestListOptions{
		State:       "all",
		Sort:        "created",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 1000},
	}

	var openCount, closedCount, inProgressCount int
	var openTitles, closedTitles, inProgressTitles []string
	for {
		pulls, resp, err := client.PullRequests.List(context.TODO(), repoOwner, repoName, options)
		if err != nil {
			return err
		}

		for _, pull := range pulls {
			if pull.UpdatedAt.After(time.Now().AddDate(0, 0, -7)) {
				switch *pull.State {
				case "open":
					if pull.CreatedAt.After(time.Now().AddDate(0, 0, -7)) {
						openCount++
						openTitles = append(openTitles, *pull.Title)
					}
				case "closed":
					if (pull.MergedAt != nil && pull.MergedAt.After(time.Now().AddDate(0, 0, -7))) ||
						(pull.ClosedAt != nil && pull.ClosedAt.After(time.Now().AddDate(0, 0, -7))) {
						closedCount++
						closedTitles = append(closedTitles, *pull.Title)
					}
				default:
					inProgressCount++
					inProgressTitles = append(inProgressTitles, *pull.Title)
				}
			}
		}
		fmt.Printf("From:: %s", fromEmail)
		fmt.Printf("To:: %s", toEmail)
		fmt.Printf("Summary of repo %s/%s for last week\n", repoOwner, repoName)
		fmt.Printf("Open Pull Requests count = %d\n", openCount)
		for _, title := range openTitles {
			fmt.Printf("  %s\n", title)
		}
		fmt.Printf("Closed Pull Requests count = %d\n", closedCount)
		for _, title := range closedTitles {
			fmt.Printf("  %s\n", title)
		}
		fmt.Printf("In Progress Pull Requests count = %d\n", inProgressCount)
		for _, title := range inProgressTitles {
			fmt.Printf("  %s\n", title)
		}
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	return nil
}
