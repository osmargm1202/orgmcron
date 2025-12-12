package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ConfigDir  = ".config/orgmcron"
	JobsFile   = "jobs.json"
	ConfigFile = "config.json"
	LogsDir    = "logs"
)

type Job struct {
	Name           string   `json:"name"`
	Schedule       string   `json:"schedule"`
	Commands       []string `json:"commands"`
	HealthcheckURL string   `json:"healthcheck_url"`
}

type JobsConfig struct {
	Jobs []Job `json:"jobs"`
}

type AppConfig struct {
	PingKey string `json:"pingkey"`
}

// GetConfigDir retorna el directorio de configuración completo
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error obteniendo directorio home: %w", err)
	}
	return filepath.Join(home, ConfigDir), nil
}

// GetLogsDir retorna el directorio de logs completo
func GetLogsDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, LogsDir), nil
}

// EnsureConfigDir crea el directorio de configuración si no existe
func EnsureConfigDir() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de configuración: %w", err)
	}
	return nil
}

// EnsureLogsDir crea el directorio de logs si no existe
func EnsureLogsDir() error {
	logsDir, err := GetLogsDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de logs: %w", err)
	}
	return nil
}

// LoadJobs carga los jobs desde el archivo jobs.json
func LoadJobs() (*JobsConfig, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	jobsPath := filepath.Join(configDir, JobsFile)
	data, err := os.ReadFile(jobsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Si no existe, retornar configuración vacía
			return &JobsConfig{Jobs: []Job{}}, nil
		}
		return nil, fmt.Errorf("error leyendo jobs.json: %w", err)
	}

	var config JobsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parseando jobs.json: %w", err)
	}

	return &config, nil
}

// SaveJobs guarda los jobs en el archivo jobs.json
func SaveJobs(config *JobsConfig) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	jobsPath := filepath.Join(configDir, JobsFile)
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando jobs: %w", err)
	}

	if err := os.WriteFile(jobsPath, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo jobs.json: %w", err)
	}

	return nil
}

// LoadConfig carga la configuración de la aplicación
func LoadConfig() (*AppConfig, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, ConfigFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Si no existe, retornar configuración vacía
			return &AppConfig{}, nil
		}
		return nil, fmt.Errorf("error leyendo config.json: %w", err)
	}

	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parseando config.json: %w", err)
	}

	return &config, nil
}

// SaveConfig guarda la configuración de la aplicación
func SaveConfig(config *AppConfig) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, ConfigFile)
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando configuración: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo config.json: %w", err)
	}

	return nil
}

// GetJobByName busca un job por nombre
func GetJobByName(name string) (*Job, error) {
	config, err := LoadJobs()
	if err != nil {
		return nil, err
	}

	for i := range config.Jobs {
		if config.Jobs[i].Name == name {
			return &config.Jobs[i], nil
		}
	}

	return nil, fmt.Errorf("job '%s' no encontrado", name)
}

// AddJob agrega un nuevo job
func AddJob(job Job) error {
	config, err := LoadJobs()
	if err != nil {
		return err
	}

	// Verificar que no exista un job con el mismo nombre
	for _, j := range config.Jobs {
		if j.Name == job.Name {
			return fmt.Errorf("ya existe un job con el nombre '%s'", job.Name)
		}
	}

	config.Jobs = append(config.Jobs, job)
	return SaveJobs(config)
}

// UpdateJob actualiza un job existente
func UpdateJob(name string, job Job) error {
	config, err := LoadJobs()
	if err != nil {
		return err
	}

	found := false
	for i := range config.Jobs {
		if config.Jobs[i].Name == name {
			config.Jobs[i] = job
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("job '%s' no encontrado", name)
	}

	return SaveJobs(config)
}

// DeleteJob elimina un job
func DeleteJob(name string) error {
	config, err := LoadJobs()
	if err != nil {
		return err
	}

	found := false
	for i, j := range config.Jobs {
		if j.Name == name {
			config.Jobs = append(config.Jobs[:i], config.Jobs[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("job '%s' no encontrado", name)
	}

	return SaveJobs(config)
}

