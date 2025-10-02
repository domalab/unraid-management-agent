# Installation Guide - Unraid Management Agent for Home Assistant

This guide will walk you through installing and configuring the Unraid Management Agent Home Assistant integration.

## Prerequisites

### 1. Unraid Management Agent

First, ensure the Unraid Management Agent is installed and running on your Unraid server:

1. Install the plugin from Unraid Community Applications
2. Verify it's running: `http://<unraid-ip>:8043/api/v1/health`
3. You should see: `{"status": "healthy"}`

### 2. Home Assistant

- Home Assistant 2023.1 or newer
- Network access to your Unraid server
- HACS installed (for HACS installation method)

## Installation Methods

### Method 1: HACS (Recommended)

1. **Add Custom Repository**
   - Open HACS in Home Assistant
   - Click on "Integrations"
   - Click the three dots (⋮) in the top right
   - Select "Custom repositories"
   - Add repository URL: `https://github.com/ruaandeysel/unraid-management-agent`
   - Category: "Integration"
   - Click "Add"

2. **Install Integration**
   - Search for "Unraid Management Agent" in HACS
   - Click "Download"
   - Restart Home Assistant

3. **Configure Integration**
   - Go to **Settings** → **Devices & Services**
   - Click **+ Add Integration**
   - Search for "Unraid Management Agent"
   - Follow the configuration steps below

### Method 2: Manual Installation

1. **Download Integration**
   ```bash
   cd /config
   mkdir -p custom_components
   cd custom_components
   git clone https://github.com/ruaandeysel/unraid-management-agent.git
   mv unraid-management-agent/ha-unraid-management-agent/custom_components/unraid_management_agent .
   rm -rf unraid-management-agent
   ```

2. **Restart Home Assistant**
   - Go to **Settings** → **System** → **Restart**

3. **Configure Integration**
   - Go to **Settings** → **Devices & Services**
   - Click **+ Add Integration**
   - Search for "Unraid Management Agent"
   - Follow the configuration steps below

## Configuration

### Initial Setup

1. **Add Integration**
   - Navigate to **Settings** → **Devices & Services**
   - Click **+ Add Integration**
   - Search for "Unraid Management Agent"

2. **Enter Connection Details**
   - **Host**: IP address or hostname of your Unraid server (e.g., `192.168.1.100`)
   - **Port**: Port number (default: `8043`)
   - **Update Interval**: Polling interval in seconds (default: `30`)
   - **Enable WebSocket**: Enable real-time updates (recommended: `true`)

3. **Verify Connection**
   - The integration will test the connection
   - If successful, you'll see a success message
   - If failed, check the error message and verify:
     - Unraid server is accessible
     - Port is correct
     - Firewall allows connections

### Configuration Options

After initial setup, you can modify settings:

1. Go to **Settings** → **Devices & Services**
2. Find "Unraid Management Agent"
3. Click **Configure**
4. Modify:
   - **Update Interval**: How often to poll (when WebSocket unavailable)
   - **Enable WebSocket**: Toggle real-time updates

## Verification

### Check Integration Status

1. **Device Page**
   - Go to **Settings** → **Devices & Services**
   - Click on "Unraid Management Agent"
   - You should see the Unraid device with all entities

2. **Entity Count**
   - System sensors: 4 (CPU, RAM, temp, uptime)
   - Array sensors: 2 (usage, parity progress)
   - Array binary sensors: 3 (started, parity check, parity valid)
   - GPU sensors: 4 (if GPU available)
   - UPS sensors: 3 (if UPS connected)
   - UPS binary sensor: 1 (if UPS connected)
   - Container entities: 2 per container (binary sensor + switch)
   - VM entities: 2 per VM (binary sensor + switch)
   - Network entities: 3 per interface (RX sensor, TX sensor, binary sensor)
   - Array buttons: 4 (start/stop array, start/stop parity check)

3. **Test Entities**
   - Check sensor values are updating
   - Verify binary sensors show correct states
   - Test a switch (start/stop a container)
   - Test a button (if safe to do so)

### Check Logs

If you encounter issues, check the logs:

1. Go to **Settings** → **System** → **Logs**
2. Filter by "unraid_management_agent"
3. Look for errors or warnings

## Network Configuration

### Firewall Rules

Ensure your firewall allows:
- **REST API**: TCP port 8043 (or your configured port)
- **WebSocket**: TCP port 8043 (same port, WebSocket upgrade)

### VLAN Configuration

If Home Assistant and Unraid are on different VLANs:
1. Ensure routing between VLANs
2. Allow TCP port 8043 in firewall rules
3. Test connectivity: `curl http://<unraid-ip>:8043/api/v1/health`

## Troubleshooting

### Cannot Connect to Server

**Error**: "Failed to connect to the Unraid server"

**Solutions**:
1. Verify Unraid Management Agent is running:
   ```bash
   curl http://<unraid-ip>:8043/api/v1/health
   ```
2. Check firewall rules
3. Verify IP address and port
4. Ensure Home Assistant can reach Unraid server

### Timeout Errors

**Error**: "Connection timeout"

**Solutions**:
1. Check network latency
2. Verify no network issues
3. Increase timeout in integration code (advanced)

### WebSocket Not Working

**Symptoms**: Entities not updating in real-time

**Solutions**:
1. Check WebSocket is enabled in options
2. Verify no proxy blocking WebSocket
3. Check logs for WebSocket errors
4. Integration will fall back to REST polling

### Entities Not Appearing

**Symptoms**: Missing sensors or controls

**Solutions**:
1. Verify resources exist on Unraid:
   - Containers: Check Docker tab
   - VMs: Check VMs tab
   - GPU: Verify GPU is detected
   - UPS: Verify UPS is connected
2. Reload integration
3. Check logs for errors during entity creation

### State Not Updating

**Symptoms**: Sensor values are stale

**Solutions**:
1. Check update interval in options
2. Verify WebSocket connection (check logs)
3. Test REST API manually:
   ```bash
   curl http://<unraid-ip>:8043/api/v1/system
   ```
4. Reload integration

## Advanced Configuration

### Custom Update Intervals

For different update frequencies:
1. Go to integration options
2. Set update interval (seconds)
3. Recommended values:
   - Fast updates: 10-15 seconds
   - Normal: 30 seconds (default)
   - Slow: 60+ seconds

### Disable WebSocket

If WebSocket causes issues:
1. Go to integration options
2. Disable "Enable WebSocket"
3. Integration will use REST polling only

### Multiple Unraid Servers

To monitor multiple servers:
1. Add integration multiple times
2. Each server gets its own device
3. Entities are prefixed with hostname

## Next Steps

After installation:
1. Review [README.md](custom_components/unraid_management_agent/README.md) for entity list
2. Check [EXAMPLES.md](EXAMPLES.md) for automation examples
3. Create dashboards to visualize your Unraid server
4. Set up automations for monitoring and control

## Support

- **Issues**: [GitHub Issues](https://github.com/ruaandeysel/unraid-management-agent/issues)
- **Documentation**: [GitHub Wiki](https://github.com/ruaandeysel/unraid-management-agent/wiki)
- **Discussions**: [GitHub Discussions](https://github.com/ruaandeysel/unraid-management-agent/discussions)

