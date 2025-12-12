package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	DebugLogFile = "debug.log"
)

// GetDebugLogPath retorna la ruta del archivo de log de depuración
func GetDebugLogPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error obteniendo directorio home: %w", err)
	}
	logsDir := filepath.Join(home, ".config", "orgmcron", "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return "", fmt.Errorf("error creando directorio de logs: %w", err)
	}
	return filepath.Join(logsDir, DebugLogFile), nil
}

// DebugLog escribe un mensaje de depuración al archivo de log
func DebugLog(format string, args ...interface{}) error {
	logPath, err := GetDebugLogPath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error abriendo archivo de log de depuración: %w", err)
	}
	defer file.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] %s\n", timestamp, message)
	
	if _, err := file.WriteString(logLine); err != nil {
		return fmt.Errorf("error escribiendo log de depuración: %w", err)
	}

	return nil
}

// DebugLogAndPrint escribe un mensaje de depuración al archivo y lo imprime en consola si debug=true
func DebugLogAndPrint(debug bool, format string, args ...interface{}) error {
	err := DebugLog(format, args...)
	if err != nil {
		return err
	}
	
	if debug {
		fmt.Printf("[DEBUG] %s\n", fmt.Sprintf(format, args...))
	}
	
	return nil
}


