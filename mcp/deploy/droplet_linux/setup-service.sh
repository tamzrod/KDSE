#!/bin/bash
# KDSE MCP Server - Systemd Service Setup Script
# Usage: sudo ./setup-service.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "ERROR: Please run as root (use sudo)"
    exit 1
fi

echo "========================================"
echo "KDSE MCP Server - Systemd Service Setup"
echo "========================================"

# Create systemd service file
SERVICE_FILE="/etc/systemd/system/kdse-mcp.service"

cat > "$SERVICE_FILE" << EOF
[Unit]
Description=KDSE MCP Server
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=${SCRIPT_DIR}
ExecStart=/usr/bin/docker compose --file ${SCRIPT_DIR}/docker-compose.yml --profile http up -d
ExecStop=/usr/bin/docker compose --file ${SCRIPT_DIR}/docker-compose.yml down
TimeoutStartSec=0
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

echo "✓ Created systemd service file: $SERVICE_FILE"

# Reload systemd daemon
systemctl daemon-reload

echo "✓ Reloaded systemd daemon"

# Enable service
systemctl enable kdse-mcp.service

echo "✓ Service enabled"
echo ""
echo "========================================"
echo "Service Setup Complete!"
echo "========================================"
echo ""
echo "Commands:"
echo "  Start:    sudo systemctl start kdse-mcp"
echo "  Stop:     sudo systemctl stop kdse-mcp"
echo "  Status:   sudo systemctl status kdse-mcp"
echo "  Logs:     sudo journalctl -u kdse-mcp -f"
echo "  Restart:  sudo systemctl restart kdse-mcp"
echo ""
