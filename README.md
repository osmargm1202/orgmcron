# orgmcron

CLI en Go para **ejecutar cronjobs como daemon de usuario** (systemd `--user`) usando **robfig/cron**, ejecutando **comandos arbitrarios** (ej. `rsync`) y notificando un **healthcheck URL** cuando el job termina correctamente.

## Instalación

### Opción A: Script (recomendado)

```bash
curl -fsSL custom.or-gm.com/orgmcron.sh | sh
```

Instala el binario en `~/.local/bin/orgmcron`.

### Opción B: Compilar local

```bash
go build -o orgmcron .
```

## Configuración

Los archivos viven en:
- **Config**: `~/.config/orgmcron/config.json`
- **Jobs**: `~/.config/orgmcron/jobs.json`
- **Logs**: `~/.config/orgmcron/logs/`
  - `~/.config/orgmcron/logs/<job>.log`
  - `~/.config/orgmcron/logs/debug.log`

### Configurar pingkey (healthchecks)

```bash
orgmcron config pingkey zj46yb44fqw2bmlyt2bdgg
```

La URL de healthcheck soporta `{pingkey}` y se reemplaza automáticamente.

Ejemplo (de prueba):
`https://hc.or-gm.com/ping/zj46yb44fqw2bmlyt2bdgg/prueba`

## Uso (comandos)

### Crear un job (interactivo)

```bash
orgmcron add
```

Te pedirá:
- nombre del job
- tipo de schedule (`@every ...` o expresión cron)
- comandos (1 por línea, se ejecutan en orden)
- nombre para healthcheck (construye `https://hc.or-gm.com/ping/{pingkey}/<nombre>`)

### Editar un job existente (interactivo)

```bash
orgmcron edit <job_name>
```

### Listar jobs

```bash
orgmcron list
```

### Ver logs de un job (en tiempo real)

```bash
orgmcron log <job_name>
```

### Ejecutar el daemon (foreground)

```bash
orgmcron start
```

### Aplicar cambios de configuración (manual)

Cuando actualizas `jobs.json`, recarga el servicio:

```bash
orgmcron reload
```

## Schedules soportados

- **Intervalos**: `@every 1m`, `@every 1h`, `@daily`, `@weekly`, etc.
- **Cron estándar**: `* * * * *` (5 campos)
- **Cron con segundos**: `* * * * * *` (6 campos)

Nota: internamente el scheduler usa `WithSeconds()`. Si configuras un cron de **5 campos**, se normaliza automáticamente agregando `0` segundos al inicio.

## Servicio systemd (usuario)

Instala el servicio:

```bash
orgmcron install
systemctl --user daemon-reload
systemctl --user enable orgmcron
systemctl --user start orgmcron
```

Ver estado:

```bash
systemctl --user status orgmcron
```

## Formato de jobs.json (referencia)

```json
{
  "jobs": [
    {
      "name": "prueba",
      "schedule": "@every 1h",
      "commands": ["rsync -avz /origen /destino"],
      "healthcheck_url": "https://hc.or-gm.com/ping/{pingkey}/prueba"
    }
  ]
}
```

## Comportamiento del healthcheck

- Se ejecutan los comandos en orden y se guardan en `~/.config/orgmcron/logs/<job>.log`
- **Solo si el job termina con código 0** y `healthcheck_url` no está vacío, se envía un GET al healthcheck.
- Si el job falla, **solo se registra** (no se envía healthcheck).



