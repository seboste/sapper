package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage sources for service templates and building bricks",
}

var addRemoteCmd = &cobra.Command{
	Use:   "add remote_name remote_url [--insert=position]",
	Short: "Add a new file based or git remote",

	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			fmt.Println("remote_name and/or remote_src arguments are missing")
			return
		}
		position, _ := cmd.Flags().GetInt("insert")
		if err := remoteApi.Add(args[0], args[1], position); err != nil {
			fmt.Println(err)
		}
	},
}

var removeRemoteCmd = &cobra.Command{
	Use:   "remove remote_name",
	Short: "Remove a remote",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("remote_name argument is missing")
			return
		}
		if err := remoteApi.Remove(args[0]); err != nil {
			fmt.Println(err)
		}
	},
}

var updateRemoteCmd = &cobra.Command{
	Use:   "update git_remote_name",
	Short: "Pulls latest version for git remotes. No effect on file based remotes.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("remote_name argument is missing")
			return
		}
		if err := remoteApi.Update(args[0]); err != nil {
			fmt.Println(err)
		}
	},
}

var upgradeRemoteCmd = &cobra.Command{
	Use:   "upgrade remote_name",
	Short: "Upgrades the dependencies of all bricks and templates in a remote.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("remote_name argument is missing")
			return
		}
		if err := remoteApi.Upgrade(args[0]); err != nil {
			fmt.Println(err)
		}
	},
}

var listRemoteCmd = &cobra.Command{
	Use:   "list",
	Short: "List current remotes",
	Run: func(cmd *cobra.Command, args []string) {
		for _, r := range remoteApi.List() {
			fmt.Printf("%s: %s\n", r.Name, r.Src)
		}
	},
}

func init() {
	remoteCmd.AddCommand(addRemoteCmd)
	remoteCmd.AddCommand(removeRemoteCmd)
	remoteCmd.AddCommand(updateRemoteCmd)
	remoteCmd.AddCommand(upgradeRemoteCmd)
	remoteCmd.AddCommand(listRemoteCmd)

	rootCmd.AddCommand(remoteCmd)

	addRemoteCmd.Flags().IntP("insert", "i", -1, "Insert the remote at a given position")
}
