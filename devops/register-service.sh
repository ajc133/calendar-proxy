#!/bin/bash

set -eu -o pipefail

SERVICE_NAME="calendar-proxy.service"
SYSTEMD_DIR="${XDG_CONFIG_HOME}/systemd/user"
mkdir -p "$SYSTEMD_DIR"
cp "./devops/${SERVICE_NAME}"  "${SYSTEMD_DIR}/${SERVICE_NAME}"
systemctl --user daemon-reload
systemctl --user enable "${SERVICE_NAME}"
systemctl --user start "${SERVICE_NAME}"
