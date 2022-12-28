package cmd

import (
	"fmt"

	"github.com/seboste/sapper/adapters"
	"github.com/seboste/sapper/ports"
	"github.com/spf13/cobra"
)

func Print(b ports.Brick) {
	fmt.Println(b.Id, b.Version, b.Description)
}

var brickCmd = &cobra.Command{
	Use:   "brick",
	Short: "Manage building bricks of your C++ microservice",
}

var addBrickCmd = &cobra.Command{
	Use:   "add [template]",
	Short: "Adds another building brick to the C++ microservice",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("brick id argument is missing")
			return
		}

		brickId := args[0]
		service, _ := cmd.Flags().GetString("service")

		r, err := adapters.MakeSapperParameterResolver(cmd.Flags(), "")
		if err != nil {
			fmt.Println(err)
			return
		}

		err = brickApi.Add(service, brickId, r)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var listBrickCmd = &cobra.Command{
	Use:   "list [template]",
	Short: "Displays information about building bricks",
	Run: func(cmd *cobra.Command, args []string) {
		bricks := brickApi.List()
		fmt.Printf("found a total of %d bricks\n", len(bricks))
		for _, b := range bricks {
			Print(b)
		}
	},
}

var searchBrickCmd = &cobra.Command{
	Use:   "search <term>",
	Short: "Searches for building bricks in the id or description",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Unable to search for building bricks. The search term is missing")
			return
		}
		bricks := brickApi.Search(args[0])
		fmt.Printf("found a total of %d bricks\n", len(bricks))
		for _, b := range brickApi.Search(args[0]) {
			Print(b)
		}
	},
}

func init() {
	brickCmd.AddCommand(addBrickCmd)
	brickCmd.AddCommand(listBrickCmd)
	brickCmd.AddCommand(searchBrickCmd)

	rootCmd.AddCommand(brickCmd)

	addBrickCmd.PersistentFlags().StringP("service", "s", ".", "Path to the service that the brick shall be added to.")
	adapters.RegisterSapperParameterResolver(addBrickCmd.PersistentFlags())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// brickCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// brickCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
