package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/osmargm1202/orgmcron/internal/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Agrega un nuevo job",
	Long:  "Agrega un nuevo job usando una interfaz interactiva",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			jobName         string
			scheduleType    string
			cronExpr        string
			intervalExpr    string
			commandsStr     string
			healthcheckName string
		)

		// Primer formulario: nombre y tipo de schedule
		form1 := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Nombre del job").
					Description("Nombre único para identificar este job").
					Value(&jobName).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("el nombre no puede estar vacío")
						}
						// Verificar que no exista
						_, err := config.GetJobByName(s)
						if err == nil {
							return fmt.Errorf("ya existe un job con este nombre")
						}
						return nil
					}),

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

		// Crear job
		job := config.Job{
			Name:           jobName,
			Schedule:       schedule,
			Commands:       commands,
			HealthcheckURL: healthcheckURL,
		}

		// Guardar job
		if err := config.AddJob(job); err != nil {
			return fmt.Errorf("error guardando job: %w", err)
		}

		fmt.Printf("\n✓ Job '%s' agregado exitosamente\n", jobName)
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
	rootCmd.AddCommand(addCmd)
}

