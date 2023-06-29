package cmd

import (
	"errors"
	"fmt"
	"os"

	parameterResolver "github.com/seboste/sapper/adapters/parameter-resolver"
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
	Use:           "add [template]",
	Short:         "Adds another building brick to the C++ microservice",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("brick id argument is missing")
		}
		brickId := args[0]
		service, _ := cmd.Flags().GetString("service")
		r, err := parameterResolver.MakeSapperParameterResolver(cmd.Flags(), "")
		if err != nil {
			return err
		}
		return brickApi.Add(service, brickId, r)
	},
}

var upgradeBrickCmd = &cobra.Command{
	Use:           "upgrade [brickId]",
	Short:         "upgrades the dependencies of a brick",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("brick id argument is missing")
		}
		brickId := args[0]
		return brickApi.Upgrade(brickId)
	},
}

var listBrickCmd = &cobra.Command{
	Use:           "list [template]",
	Short:         "Displays information about building bricks",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		bricks := brickApi.List()
		fmt.Printf("found a total of %d bricks\n", len(bricks))
		for _, b := range bricks {
			Print(b)
		}
	},
}

var searchBrickCmd = &cobra.Command{
	Use:           "search <term>",
	Short:         "Searches for building bricks in the id or description",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Unable to search for building bricks. The search term is missing")
		}
		bricks := brickApi.Search(args[0])
		fmt.Printf("found a total of %d bricks\n", len(bricks))
		for _, b := range brickApi.Search(args[0]) {
			Print(b)
		}
		return nil
	},
}

var describeBrickCmd = &cobra.Command{
	Use:           "describe [brickId]",
	Short:         "Shows a brick's README.md for information about the usage, its parameters, and next steps",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("brick id argument is missing")
		}
		brickId := args[0]
		return brickApi.Describe(brickId, os.Stdout)
	},
}

func init() {
	brickCmd.AddCommand(addBrickCmd)
	brickCmd.AddCommand(upgradeBrickCmd)
	brickCmd.AddCommand(listBrickCmd)
	brickCmd.AddCommand(searchBrickCmd)
	brickCmd.AddCommand(describeBrickCmd)

	rootCmd.AddCommand(brickCmd)

	addBrickCmd.PersistentFlags().StringP("service", "s", ".", "Path to the service that the brick shall be added to.")
	parameterResolver.RegisterSapperParameterResolver(addBrickCmd.PersistentFlags())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// brickCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// brickCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
