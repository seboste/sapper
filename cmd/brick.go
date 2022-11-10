package cmd

import (
	"github.com/spf13/cobra"
)

var brickCmd = &cobra.Command{
	Use:   "brick",
	Short: "Manage building bricks of your C++ microservice",
}

var addBrickCmd = &cobra.Command{
	Use:   "add [template]",
	Short: "Adds another building brick to the C++ microservice",
	Run: func(cmd *cobra.Command, args []string) {
		brickApi.Add()
	},
}

var listBrickCmd = &cobra.Command{
	Use:   "list [template]",
	Short: "Displays information about building bricks",
	Run: func(cmd *cobra.Command, args []string) {
		brickApi.List()
	},
}

var searchBrickCmd = &cobra.Command{
	Use:   "search [template]",
	Short: "Searches for building bricks",
	Run: func(cmd *cobra.Command, args []string) {
		brickApi.Search()
	},
}

func init() {
	brickCmd.AddCommand(addBrickCmd)
	brickCmd.AddCommand(listBrickCmd)
	brickCmd.AddCommand(searchBrickCmd)

	rootCmd.AddCommand(brickCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// brickCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// brickCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
