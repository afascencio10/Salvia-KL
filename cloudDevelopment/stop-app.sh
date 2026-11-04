#!/bin/bash
# Detiene la instancia de la app que esté escuchando en el puerto de desarrollo
# (por defecto 9090), para poder volver a correr ./run_dev.sh.
#
# Por qué existe: `go run` compila el backend a un binario llamado `bitsflow` y lo
# lanza como proceso hijo. Si se mata solo `go run` (kill, TaskStop, cerrar la
# terminal), el binario queda HUÉRFANO escuchando en el puerto y sigue sirviendo
# código viejo; un nuevo ./run_dev.sh fallaría con "address already in use" o,
# peor, parecería arrancar mientras responde la instancia anterior.
#
# Solo mata el proceso que realmente escucha en el puerto, y solo si es el
# binario de la app (…/exe/bitsflow): nunca por nombre ni a ciegas.
#
# Uso: cloudDevelopment/stop-app.sh        (también lo llama setup-local-postgres.sh)
#      PORT=9091 cloudDevelopment/stop-app.sh

PORT="${PORT:-9090}"
port_hex="$(printf '%04X' "$PORT")"

# Inodo del socket que escucha en el puerto (columna 4 = 0A → LISTEN).
inode="$(awk -v p=":$port_hex" '$4=="0A" && $2 ~ p"$" {print $10; exit}' /proc/net/tcp 2>/dev/null)"
if [ -z "$inode" ]; then
    echo "[stop-app] Nada escucha en :$PORT."
    exit 0
fi

# PID dueño de ese socket.
pid=""
for fd in /proc/[0-9]*/fd/*; do
    if [ "$(readlink "$fd" 2>/dev/null)" = "socket:[$inode]" ]; then
        pid="${fd#/proc/}"
        pid="${pid%%/*}"
        break
    fi
done
if [ -z "$pid" ]; then
    echo "[stop-app] :$PORT está ocupado pero no encontré el proceso (¿otro usuario?)." >&2
    exit 1
fi

cmd="$(tr '\0' ' ' < "/proc/$pid/cmdline" 2>/dev/null)"
case "$cmd" in
    */exe/bitsflow*)
        echo "[stop-app] Deteniendo la app (PID $pid) que ocupa :$PORT..."
        kill "$pid"
        for _ in $(seq 1 20); do
            [ -e "/proc/$pid" ] || { echo "[stop-app] Detenida."; exit 0; }
            sleep 0.5
        done
        echo "[stop-app] No terminó con SIGTERM; forzando..."
        kill -9 "$pid"
        ;;
    *)
        echo "[stop-app] :$PORT lo ocupa otro proceso (PID $pid: ${cmd:0:80}); no es la app, no lo toco." >&2
        exit 1
        ;;
esac
