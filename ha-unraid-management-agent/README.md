# Unraid Management Agent - Home Assistant Integration

[![hacs_badge](https://img.shields.io/badge/HACS-Custom-orange.svg)](https://github.com/custom-components/hacs)
[![GitHub Release](https://img.shields.io/github/release/ruaandeysel/unraid-management-agent.svg)](https://github.com/ruaandeysel/unraid-management-agent/releases)
[![License](https://img.shields.io/github/license/ruaandeysel/unraid-management-agent.svg)](../LICENSE)

Complete Home Assistant custom integration for monitoring and controlling Unraid servers via the Unraid Management Agent.

## Features

### 🔍 Comprehensive Monitoring
- **System Metrics**: CPU usage, RAM usage, CPU temperature, uptime
- **Array Status**: Array state, disk usage, parity check progress
- **GPU Metrics**: GPU utilization, temperature, power consumption (Intel/NVIDIA/AMD)
- **Network**: Interface status, bandwidth monitoring (RX/TX)
- **UPS**: Battery level, load, runtime estimation
- **Containers**: Docker container status and metrics
- **Virtual Machines**: VM status and resource allocation

### 🎮 Full Control
- **Docker Containers**: Start, stop, restart containers via switches
- **Virtual Machines**: Start, stop, restart VMs via switches
- **Array Management**: Start/stop array with buttons
- **Parity Checks**: Start/stop parity checks with buttons

### ⚡ Real-Time Updates
- **WebSocket Support**: Instant state updates (<1s latency)
- **Automatic Fallback**: Falls back to REST API polling if WebSocket fails
- **Exponential Backoff**: Smart reconnection strategy
- **No Data Loss**: Seamless transition between WebSocket and polling

### 🏠 Home Assistant Native
- **UI Configuration**: No YAML required for setup
- **Device Grouping**: All entities grouped under single device
- **Proper Device Classes**: Temperature, power, battery, duration, etc.
- **State Classes**: Support for statistics and long-term data
- **MDI Icons**: Beautiful Material Design Icons for all entities
- **Extra Attributes**: Contextual information for each entity

## Quick Start

### Prerequisites

1. **Unraid Management Agent** installed and running on your Unraid server
2. **Home Assistant** 2023.1 or newer
3. **Network access** between Home Assistant and Unraid server

### Installation

See [INSTALLATION.md](INSTALLATION.md) for detailed instructions.

**Quick Install via HACS:**
1. Add custom repository: `https://github.com/ruaandeysel/unraid-management-agent`
2. Install "Unraid Management Agent"
3. Restart Home Assistant
4. Add integration via UI

### Configuration

1. Go to **Settings** → **Devices & Services**
2. Click **+ Add Integration**
3. Search for "Unraid Management Agent"
4. Enter your Unraid server details:
   - Host: `192.168.1.100` (your Unraid IP)
   - Port: `8043` (default)
   - Update Interval: `30` seconds
   - Enable WebSocket: `true` (recommended)

## Documentation

- **[Installation Guide](INSTALLATION.md)** - Detailed installation and configuration
- **[Examples](EXAMPLES.md)** - Automations, dashboards, and scripts
- **[Integration README](custom_components/unraid_management_agent/README.md)** - Entity list and services
- **[WebSocket Test Results](../WEBSOCKET_TEST_RESULTS.md)** - WebSocket implementation details

## Entity Overview

### Sensors (13+ entities)

**System Sensors (4)**
- CPU Usage (%)
- RAM Usage (%)
- CPU Temperature (°C)
- Uptime (seconds)

**Array Sensors (2)**
- Array Usage (%)
- Parity Check Progress (%)

**GPU Sensors (4, conditional)**
- GPU Name
- GPU Utilization (%)
- GPU CPU Temperature (°C)
- GPU Power (W)

**UPS Sensors (3, conditional)**
- UPS Battery (%)
- UPS Load (%)
- UPS Runtime (seconds)

**Network Sensors (dynamic)**
- Network {interface} RX (bytes)
- Network {interface} TX (bytes)

### Binary Sensors (7+ entities)

**Array Binary Sensors (3)**
- Array Started (on/off)
- Parity Check Running (on/off)
- Parity Valid (problem indicator)

**UPS Binary Sensor (1, conditional)**
- UPS Connected (on/off)

**Container Binary Sensors (dynamic)**
- Container {name} (running/stopped)

**VM Binary Sensors (dynamic)**
- VM {name} (running/stopped)

**Network Binary Sensors (dynamic)**
- Network {interface} (up/down)

### Switches (dynamic)

- Container {name} - Start/stop Docker containers
- VM {name} - Start/stop virtual machines

### Buttons (4)

- Start Array
- Stop Array
- Start Parity Check
- Stop Parity Check

## Example Automations

### High CPU Alert

```yaml
automation:
  - alias: "Unraid High CPU Alert"
    trigger:
      - platform: numeric_state
        entity_id: sensor.unraid_cpu_usage
        above: 80
        for:
          minutes: 5
    action:
      - service: notify.mobile_app
        data:
          title: "⚠️ Unraid Alert"
          message: "CPU usage is {{ states('sensor.unraid_cpu_usage') }}%"
```

### UPS Graceful Shutdown

```yaml
automation:
  - alias: "Unraid UPS Critical Shutdown"
    trigger:
      - platform: numeric_state
        entity_id: sensor.unraid_ups_battery
        below: 10
    action:
      - service: switch.turn_off
        target:
          entity_id: all
        data:
          domain: switch
      - delay:
          seconds: 30
      - service: button.press
        target:
          entity_id: button.unraid_array_stop
```

See [EXAMPLES.md](EXAMPLES.md) for more automation examples, dashboard configurations, and scripts.

## Architecture

### Components

- **API Client** (`api_client.py`) - REST API communication with aiohttp
- **WebSocket Client** (`websocket_client.py`) - Real-time event streaming
- **Data Coordinator** (`__init__.py`) - Data management and updates
- **Config Flow** (`config_flow.py`) - UI-based configuration
- **Platforms**:
  - `sensor.py` - System, array, GPU, UPS, network sensors
  - `binary_sensor.py` - Status indicators
  - `switch.py` - Container and VM control
  - `button.py` - Array and parity check control

### Data Flow

```
Unraid Server (REST API + WebSocket)
           ↓
    API Client / WebSocket Client
           ↓
    Data Update Coordinator
           ↓
    Entity Platforms (Sensor, Binary Sensor, Switch, Button)
           ↓
    Home Assistant UI
```

### Update Strategy

1. **Initial Load**: REST API fetches all data on startup
2. **Real-Time Updates**: WebSocket receives events and updates coordinator
3. **Fallback Polling**: REST API polls at configured interval if WebSocket fails
4. **Control Actions**: REST API sends commands, coordinator refreshes immediately

## Troubleshooting

### Common Issues

**Cannot Connect**
- Verify Unraid Management Agent is running: `curl http://<ip>:8043/api/v1/health`
- Check firewall rules allow port 8043
- Ensure Home Assistant can reach Unraid server

**WebSocket Not Working**
- Check logs for WebSocket errors
- Verify no proxy blocking WebSocket
- Integration will fall back to REST polling automatically

**Entities Not Updating**
- Check update interval in options
- Verify WebSocket connection in logs
- Test REST API manually

**Missing Entities**
- Verify resources exist on Unraid (containers, VMs, GPU, UPS)
- Reload integration
- Check logs for entity creation errors

See [INSTALLATION.md](INSTALLATION.md) for detailed troubleshooting.

## Development

### Project Structure

```
ha-unraid-management-agent/
├── custom_components/
│   └── unraid_management_agent/
│       ├── __init__.py           # Integration setup
│       ├── api_client.py         # REST API client
│       ├── binary_sensor.py      # Binary sensor platform
│       ├── button.py             # Button platform
│       ├── config_flow.py        # Configuration flow
│       ├── const.py              # Constants
│       ├── manifest.json         # Integration metadata
│       ├── sensor.py             # Sensor platform
│       ├── strings.json          # Translations
│       ├── switch.py             # Switch platform
│       └── websocket_client.py   # WebSocket client
├── EXAMPLES.md                   # Automation examples
├── INSTALLATION.md               # Installation guide
└── README.md                     # This file
```

### Testing

1. Install in development mode
2. Enable debug logging:
   ```yaml
   logger:
     default: info
     logs:
       custom_components.unraid_management_agent: debug
   ```
3. Check logs for errors
4. Test all entity types
5. Test control operations
6. Test WebSocket reconnection

## Support

- **Issues**: [GitHub Issues](https://github.com/ruaandeysel/unraid-management-agent/issues)
- **Documentation**: [GitHub Wiki](https://github.com/ruaandeysel/unraid-management-agent/wiki)
- **Discussions**: [GitHub Discussions](https://github.com/ruaandeysel/unraid-management-agent/discussions)

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## Credits

Developed by [@ruaandeysel](https://github.com/ruaandeysel)

## Changelog

### Version 1.0.0 (2025-01-02)

**Initial Release**

✅ **Core Features**
- Complete Home Assistant integration with UI configuration
- REST API client with all monitoring and control endpoints
- WebSocket client with real-time event streaming
- Automatic fallback to REST polling

✅ **Monitoring**
- System sensors (CPU, RAM, temperature, uptime)
- Array sensors (usage, parity progress)
- GPU sensors (name, utilization, temperature, power)
- UPS sensors (battery, load, runtime)
- Network sensors (RX/TX per interface)

✅ **Status Indicators**
- Array binary sensors (started, parity check, parity valid)
- UPS binary sensor (connected)
- Container binary sensors (running state)
- VM binary sensors (running state)
- Network binary sensors (interface up/down)

✅ **Control**
- Container switches (start/stop)
- VM switches (start/stop)
- Array buttons (start/stop)
- Parity check buttons (start/stop)

✅ **Quality**
- Proper device classes and state classes
- MDI icons for all entities
- Extra state attributes
- Error handling and logging
- Comprehensive documentation
- Example automations and dashboards

