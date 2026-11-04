#!/usr/bin/env bash
# Reinicia la app local (necesario tras cambios en Go o en el frontend embebido).
cd /home/user/Salvia-KL
bash cloudDevelopment/stop-app.sh >/dev/null 2>&1
(nohup ./run_dev.sh > /tmp/run_dev.log 2>&1 &)
for i in $(seq 1 60); do
  sleep 2
  if grep -q "Servidor corriendo en :9090" /tmp/run_dev.log 2>/dev/null; then
    echo "app arriba tras $((i*2))s"; exit 0
  fi
done
echo "la app no arrancó"; tail -20 /tmp/run_dev.log; exit 1
