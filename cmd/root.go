package cmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/khos2ow/ratelimiter/internal/data"
	"github.com/khos2ow/ratelimiter/internal/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var backends []string
var options = &data.Options{}

var rootCmd = &cobra.Command{
	Args:         cobra.NoArgs,
	Use:          "ratelimiter [OPTIONS]",
	Short:        "ratelimiter service",
	Long:         "ratelimiter service with internal http server which acts as a proxy to backend services",
	SilenceUsage: true,
	Version:      version.String(),
	RunE: func(cmd *cobra.Command, args []string) error {
		store := data.NewRedis(options)
		if err := store.Connect(); err != nil {
			return err
		}
		return nil
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if url, exists := os.LookupEnv("REDIS_URL"); exists && url != "" {
			options.RedisURL = url
		}
		if port, exists := os.LookupEnv("REDIS_PORT"); exists && port != "" {
			if i, err := strconv.Atoi(port); err != nil {
				options.RedisPort = i
			}
		}
		if pwd, exists := os.LookupEnv("REDIS_PASSWORD"); exists && pwd != "" {
			options.RedisPassword = pwd
		}
		if bk, exists := os.LookupEnv("BACKEND_SERVER"); exists && bk != "" {
			backends = strings.Split(bk, ",")
		}
		return nil
	},
}

func configure() *cobra.Command {
	rootCmd.PersistentFlags().StringVar(&options.RedisURL, "redis-url", "127.0.0.1", "Redis URL")
	rootCmd.PersistentFlags().IntVar(&options.RedisPort, "redis-port", 6379, "Redis port")
	rootCmd.PersistentFlags().StringVar(&options.RedisPassword, "redis-password", "", "Redis password")

	rootCmd.PersistentFlags().StringSliceVar(&backends, "backend-server", []string{}, "List of backend servers to proxy to")

	rootCmd.AddCommand(versionCmd)

	return rootCmd
}

// Run runs the `ratelimiter` root command
func Run() error {
	return configure().Execute()
}

// Main wraps Run and sets the log formatter
func Main() {
	// let's explicitly set stdout
	logrus.SetOutput(os.Stdout)

	// this formatter is the default, but the timestamps output aren't
	// particularly useful, they're relative to the command start
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})

	if err := Run(); err != nil {
		os.Exit(1)
	}
}
