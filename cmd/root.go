/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"alinea.com/cmd/server"
	"alinea.com/internal/agenda"
	"alinea.com/internal/booking"
	"alinea.com/pkg/mongo"
	"alinea.com/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "scheduler domain as a cli application",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var agendaService *agenda.AgendaService
var bookingService *booking.BookingService

func init() {
	rootCmd.AddCommand(server.ServerCmd)

	client := utils.Must(mongo.NewClient(context.Background()))

	bookingRepository := mongo.NewBookingRepository(client)
	agendaRepository := mongo.NewAgendaRepository(client)
	blockRepository := mongo.NewBlockRepository(client)
	serviceRepository := mongo.NewServiceRepository(client)

	bookingService = booking.NewBookingService(bookingRepository, agendaRepository, blockRepository, serviceRepository, createEventPublisher())
	agendaService = agenda.NewAgendaService(agendaRepository, serviceRepository)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scheduler.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".scheduler" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".scheduler")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
