package job

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/osmargm1202/orgmcron/internal/config"
	"github.com/osmargm1202/orgmcron/internal/logger"
)

// Execute ejecuta un job y retorna el código de salida
func Execute(job config.Job, pingKey string) (int, error) {
	logger.DebugLog("Iniciando ejecución del job: %s", job.Name)
	
	logsDir, err := config.GetLogsDir()
	if err != nil {
		logger.DebugLog("Error obteniendo directorio de logs para job '%s': %v", job.Name, err)
		return 1, fmt.Errorf("error obteniendo directorio de logs: %w", err)
	}

	if err := config.EnsureLogsDir(); err != nil {
		logger.DebugLog("Error creando directorio de logs para job '%s': %v", job.Name, err)
		return 1, err
	}

	logFile := filepath.Join(logsDir, job.Name+".log")
	logger.DebugLog("Escribiendo logs del job '%s' a: %s", job.Name, logFile)
	
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.DebugLog("Error abriendo archivo de log para job '%s': %v", job.Name, err)
		return 1, fmt.Errorf("error abriendo archivo de log: %w", err)
	}
	defer file.Close()

	// Escribir timestamp de inicio
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	file.WriteString(fmt.Sprintf("\n=== Ejecución iniciada: %s ===\n", timestamp))
	logger.DebugLog("Job '%s': ejecutando %d comandos", job.Name, len(job.Commands))

	// Ejecutar comandos en orden
	var lastExitCode int
	for i, cmdStr := range job.Commands {
		logger.DebugLog("Job '%s': ejecutando comando %d/%d: %s", job.Name, i+1, len(job.Commands), cmdStr)
		file.WriteString(fmt.Sprintf("\n[Comando %d/%d] %s\n", i+1, len(job.Commands), cmdStr))

		cmd := exec.Command("sh", "-c", cmdStr)
		cmd.Stdout = file
		cmd.Stderr = file

		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				lastExitCode = exitError.ExitCode()
				logger.DebugLog("Job '%s': comando %d falló con código de salida: %d", job.Name, i+1, lastExitCode)
				file.WriteString(fmt.Sprintf("\n[ERROR] Comando falló con código de salida: %d\n", lastExitCode))
			} else {
				lastExitCode = 1
				logger.DebugLog("Job '%s': error ejecutando comando %d: %v", job.Name, i+1, err)
				file.WriteString(fmt.Sprintf("\n[ERROR] Error ejecutando comando: %v\n", err))
			}
			// Continuar con el siguiente comando aunque este haya fallado
		} else {
			lastExitCode = 0
			logger.DebugLog("Job '%s': comando %d completado exitosamente", job.Name, i+1)
			file.WriteString(fmt.Sprintf("\n[OK] Comando completado exitosamente\n"))
		}
	}

	// Escribir timestamp de fin
	timestamp = time.Now().Format("2006-01-02 15:04:05")
	file.WriteString(fmt.Sprintf("\n=== Ejecución finalizada: %s (código: %d) ===\n\n", timestamp, lastExitCode))
	logger.DebugLog("Job '%s': ejecución finalizada con código de salida: %d", job.Name, lastExitCode)

	return lastExitCode, nil
}

// GetLogPath retorna la ruta del archivo de log de un job
func GetLogPath(jobName string) (string, error) {
	logsDir, err := config.GetLogsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(logsDir, jobName+".log"), nil
}

