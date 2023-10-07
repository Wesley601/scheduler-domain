/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"alinea.com/internal/schedule"
	"alinea.com/pkg/mongo"
	"alinea.com/pkg/utils"
	"github.com/spf13/cobra"
)

var c *string
var g *string
var parser *schedule.Parser

// agendaCmd represents the agenda command
var agendaCmd = &cobra.Command{
	Use:   "agenda",
	Short: "handle the providers agendas",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if *c != "" {
			err := parser.FromJSON([]byte(*c))
			if err != nil {
				panic(err)
			}

			a, err := parser.ToAgenda()
			if err != nil {
				panic(err)
			}

			agendaService.Create(schedule.CreateScheduleDTO{
				Name: a.Name,
				Slots: func() []schedule.CreateSlotDTO {
					var slots []schedule.CreateSlotDTO

					for _, slot := range a.Slots {
						slots = append(slots, schedule.CreateSlotDTO(slot))
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

var agendaService *schedule.ScheduleService

func init() {
	rootCmd.AddCommand(agendaCmd)

	parser = &schedule.Parser{}

	agendaRepository := mongo.NewScheduleRepository(utils.Must(mongo.NewClient()))
	agendaService = schedule.NewScheduleService(agendaRepository)

	c = agendaCmd.Flags().StringP("create", "c", "", "create a new agenda")
	g = agendaCmd.Flags().StringP("get", "g", "", "get a agenda by id")
}
