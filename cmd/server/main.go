/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package server

import (
	"net/http"

	"github.com/spf13/cobra"
)

var p *string

// serverCmd represents the server command
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "WebServer",
	Long:  `WebServer based on HyperMedia Driven Build in chi and HTMX.`,
	Run: func(cmd *cobra.Command, args []string) {
		http.ListenAndServe(*p, r)
	},
}

func init() {
	p = ServerCmd.Flags().StringP("port", "p", ":3000", "server port")
}
