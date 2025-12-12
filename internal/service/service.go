package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	ServiceName = "orgmcron.service"
	ServiceDir  = ".config/systemd/user"
)

// GetServicePath retorna la ruta completa del archivo de servicio
func GetServicePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error obteniendo directorio home: %w", err)
	}
	return filepath.Join(home, ServiceDir, ServiceName), nil
}

// GetBinaryPath retorna la ruta del binario
func GetBinaryPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error obteniendo ruta del ejecutable: %w", err)
	}
	// Resolver enlaces simbólicos
	realPath, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		return execPath, nil
	}
	return realPath, nil
}

// CreateService crea el archivo de servicio systemd
func CreateService() error {
	servicePath, err := GetServicePath()
	if err != nil {
		return err
	}

	// Crear directorio si no existe
	serviceDir := filepath.Dir(servicePath)
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de servicios: %w", err)
	}

	binaryPath, err := GetBinaryPath()
	if err != nil {
		return err
	}

	// Para servicios --user, no se especifica User en el archivo de servicio
	serviceContent := fmt.Sprintf(`[Unit]
Description=orgmcron - Gestor de cronjobs con healthchecks
After=network.target

[Service]
Type=simple
ExecStart=%s start
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
`, binaryPath)

	if err := os.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo de servicio: %w", err)
	}

	return nil
}

// EnableService habilita el servicio systemd
func EnableService() error {
	cmd := exec.Command("systemctl", "--user", "enable", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error habilitando servicio: %w", err)
	}
	return nil
}

// StartService inicia el servicio systemd
func StartService() error {
	cmd := exec.Command("systemctl", "--user", "start", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error iniciando servicio: %w", err)
	}
	return nil
}

// StopService detiene el servicio systemd
func StopService() error {
	cmd := exec.Command("systemctl", "--user", "stop", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error deteniendo servicio: %w", err)
	}
	return nil
}

// RestartService reinicia el servicio systemd
func RestartService() error {
	cmd := exec.Command("systemctl", "--user", "restart", ServiceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error reiniciando servicio: %w", err)
	}
	return nil
}

// ReloadDaemon recarga la configuración de systemd
func ReloadDaemon() error {
	cmd := exec.Command("systemctl", "--user", "daemon-reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error recargando daemon: %w", err)
	}
	return nil
}

// ServiceExists verifica si el servicio existe
func ServiceExists() bool {
	servicePath, err := GetServicePath()
	if err != nil {
		return false
	}
	_, err = os.Stat(servicePath)
	return err == nil
}

// IsServiceRunning verifica si el servicio está corriendo
func IsServiceRunning() bool {
	cmd := exec.Command("systemctl", "--user", "is-active", "--quiet", ServiceName)
	return cmd.Run() == nil
}

