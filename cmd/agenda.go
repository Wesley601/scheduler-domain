/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	mongodb "go.mongodb.org/mongo-driver/mongo"

	"alinea.com/internal/agenda"
	"github.com/spf13/cobra"
)

var c *string
var g *string
var l *bool

// agendaCmd represents the agenda command
var agendaCmd = &cobra.Command{
	Use:   "agenda",
	Short: "handle the providers agendas",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var parser agenda.Parser

		if *c != "" {
			parser, err := parser.FromJSON([]byte(*c))
			if err != nil {
				panic(err)
			}

			a, err := parser.ToJSONStruct()
			if err != nil {
				panic(err)
			}

			agendaService.Create(context.Background(), agenda.CreateAgendaDTO{
				Name: a.Name,
				Slots: func() []agenda.CreateSlotDTO {
					var slots []agenda.CreateSlotDTO

					for _, slot := range a.Slots {
						slots = append(slots, agenda.CreateSlotDTO(slot))
					}

					return slots
				}(),
			})

			fmt.Println("agenda created")
		}

		if *g != "" {
			s, err := agendaService.FindByID(context.Background(), *g)

			if err == mongodb.ErrNoDocuments {
				fmt.Println("agenda not found")
				return
			}

			if err != nil {
				panic(err)
			}

			j, err := s.ToJSON()
			if err != nil {
				panic(err)
			}

			fmt.Println(string(j))
		}

		if *l {
			s, err := agendaService.List(context.Background())
			if err != nil {
				panic(err)
			}

			j, err := s.ToJSON()
			if err != nil {
				panic(err)
			}

			fmt.Println(string(j))
		}
	},
}

func init() {
	rootCmd.AddCommand(agendaCmd)

	c = agendaCmd.Flags().StringP("create", "c", "", "create a new agenda")
	g = agendaCmd.Flags().StringP("get", "g", "", "get a agenda by id")
	l = agendaCmd.Flags().BoolP("list", "l", false, "list all agendas")
	agendaCmd.Flags().Lookup("list").NoOptDefVal = "true"
}
