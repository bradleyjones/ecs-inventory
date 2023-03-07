package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/anchore/anchore-ecs-inventory/ecg"
	"github.com/anchore/anchore-ecs-inventory/internal/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "anchore-ecs-inventory",
	Short: "anchore-ecs-inventory tells Anchore which images are in use in your ECS clusters",
	Long:  "anchore-ecs-inventory can poll Amazon ECS (Elastic Container Service) APIs to tell Anchore which Images are currently in-use",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			err := cmd.Help()
			if err != nil {
				log.Error("error running help command", err)
				os.Exit(1)
			}
			os.Exit(1)
		}

		// TODO(bradjones) Validate anchore connection details here
		/*
			if appConfig.AnchoreDetails.IsValid() {
				dummyReport := inventory.Report{
					Results: []inventory.ReportItem{},
				}
				err := reporter.Post(dummyReport, appConfig.AnchoreDetails, appConfig)
				if err != nil {
					log.Error("Failed to validate connection to Anchore", err)
				}
			} else {
				log.Debug("Anchore details not specified, will not report inventory")
			}
		*/

		ecg.PeriodicallyGetInventoryReport(appConfig.PollingIntervalSeconds, appConfig.AnchoreDetails, appConfig.Region)
	},
}

func init() {
	opt := "polling-interval-seconds"
	rootCmd.Flags().
		StringP(opt, "p", strconv.Itoa(config.DefaultConfigValues.PollingIntervalSeconds), "This specifies the polling interval of the ECS API in seconds")
	if err := viper.BindPFlag(opt, rootCmd.Flags().Lookup(opt)); err != nil {
		fmt.Printf("unable to bind flag '%s': %+v", opt, err)
		os.Exit(1)
	}

	opt = "region"
	rootCmd.Flags().
		StringP(opt, "r", config.DefaultConfigValues.Region, "If set overrides the AWS_REGION environment variable/region specified in anchore-ecs-inventory config")
	if err := viper.BindPFlag(opt, rootCmd.Flags().Lookup(opt)); err != nil {
		fmt.Printf("unable to bind flag '%s': %+v", opt, err)
		os.Exit(1)
	}
}
