package cmd

import (
	"Z3NTL3/ProxyBeast/proxy"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use: "proxybeast",
	Short: "The ultimate proxy checker",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(time.Duration(proxy.Timeout).Seconds())
		return nil
	},
}

func Execute() error {
	flags := map[string]map[string]any{
		"timeout": {
			"ptr": &proxy.Timeout,
			"short": "t",
			"usage": "Maximum proxy latency/timeout",
		},
		"file": {
			"ptr": &proxy.ProxyFilePath,
			"short": "f",
			"usage": "File to input proxies, aka proxy list",
		},
	}

	for k, v := range flags {
		rootCmd.PersistentFlags().VarP(
					v["ptr"].(pflag.Value),
					k,
					v["short"].(string),
					v["usage"].(string),
				)
	}	

	return rootCmd.Execute()
}