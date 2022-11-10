package cmd

import (
	"github.com/spf13/cobra"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage sources for service templates and building bricks",
}

var addRemoteCmd = &cobra.Command{
	Use:   "add [template]",
	Short: "Adds a remote",
	Run: func(cmd *cobra.Command, args []string) {
		remoteApi.Add()
	},
}

var listRemoteCmd = &cobra.Command{
	Use:   "list [template]",
	Short: "Displays information about building bricks",
	Run: func(cmd *cobra.Command, args []string) {
		remoteApi.List()
	},
}

func init() {
	remoteCmd.AddCommand(addRemoteCmd)
	remoteCmd.AddCommand(listRemoteCmd)

	rootCmd.AddCommand(remoteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// brickCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// brickCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
