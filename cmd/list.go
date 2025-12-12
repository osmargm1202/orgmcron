package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/osmargm1202/orgmcron/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos los jobs configurados",
	Long:  "Muestra una lista de todos los jobs configurados con sus detalles",
	RunE: func(cmd *cobra.Command, args []string) error {
		jobsConfig, err := config.LoadJobs()
		if err != nil {
			return fmt.Errorf("error cargando jobs: %w", err)
		}

		if len(jobsConfig.Jobs) == 0 {
			fmt.Println("No hay jobs configurados.")
			fmt.Println("Usa 'orgmcron add' para agregar un nuevo job.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NOMBRE\tSCHEDULE\tCOMANDOS\tHEALTHCHECK")
		fmt.Fprintln(w, "------\t--------\t--------\t-----------")

		for _, job := range jobsConfig.Jobs {
			commandsCount := fmt.Sprintf("%d", len(job.Commands))
			healthcheck := "No"
			if job.HealthcheckURL != "" {
				healthcheck = "Sí"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", job.Name, job.Schedule, commandsCount, healthcheck)
		}

		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

