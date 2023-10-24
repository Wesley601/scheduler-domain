/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"alinea.com/internal/app"
	"alinea.com/internal/core"
	"alinea.com/pkg/utils"
	"github.com/spf13/cobra"
)

// slotsCmd represents the slots command
var slotsCmd = &cobra.Command{
	Use:   "slots",
	Short: "handle the a agenda slots",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if *i == "" {
			fmt.Println("i is required!")
			os.Exit(1)
		}

		if *f == "" || *t == "" {
			fmt.Println("f and t is required!")
			os.Exit(1)
		}

		s, err := app.AgendaService.ListSlots(context.Background(), *i, *sID, utils.Must(core.NewWindowFromString(*f, *t)))
		if err != nil {
			panic(err)
		}

		fmt.Println(s.ToJSON())
	},
}

var f *string
var t *string
var i *string
var sID *string

func init() {
	agendaCmd.AddCommand(slotsCmd)

	f = slotsCmd.Flags().StringP("from", "f", "", "when the window starts")
	t = slotsCmd.Flags().StringP("to", "t", "", "when the window ends")
	i = slotsCmd.Flags().StringP("id", "i", "", "agenda id")
	sID = slotsCmd.Flags().StringP("serviceId", "s", "", "service id")
}
