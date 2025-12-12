package healthcheck

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/osmargm1202/orgmcron/internal/logger"
)

// SendHealthcheck envía un request GET al healthcheck URL
func SendHealthcheck(url string, pingKey string) error {
	// Reemplazar {pingkey} en la URL
	finalURL := strings.ReplaceAll(url, "{pingkey}", pingKey)
	logger.DebugLog("Enviando healthcheck a URL: %s", finalURL)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(finalURL)
	if err != nil {
		logger.DebugLog("Error enviando healthcheck a %s: %v", finalURL, err)
		return fmt.Errorf("error enviando healthcheck: %w", err)
	}
	defer resp.Body.Close()

	logger.DebugLog("Healthcheck respondió con código de estado: %d", resp.StatusCode)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.DebugLog("Healthcheck falló con código de estado: %d", resp.StatusCode)
		return fmt.Errorf("healthcheck retornó código de estado: %d", resp.StatusCode)
	}

	logger.DebugLog("Healthcheck enviado exitosamente")
	return nil
}

