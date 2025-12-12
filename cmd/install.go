package cmd

import (
	"fmt"

	"github.com/osmargm1202/orgmcron/internal/service"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Instala el servicio systemd --user",
	Long:  "Crea el archivo de servicio systemd y lo habilita para ejecutarse automáticamente",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Instalando servicio systemd...")

		// Crear archivo de servicio
		if err := service.CreateService(); err != nil {
			return fmt.Errorf("error creando servicio: %w", err)
		}
		fmt.Println("✓ Archivo de servicio creado")

		// Recargar daemon
		if err := service.ReloadDaemon(); err != nil {
			return fmt.Errorf("error recargando daemon: %w", err)
		}
		fmt.Println("✓ Daemon recargado")

		// Habilitar servicio
		if err := service.EnableService(); err != nil {
			return fmt.Errorf("error habilitando servicio: %w", err)
		}
		fmt.Println("✓ Servicio habilitado")

		fmt.Println("\nServicio instalado exitosamente.")
		fmt.Println("Para iniciar el servicio, ejecuta:")
		fmt.Println("  systemctl --user start orgmcron")
		fmt.Println("\nPara ver el estado del servicio:")
		fmt.Println("  systemctl --user status orgmcron")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

