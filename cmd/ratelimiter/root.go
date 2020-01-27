package ratelimiter

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/khos2ow/ratelimiter"
	"github.com/khos2ow/ratelimiter/internal/data"
	"github.com/khos2ow/ratelimiter/internal/server"
	"github.com/khos2ow/ratelimiter/internal/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ruleflag struct {
	Limit    int
	Interval int
	Unit     string
}

var backends []string
var options = &data.Options{}
var rf = &ruleflag{}

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
		timeunit, err := convert(rf.Unit)
		if err != nil {
			return err
		}
		rule := ratelimiter.NewRule(rf.Limit, rf.Interval, timeunit)
		limiter := ratelimiter.NewLimiter(rule, store)
		return server.Start(backends, limiter)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if rl, exists := os.LookupEnv("RATE_LIMIT"); exists && rl != "" {
			if i, err := strconv.Atoi(rl); err == nil {
				rf.Limit = i
			}
		}
		if ri, exists := os.LookupEnv("RATE_INTERVAL"); exists && ri != "" {
			if i, err := strconv.Atoi(ri); err == nil {
				rf.Interval = i
			}
		}
		if ru, exists := os.LookupEnv("RATE_TIMEUNIT"); exists && ru != "" {
			rf.Unit = ru
		}
		if rf.Limit == 0 {
			return fmt.Errorf("invalid value '0' for --rate-limit")
		}
		if rf.Interval == 0 {
			return fmt.Errorf("invalid value '0' for --rate-interval")
		}
		if url, exists := os.LookupEnv("REDIS_URL"); exists && url != "" {
			options.RedisURL = url
		}
		if port, exists := os.LookupEnv("REDIS_PORT"); exists && port != "" {
			if i, err := strconv.Atoi(port); err == nil {
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

func convert(u string) (time.Duration, error) {
	switch u {
	case "ms":
		return time.Millisecond, nil
	case "s":
		return time.Second, nil
	case "m":
		return time.Minute, nil
	case "h":
		return time.Hour, nil
	}
	return 0, fmt.Errorf("unknown value '%s' for --rate-timeunit", u)
}

func configure() *cobra.Command {
	rootCmd.PersistentFlags().IntVar(&rf.Limit, "rate-limit", 0, "Maximum number of hits to allow in every unit of time")
	rootCmd.PersistentFlags().IntVar(&rf.Interval, "rate-interval", 0, "Interval for limiting hits every unit of time in")
	rootCmd.PersistentFlags().StringVar(&rf.Unit, "rate-timeunit", "s", "Unit of time for limiting hits in each interval [ms, s, m, h]")

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
