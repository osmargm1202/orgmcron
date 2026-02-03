# Resumen de uso - orgmcron

CLI para ejecutar cronjobs como daemon de usuario con notificaciones de healthcheck.

## Instalación rápida

```bash
curl -fsSL custom.or-gm.com/orgmcron.sh | sh
```

El binario se instala en `~/.local/bin/orgmcron`.

## Configuración inicial

1. **Configurar pingkey** (para healthchecks):
```bash
orgmcron config pingkey <tu_pingkey>
```

2. **Instalar servicio systemd**:
```bash
orgmcron install
systemctl --user daemon-reload
systemctl --user enable orgmcron
systemctl --user start orgmcron
```

## Comandos principales

### Gestión de jobs

| Comando | Descripción |
|---------|-------------|
| `orgmcron add` | Crear un nuevo job (interactivo) |
| `orgmcron edit <nombre>` | Editar un job existente |
| `orgmcron list` | Listar todos los jobs |
| `orgmcron reload` | Recargar configuración después de cambios manuales |

### Monitoreo

| Comando | Descripción |
|---------|-------------|
| `orgmcron log <nombre>` | Ver logs de un job en tiempo real |
| `systemctl --user status orgmcron` | Ver estado del servicio |

## Tipos de schedule

- **Intervalos**: `@every 1m`, `@every 1h`, `@daily`, `@weekly`
- **Cron estándar**: `* * * * *` (5 campos)
- **Cron con segundos**: `* * * * * *` (6 campos)

## Ubicación de archivos

- **Configuración**: `~/.config/orgmcron/config.json`
- **Jobs**: `~/.config/orgmcron/jobs.json`
- **Logs**: `~/.config/orgmcron/logs/`
  - `~/.config/orgmcron/logs/<job>.log`
  - `~/.config/orgmcron/logs/debug.log`

## Ejemplo de uso

```bash
# 1. Configurar pingkey
orgmcron config pingkey zj46yb44fqw2bmlyt2bdgg

# 2. Crear un job
orgmcron add
# Nombre: backup_diario
# Schedule: @daily
# Comandos: rsync -avz /origen /destino
# Healthcheck: backup_diario

# 3. Ver jobs configurados
orgmcron list

# 4. Ver logs
orgmcron log backup_diario
```

## Notas importantes

- El healthcheck solo se envía si el job termina con código 0
- Los comandos se ejecutan en orden secuencial
- Después de editar `jobs.json` manualmente, ejecutar `orgmcron reload`
- El servicio corre como daemon de usuario (no requiere root)
