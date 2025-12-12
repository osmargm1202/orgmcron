package cmd

import (
	"fmt"
	"os"

	"github.com/osmargm1202/orgmcron/internal/config"
	"github.com/osmargm1202/orgmcron/internal/scheduler"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia el daemon que ejecuta los jobs",
	Long:  "Inicia el daemon que lee la configuración y ejecuta los jobs programados",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Cargar configuración para obtener pingkey
		appConfig, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("error cargando configuración: %w", err)
		}

		if appConfig.PingKey == "" {
			fmt.Fprintf(os.Stderr, "Advertencia: pingkey no configurado. Usa 'orgmcron config pingkey <key>' para configurarlo.\n")
		}

		// Crear scheduler
		sched := scheduler.NewScheduler(appConfig.PingKey)
		
		fmt.Println("Iniciando daemon orgmcron...")
		fmt.Printf("PingKey configurado: %s\n", appConfig.PingKey)
		fmt.Println("Presiona Ctrl+C para detener el daemon")

		// Iniciar scheduler (bloquea hasta recibir señal de parada)
		if err := sched.Start(); err != nil {
			return fmt.Errorf("error iniciando scheduler: %w", err)
		}

		fmt.Println("Daemon detenido.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

