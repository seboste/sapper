package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	parameterResolver "github.com/seboste/sapper/adapters/parameter-resolver"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage C++ microservices",
}

var keepMajorVersion *bool

var addServiceCmd = &cobra.Command{
	Use:           "add [folder]",
	Short:         "Adds a new C++ microservice",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("service folder argument is missing")
		}

		path, name := filepath.Split(args[0])
		template, _ := cmd.Flags().GetString("template")

		r, err := parameterResolver.MakeSapperParameterResolver(cmd.Flags(), name)
		if err != nil {
			return err
		}

		_, err = serviceApi.Add(template, path, r)
		return err
	},
}

var describeServiceCmd = &cobra.Command{
	Use:           "describe [service folder]",
	Short:         "Prints information about a service",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("service folder argument is missing")
		}
		return serviceApi.Describe(args[0], os.Stdout)
	},
}

var upgradeServiceCmd = &cobra.Command{
	Use:           "upgrade [service folder]",
	Short:         "upgrades the dependencies of the service",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.PersistentFlags()
		return serviceApi.Upgrade(args[0], *keepMajorVersion)
	},
}

var buildServiceCmd = &cobra.Command{
	Use:           "build [service folder]",
	Short:         "Builds the service",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return errors.New("service folder argument is missing")
		}

		fmt.Printf("building service...")
		buildLogFilename, err := serviceApi.Build(args[0])
		if err != nil {
			fmt.Printf("failed (see %s for details)\n", buildLogFilename)
			return err
		} else {
			fmt.Println("success")
		}
		return nil
	},
}

var testServiceCmd = &cobra.Command{
	Use:           "test [service folder]",
	Short:         "Tests the service",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("service folder argument is missing")
		}

		fmt.Printf("tesing service...")
		err := serviceApi.Test(args[0])
		if err != nil {
			return err
		} else {
			fmt.Println("success")
		}
		return nil
	},
}

var deployServiceCmd = &cobra.Command{
	Use:           "deploy [service folder]",
	Short:         "Deploy the service",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("service folder argument is missing")
		}

		fmt.Printf("deploying service...")
		err := serviceApi.Deploy(args[0])
		if err != nil {
			return err
		} else {
			fmt.Println("success")
		}
		return nil
	},
}

var runServiceCmd = &cobra.Command{
	Use:           "run [service folder]",
	Short:         "Run the service",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("service folder argument is missing")
		}

		fmt.Printf("tesing service...")
		err := serviceApi.Run(args[0])
		if err != nil {
			return err
		} else {
			fmt.Println("success")
		}
		return nil
	},
}

func init() {
	serviceCmd.AddCommand(addServiceCmd)
	serviceCmd.AddCommand(describeServiceCmd)
	serviceCmd.AddCommand(upgradeServiceCmd)
	serviceCmd.AddCommand(buildServiceCmd)
	serviceCmd.AddCommand(testServiceCmd)
	serviceCmd.AddCommand(deployServiceCmd)
	serviceCmd.AddCommand(runServiceCmd)

	rootCmd.AddCommand(serviceCmd)

	addServiceCmd.PersistentFlags().StringP("template", "t", "base-hexagonal-skeleton", "The id of a service template.")
	parameterResolver.RegisterSapperParameterResolver(addServiceCmd.PersistentFlags())

	keepMajorVersion = upgradeServiceCmd.PersistentFlags().Bool("keep-major", false, "Upgrades are only conducted within the same major version of a dependency's semantic version")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// brickCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// brickCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
