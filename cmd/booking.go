/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var b *string

// bookingCmd represents the booking command
var bookingCmd = &cobra.Command{
	Use:   "booking",
	Short: "handle the bookings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("booking called")
	},
}

func init() {
	rootCmd.AddCommand(bookingCmd)

	b = agendaCmd.Flags().StringP("book", "b", "", "book an agenda")
}
