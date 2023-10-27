/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	fromEmail string
	toEmail   string
	repoName  string
	repoOwner string
)

func initializeFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&fromEmail, "fromEmail", "f", "abc@gmail.com", "From email")
	cmd.Flags().StringVarP(&toEmail, "toEmail", "t", "xyz@gmail.com", "To email")
	cmd.Flags().StringVar(&repoName, "repoName", "", "Public repo Name")
	cmd.Flags().StringVar(&repoOwner, "repoOwner", "", "Public repo Owner")
}

func getEnvValues() {
	if fromEmail == "" {
		fromEmail = os.Getenv("FROM_EMAIL")
	}
	if toEmail == "" {
		toEmail = os.Getenv("TO_EMAIL")
	}
	if repoName == "" {
		repoName = os.Getenv("REPO_NAME")
	}
	if repoOwner == "" {
		repoOwner = os.Getenv("REPO_OWNER")
	}
}
