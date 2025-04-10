package main

import (
	"fmt"

	server "github.com/MenD32/allpaca/pkg/server"
	config "github.com/MenD32/allpaca/pkg/server/config"
	"github.com/spf13/cobra"
)

var (
	Version        = "dev"
	CommitHash     = "none"
	BuildTimestamp = "unknown"
)

func BuildVersion() string {
	return fmt.Sprintf("%s-%s (%s)", Version, CommitHash, BuildTimestamp)
}

func main() {
	var configFile string

	rootCmd := &cobra.Command{
		Use:   "allpaca",
		Short: "Allpaca is a server application",
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of Allpaca",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(BuildVersion())
		},
	}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the Allpaca server",
		Run: func(cmd *cobra.Command, args []string) {
			var c *config.Config
			var err error

			if configFile != "" {
				c, err = config.ParseConfigFromFile(configFile)
				if err != nil {
					fmt.Printf("Error loading config: %v\n", err)
					return
				}
			} else {
				c = config.NewRecommendedConfig()
			}

			s := server.NewServer(c)
			s.Start()
		},
	}

	runCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the configuration file (optional, defaults to none)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
