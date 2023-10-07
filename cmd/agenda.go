/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"alinea.com/internal/agenda"
	"alinea.com/pkg/mongo"
	"alinea.com/pkg/utils"
	"github.com/spf13/cobra"
)

var c *string
var g *string

// agendaCmd represents the agenda command
var agendaCmd = &cobra.Command{
	Use:   "agenda",
	Short: "handle the providers agendas",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if *c != "" {
			parser, err := agenda.FromJSON([]byte(*c))
			if err != nil {
				panic(err)
			}

			a, err := parser.ToJSONStruct()
			if err != nil {
				panic(err)
			}

			agendaService.Create(agenda.CreateAgendaDTO{
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
			s, err := agendaService.FindById(*g)
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

var agendaService *agenda.AgendaService

func init() {
	rootCmd.AddCommand(agendaCmd)

	agendaRepository := mongo.NewAgendaRepository(utils.Must(mongo.NewClient()))
	agendaService = agenda.NewAgendaService(agendaRepository)

	c = agendaCmd.Flags().StringP("create", "c", "", "create a new agenda")
	g = agendaCmd.Flags().StringP("get", "g", "", "get a agenda by id")
}
