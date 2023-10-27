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

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Provides repo PR summary for last week",
	RunE:  repoSummary,
}

func init() {
	rootCmd.AddCommand(summaryCmd)
	initializeFlags(summaryCmd)
}

func repoSummary(cmd *cobra.Command, args []string) error {
	getEnvValues()

	client := github.NewClient(nil)
	options := &github.PullRequestListOptions{
		State:       "all",
		Sort:        "created",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var openCount, closedCount, inProgressCount int
	var openTitles, closedTitles, inProgressTitles []string
	for {
		pulls, resp, err := client.PullRequests.List(context.TODO(), repoOwner, repoName, options)
		if err != nil {
			return err
		}
		fmt.Println("PULL COUNT", len(pulls))

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
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	fmt.Printf("From:: %s\n", fromEmail)
	fmt.Printf("To:: %s\n", toEmail)
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

	return nil
}
