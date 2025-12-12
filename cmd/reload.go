package cmd

import (
	"fmt"

	"github.com/osmargm1202/orgmcron/internal/service"
	"github.com/spf13/cobra"
)

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Recarga la configuración y reinicia el servicio",
	Long:  "Recarga la configuración de jobs y reinicia el servicio systemd si está corriendo",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Verificar si el servicio existe
		if !service.ServiceExists() {
			fmt.Println("El servicio no está instalado. Ejecuta 'orgmcron install' primero.")
			return nil
		}

		// Verificar si el servicio está corriendo
		if service.IsServiceRunning() {
			fmt.Println("Reiniciando servicio systemd...")
			if err := service.RestartService(); err != nil {
				return fmt.Errorf("error reiniciando servicio: %w", err)
			}
			fmt.Println("✓ Servicio reiniciado exitosamente")
		} else {
			fmt.Println("El servicio no está corriendo. Los cambios se aplicarán cuando se inicie el servicio.")
			fmt.Println("Para iniciar el servicio, ejecuta:")
			fmt.Println("  systemctl --user start orgmcron")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(reloadCmd)
}

