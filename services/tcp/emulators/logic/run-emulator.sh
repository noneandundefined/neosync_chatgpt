#!/usr/bin/env bash
# NeoSync device emulator (Go)
#
# Usage:
#   bash tcp-server/emulators/run-emulator.sh
#   bash tcp-server/emulators/run-emulator.sh -- -host 127.0.0.1 -reconnect-sec 1800 -sync-sec 300
#
# State files:
#   ./emulator-state/state.json   — imei, hash, last_mod, firmware
#   ./emulator-state/config.bin   — конфиг (обновляется при SET_CFG с сервера)
#   ./emulator-state/config.hex   — hex-дамп config.bin
#
# Import config once:
#   bash tcp-server/emulators/run-emulator.sh -- -import-config-hex tcp-server/emulators/default-adm333.cfg.hex

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

go run ./logic "$@"
