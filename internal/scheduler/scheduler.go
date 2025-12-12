package scheduler

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/osmargm1202/orgmcron/internal/config"
	"github.com/osmargm1202/orgmcron/internal/healthcheck"
	"github.com/osmargm1202/orgmcron/internal/job"
	"github.com/osmargm1202/orgmcron/internal/logger"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron      *cron.Cron
	jobs      map[string]cron.EntryID
	mu        sync.RWMutex
	pingKey   string
	stopChan  chan struct{}
	reloadChan chan struct{}
}

// NewScheduler crea un nuevo scheduler
func NewScheduler(pingKey string) *Scheduler {
	// Usar WithSeconds para soportar expresiones con segundos
	// Las expresiones sin segundos también funcionan (se asume 0 segundos)
	c := cron.New(cron.WithSeconds())
	return &Scheduler{
		cron:       c,
		jobs:       make(map[string]cron.EntryID),
		pingKey:    pingKey,
		stopChan:   make(chan struct{}),
		reloadChan: make(chan struct{}),
	}
}

// LoadJobs carga los jobs desde la configuración y los programa
func (s *Scheduler) LoadJobs() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.DebugLog("Iniciando carga de jobs")

	// Detener y limpiar jobs existentes
	s.cron.Stop()
	// Usar WithSeconds para soportar expresiones con segundos
	s.cron = cron.New(cron.WithSeconds())
	s.jobs = make(map[string]cron.EntryID)

	// Cargar configuración
	config, err := config.LoadJobs()
	if err != nil {
		logger.DebugLog("Error cargando jobs: %v", err)
		return fmt.Errorf("error cargando jobs: %w", err)
	}

	logger.DebugLog("Cargados %d jobs desde la configuración", len(config.Jobs))

	// Programar cada job
	for _, j := range config.Jobs {
		if err := s.scheduleJob(j); err != nil {
			logger.DebugLog("Error programando job '%s': %v", j.Name, err)
			fmt.Fprintf(os.Stderr, "Error programando job '%s': %v\n", j.Name, err)
			continue
		}
		logger.DebugLog("Job '%s' programado exitosamente con schedule '%s'", j.Name, j.Schedule)
	}

	// Iniciar el cron
	s.cron.Start()
	logger.DebugLog("Scheduler iniciado con %d jobs programados", len(s.jobs))
	return nil
}

// normalizeSchedule normaliza una expresión cron para que funcione con WithSeconds
// Si es una expresión de 5 campos (sin segundos), agrega "0" al inicio
func normalizeSchedule(schedule string) string {
	// Si ya empieza con @, es un intervalo especial, no necesita normalización
	if strings.HasPrefix(schedule, "@") {
		return schedule
	}
	
	// Contar campos separados por espacios
	fields := strings.Fields(schedule)
	if len(fields) == 5 {
		// Es una expresión cron estándar sin segundos, agregar "0" al inicio
		return "0 " + schedule
	}
	
	// Ya tiene 6 campos o es inválida, retornar tal cual
	return schedule
}

// scheduleJob programa un job individual
func (s *Scheduler) scheduleJob(j config.Job) error {
	// Normalizar el schedule para que funcione con WithSeconds
	normalizedSchedule := normalizeSchedule(j.Schedule)
	
	entryID, err := s.cron.AddFunc(normalizedSchedule, func() {
		logger.DebugLog("Ejecutando job: %s (schedule: %s)", j.Name, j.Schedule)
		fmt.Fprintf(os.Stdout, "[%s] Ejecutando job: %s\n", j.Schedule, j.Name)
		
		exitCode, err := job.Execute(j, s.pingKey)
		if err != nil {
			logger.DebugLog("Error ejecutando job '%s': %v", j.Name, err)
			fmt.Fprintf(os.Stderr, "[%s] Error ejecutando job: %v\n", j.Name, err)
			return
		}

		logger.DebugLog("Job '%s' completado con código de salida: %d", j.Name, exitCode)

		// Solo enviar healthcheck si el job fue exitoso
		if exitCode == 0 && j.HealthcheckURL != "" {
			logger.DebugLog("Enviando healthcheck para job '%s' a URL: %s", j.Name, j.HealthcheckURL)
			if err := healthcheck.SendHealthcheck(j.HealthcheckURL, s.pingKey); err != nil {
				logger.DebugLog("Error enviando healthcheck para job '%s': %v", j.Name, err)
				fmt.Fprintf(os.Stderr, "[%s] Error enviando healthcheck: %v\n", j.Name, err)
			} else {
				logger.DebugLog("Healthcheck enviado exitosamente para job '%s'", j.Name)
				fmt.Fprintf(os.Stdout, "[%s] Healthcheck enviado exitosamente\n", j.Name)
			}
		} else if exitCode != 0 {
			logger.DebugLog("Job '%s' falló con código %d, no se envía healthcheck", j.Name, exitCode)
			fmt.Fprintf(os.Stderr, "[%s] Job falló con código %d, no se envía healthcheck\n", j.Name, exitCode)
		}
	})

	if err != nil {
		return fmt.Errorf("error agregando job al cron: %w", err)
	}

	s.jobs[j.Name] = entryID
	fmt.Fprintf(os.Stdout, "Job '%s' programado con schedule '%s'\n", j.Name, j.Schedule)
	return nil
}

// UpdatePingKey actualiza la pingkey
func (s *Scheduler) UpdatePingKey(pingKey string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pingKey = pingKey
}

// Reload recarga los jobs desde la configuración
func (s *Scheduler) Reload() error {
	return s.LoadJobs()
}

// Stop detiene el scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cron.Stop()
	close(s.stopChan)
}

// Start inicia el scheduler y espera señales
func (s *Scheduler) Start() error {
	logger.DebugLog("Iniciando scheduler con pingkey: %s", s.pingKey)
	// Cargar jobs iniciales
	if err := s.LoadJobs(); err != nil {
		logger.DebugLog("Error cargando jobs iniciales: %v", err)
		return err
	}
	logger.DebugLog("Scheduler iniciado correctamente")

	// Configurar manejo de señales
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Esperar señales o stop
	for {
		select {
		case sig := <-sigChan:
			switch sig {
			case syscall.SIGHUP:
				// Recargar configuración
				logger.DebugLog("Recibida señal SIGHUP, recargando configuración")
				fmt.Fprintf(os.Stdout, "Recibida señal SIGHUP, recargando configuración...\n")
				if err := s.Reload(); err != nil {
					logger.DebugLog("Error recargando configuración: %v", err)
					fmt.Fprintf(os.Stderr, "Error recargando configuración: %v\n", err)
				} else {
					logger.DebugLog("Configuración recargada exitosamente")
				}
			case syscall.SIGINT, syscall.SIGTERM:
				// Detener scheduler
				logger.DebugLog("Recibida señal %v, deteniendo scheduler", sig)
				fmt.Fprintf(os.Stdout, "Recibida señal %v, deteniendo scheduler...\n", sig)
				s.Stop()
				logger.DebugLog("Scheduler detenido")
				return nil
			}
		case <-s.stopChan:
			return nil
		case <-s.reloadChan:
			logger.DebugLog("Recarga manual solicitada")
			if err := s.Reload(); err != nil {
				logger.DebugLog("Error recargando configuración: %v", err)
				fmt.Fprintf(os.Stderr, "Error recargando configuración: %v\n", err)
			} else {
				logger.DebugLog("Configuración recargada exitosamente (manual)")
			}
		}
	}
}

// TriggerReload dispara una recarga manual
func (s *Scheduler) TriggerReload() {
	select {
	case s.reloadChan <- struct{}{}:
	default:
	}
}

