package cmd

import (
	"fmt"

	"github.com/osmargm1202/orgmcron/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Gestiona la configuración",
	Long:  "Gestiona la configuración de orgmcron",
}

var pingkeyCmd = &cobra.Command{
	Use:   "pingkey [key]",
	Short: "Configura o muestra la pingkey",
	Long:  "Configura la pingkey para los healthchecks. Si no se proporciona key, muestra la actual.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		appConfig, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("error cargando configuración: %w", err)
		}

		if len(args) == 0 {
			// Mostrar pingkey actual
			if appConfig.PingKey == "" {
				fmt.Println("PingKey no configurado")
			} else {
				fmt.Printf("PingKey actual: %s\n", appConfig.PingKey)
			}
			return nil
		}

		// Configurar nueva pingkey
		appConfig.PingKey = args[0]
		if err := config.SaveConfig(appConfig); err != nil {
			return fmt.Errorf("error guardando configuración: %w", err)
		}

		fmt.Printf("PingKey configurado: %s\n", appConfig.PingKey)
		return nil
	},
}

func init() {
	configCmd.AddCommand(pingkeyCmd)
	rootCmd.AddCommand(configCmd)
}

