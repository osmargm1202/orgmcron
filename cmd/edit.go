package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/osmargm1202/orgmcron/internal/config"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [job_name]",
	Short: "Edita un job existente",
	Long:  "Edita un job existente usando una interfaz interactiva",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		jobName := args[0]

		// Cargar job existente
		existingJob, err := config.GetJobByName(jobName)
		if err != nil {
			return fmt.Errorf("error cargando job: %w", err)
		}

		var (
			scheduleType    string
			cronExpr        string
			intervalExpr    string
			commandsStr     string
			healthcheckName string
		)

		// Determinar tipo de schedule
		if strings.HasPrefix(existingJob.Schedule, "@") {
			scheduleType = "interval"
			intervalExpr = existingJob.Schedule
		} else {
			scheduleType = "cron"
			cronExpr = existingJob.Schedule
		}

		// Preparar comandos
		commandsStr = strings.Join(existingJob.Commands, "\n")

		// Extraer nombre del healthcheck de la URL
		if existingJob.HealthcheckURL != "" {
			parts := strings.Split(existingJob.HealthcheckURL, "/")
			if len(parts) > 0 {
				healthcheckName = parts[len(parts)-1]
			}
		}

		// Primer formulario: tipo de schedule
		form1 := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Tipo de schedule").
					Description("Selecciona el tipo de programación").
					Options(
						huh.NewOption("Intervalo (@every)", "interval"),
						huh.NewOption("Expresión Cron (* * * * *)", "cron"),
					).
					Value(&scheduleType),
			),
		)

		if err := form1.Run(); err != nil {
			return fmt.Errorf("error en el formulario: %w", err)
		}

		// Segundo formulario: schedule específico según el tipo
		var form2 *huh.Form
		if scheduleType == "cron" {
			form2 = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Expresión Cron").
						Description("Formato: minuto hora día mes día-semana (ej: '0 * * * *' para cada hora) o con segundos: segundo minuto hora día mes día-semana").
						Value(&cronExpr).
						Validate(func(s string) error {
							if s == "" {
								return fmt.Errorf("la expresión cron es requerida")
							}
							return nil
						}),
				),
			)
		} else {
			form2 = huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Intervalo").
						Description("Selecciona el intervalo de ejecución").
						Options(
							huh.NewOption("Cada minuto", "@every 1m"),
							huh.NewOption("Cada 5 minutos", "@every 5m"),
							huh.NewOption("Cada 10 minutos", "@every 10m"),
							huh.NewOption("Cada 15 minutos", "@every 15m"),
							huh.NewOption("Cada 30 minutos", "@every 30m"),
							huh.NewOption("Cada hora", "@every 1h"),
							huh.NewOption("Cada 3 horas", "@every 3h"),
							huh.NewOption("Cada 6 horas", "@every 6h"),
							huh.NewOption("Cada 10 horas", "@every 10h"),
							huh.NewOption("Cada 12 horas", "@every 12h"),
							huh.NewOption("Diario", "@daily"),
							huh.NewOption("Semanal", "@weekly"),
						).
						Value(&intervalExpr),
				),
			)
		}

		if err := form2.Run(); err != nil {
			return fmt.Errorf("error en el formulario: %w", err)
		}

		// Tercer formulario: comandos y healthcheck
		form3 := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Comandos").
					Description("Un comando por línea. Se ejecutarán en orden.").
					Value(&commandsStr).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("debe proporcionar al menos un comando")
						}
						return nil
					}).
					CharLimit(10000),

				huh.NewInput().
					Title("Nombre para healthcheck").
					Description("Nombre del proyecto para el healthcheck (se construirá la URL automáticamente)").
					Value(&healthcheckName).
					Placeholder("prueba"),
			),
		)

		if err := form3.Run(); err != nil {
			return fmt.Errorf("error en el formulario: %w", err)
		}

		// Determinar schedule final
		var schedule string
		if scheduleType == "cron" {
			schedule = cronExpr
		} else {
			schedule = intervalExpr
		}

		// Parsear comandos
		commandLines := strings.Split(commandsStr, "\n")
		commands := []string{}
		for _, line := range commandLines {
			line = strings.TrimSpace(line)
			if line != "" {
				commands = append(commands, line)
			}
		}

		// Construir healthcheck URL
		var healthcheckURL string
		if healthcheckName != "" {
			healthcheckURL = fmt.Sprintf("https://hc.or-gm.com/ping/{pingkey}/%s", healthcheckName)
		}

		// Actualizar job
		updatedJob := config.Job{
			Name:           jobName,
			Schedule:       schedule,
			Commands:       commands,
			HealthcheckURL: healthcheckURL,
		}

		if err := config.UpdateJob(jobName, updatedJob); err != nil {
			return fmt.Errorf("error actualizando job: %w", err)
		}

		fmt.Printf("\n✓ Job '%s' actualizado exitosamente\n", jobName)
		fmt.Printf("  Schedule: %s\n", schedule)
		fmt.Printf("  Comandos: %d\n", len(commands))
		if healthcheckURL != "" {
			fmt.Printf("  Healthcheck: %s\n", healthcheckURL)
		}
		fmt.Println("\nPara aplicar los cambios, ejecuta:")
		fmt.Println("  orgmcron reload")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

