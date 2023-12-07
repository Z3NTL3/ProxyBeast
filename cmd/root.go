package cmd

import (
	"Z3NTL3/proxy-checker/globals"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {},
}

func Init() {
	RootCmd.Flags().IntVar(
		&globals.Timeout, "timeout", 5,
		"Sets custom timeout in seconds",
	)
	RootCmd.MarkFlagRequired("timeout")

	RootCmd.Flags().StringVar(
		&globals.Protocol, "protocol",
		"http", "The proxy protocol to check against",
	)
	RootCmd.MarkFlagRequired("protocol")

	RootCmd.Flags().StringVar(
		&globals.ProxyFile, "file", "proxies.txt",
		"Determines your proxy file name requires to be *.txt matching",
	)

	RootCmd.Flags().IntVar(
		&globals.Retries, "retry", 2,
		"The amount of tries to retry to connect to a failure proxy",
	)

	RootCmd.Flags().BoolVar(
		&globals.Multi, "multi",
		false, "If passed as arg, it will check for all protocols, will tear down the accuracy",
	)
}
