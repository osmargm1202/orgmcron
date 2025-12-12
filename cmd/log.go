package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/osmargm1202/orgmcron/internal/job"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log [job_name]",
	Short: "Muestra los logs de un job",
	Long:  "Muestra los logs de un job. Usa tail -f para seguimiento en tiempo real.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jobName := args[0]

		logPath, err := job.GetLogPath(jobName)
		if err != nil {
			return fmt.Errorf("error obteniendo ruta de log: %w", err)
		}

		// Verificar si el archivo existe
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			return fmt.Errorf("no se encontraron logs para el job '%s'", jobName)
		}

		// Ejecutar tail -f
		tailCmd := exec.Command("tail", "-f", logPath)
		tailCmd.Stdout = os.Stdout
		tailCmd.Stderr = os.Stderr

		if err := tailCmd.Run(); err != nil {
			return fmt.Errorf("error ejecutando tail: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}

