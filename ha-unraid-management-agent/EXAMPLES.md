# Examples - Unraid Management Agent for Home Assistant

This document provides example automations, dashboards, and scripts for the Unraid Management Agent integration.

## Table of Contents

- [Automations](#automations)
- [Dashboard Cards](#dashboard-cards)
- [Scripts](#scripts)
- [Notifications](#notifications)

## Automations

### System Monitoring

#### High CPU Usage Alert

```yaml
automation:
  - alias: "Unraid: High CPU Usage Alert"
    description: "Alert when CPU usage is above 80% for 5 minutes"
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
          message: "CPU usage is {{ states('sensor.unraid_cpu_usage') }}% for 5 minutes"
          data:
            priority: high
```

#### High RAM Usage Alert

```yaml
automation:
  - alias: "Unraid: High RAM Usage Alert"
    description: "Alert when RAM usage exceeds 90%"
    trigger:
      - platform: numeric_state
        entity_id: sensor.unraid_ram_usage
        above: 90
    action:
      - service: notify.mobile_app
        data:
          title: "⚠️ Unraid Alert"
          message: "RAM usage is {{ states('sensor.unraid_ram_usage') }}%"
```

#### High CPU Temperature Alert

```yaml
automation:
  - alias: "Unraid: High CPU Temperature"
    description: "Alert when CPU temperature exceeds 75°C"
    trigger:
      - platform: numeric_state
        entity_id: sensor.unraid_cpu_temperature
        above: 75
    action:
      - service: notify.mobile_app
        data:
          title: "🌡️ Unraid Temperature Alert"
          message: "CPU temperature is {{ states('sensor.unraid_cpu_temperature') }}°C"
```

### Array Management

#### Array Started Notification

```yaml
automation:
  - alias: "Unraid: Array Started"
    description: "Notify when array starts"
    trigger:
      - platform: state
        entity_id: binary_sensor.unraid_array_started
        to: "on"
    action:
      - service: notify.mobile_app
        data:
          title: "✅ Unraid Array Started"
          message: "The Unraid array has been started"
```

#### Array Stopped Notification

```yaml
automation:
  - alias: "Unraid: Array Stopped"
    description: "Notify when array stops"
    trigger:
      - platform: state
        entity_id: binary_sensor.unraid_array_started
        to: "off"
    action:
      - service: notify.mobile_app
        data:
          title: "⚠️ Unraid Array Stopped"
          message: "The Unraid array has been stopped"
          data:
            priority: high
```

#### Parity Check Started

```yaml
automation:
  - alias: "Unraid: Parity Check Started"
    description: "Notify when parity check starts"
    trigger:
      - platform: state
        entity_id: binary_sensor.unraid_parity_check_running
        to: "on"
    action:
      - service: notify.mobile_app
        data:
          title: "🔍 Unraid Parity Check Started"
          message: "Parity check has been initiated"
```

#### Parity Check Completed

```yaml
automation:
  - alias: "Unraid: Parity Check Completed"
    description: "Notify when parity check completes"
    trigger:
      - platform: state
        entity_id: binary_sensor.unraid_parity_check_running
        to: "off"
        for:
          seconds: 30
    condition:
      - condition: template
        value_template: "{{ trigger.from_state.state == 'on' }}"
    action:
      - service: notify.mobile_app
        data:
          title: "✅ Unraid Parity Check Complete"
          message: "Parity check has finished"
```

#### Parity Invalid Alert

```yaml
automation:
  - alias: "Unraid: Parity Invalid"
    description: "Alert when parity becomes invalid"
    trigger:
      - platform: state
        entity_id: binary_sensor.unraid_parity_valid
        to: "on"  # ON means problem (PROBLEM device class)
    action:
      - service: notify.mobile_app
        data:
          title: "🚨 Unraid Parity Invalid"
          message: "Parity is invalid! Check your array immediately."
          data:
            priority: high
```

### Container Management

#### Start Container When Array Starts

```yaml
automation:
  - alias: "Unraid: Start Plex on Array Start"
    description: "Start Plex container when array starts"
    trigger:
      - platform: state
        entity_id: binary_sensor.unraid_array_started
        to: "on"
    action:
      - service: switch.turn_on
        target:
          entity_id: switch.unraid_container_plex
```

#### Stop Containers Before Array Stop

```yaml
automation:
  - alias: "Unraid: Stop Containers Before Array Stop"
    description: "Stop all containers before stopping array"
    trigger:
      - platform: event
        event_type: call_service
        event_data:
          domain: button
          service: press
          service_data:
            entity_id: button.unraid_array_stop
    action:
      - service: switch.turn_off
        target:
          entity_id:
            - switch.unraid_container_plex
            - switch.unraid_container_nginx
            - switch.unraid_container_homeassistant
```

#### Container Stopped Alert

```yaml
automation:
  - alias: "Unraid: Critical Container Stopped"
    description: "Alert when critical container stops unexpectedly"
    trigger:
      - platform: state
        entity_id: binary_sensor.unraid_container_plex
        to: "off"
        for:
          minutes: 2
    condition:
      - condition: state
        entity_id: binary_sensor.unraid_array_started
        state: "on"
    action:
      - service: notify.mobile_app
        data:
          title: "⚠️ Unraid Container Stopped"
          message: "Plex container has stopped unexpectedly"
```

### UPS Monitoring

#### UPS Battery Low

```yaml
automation:
  - alias: "Unraid: UPS Battery Low"
    description: "Alert when UPS battery is below 20%"
    trigger:
      - platform: numeric_state
        entity_id: sensor.unraid_ups_battery
        below: 20
    action:
      - service: notify.mobile_app
        data:
          title: "🔋 UPS Battery Low"
          message: "UPS battery is at {{ states('sensor.unraid_ups_battery') }}%"
          data:
            priority: high
```

#### UPS Disconnected

```yaml
automation:
  - alias: "Unraid: UPS Disconnected"
    description: "Alert when UPS disconnects"
    trigger:
      - platform: state
        entity_id: binary_sensor.unraid_ups_connected
        to: "off"
    action:
      - service: notify.mobile_app
        data:
          title: "⚠️ UPS Disconnected"
          message: "UPS has been disconnected from Unraid server"
          data:
            priority: high
```

#### UPS Power Failure - Graceful Shutdown

```yaml
automation:
  - alias: "Unraid: UPS Power Failure Shutdown"
    description: "Gracefully shutdown when UPS battery is critical"
    trigger:
      - platform: numeric_state
        entity_id: sensor.unraid_ups_battery
        below: 10
    condition:
      - condition: state
        entity_id: binary_sensor.unraid_ups_connected
        state: "on"
    action:
      # Stop all containers
      - service: switch.turn_off
        target:
          entity_id: all
        data:
          domain: switch
      # Wait for containers to stop
      - delay:
          seconds: 30
      # Stop array
      - service: button.press
        target:
          entity_id: button.unraid_array_stop
      # Notify
      - service: notify.mobile_app
        data:
          title: "🚨 UPS Critical - Shutdown Initiated"
          message: "UPS battery at {{ states('sensor.unraid_ups_battery') }}%. Shutting down Unraid."
```

## Dashboard Cards

### System Overview Card

```yaml
type: entities
title: Unraid System
entities:
  - entity: sensor.unraid_cpu_usage
    name: CPU Usage
  - entity: sensor.unraid_ram_usage
    name: RAM Usage
  - entity: sensor.unraid_cpu_temperature
    name: CPU Temperature
  - entity: sensor.unraid_uptime
    name: Uptime
  - entity: binary_sensor.unraid_array_started
    name: Array Status
```

### Array Status Card

```yaml
type: entities
title: Unraid Array
entities:
  - entity: binary_sensor.unraid_array_started
    name: Array Started
  - entity: sensor.unraid_array_usage
    name: Array Usage
  - entity: binary_sensor.unraid_parity_check_running
    name: Parity Check
  - entity: sensor.unraid_parity_progress
    name: Parity Progress
  - entity: binary_sensor.unraid_parity_valid
    name: Parity Valid
  - type: divider
  - entity: button.unraid_array_start
    name: Start Array
  - entity: button.unraid_array_stop
    name: Stop Array
  - entity: button.unraid_parity_check_start
    name: Start Parity Check
  - entity: button.unraid_parity_check_stop
    name: Stop Parity Check
```

### GPU Monitoring Card

```yaml
type: entities
title: Unraid GPU
entities:
  - entity: sensor.unraid_gpu_name
    name: GPU Model
  - entity: sensor.unraid_gpu_utilization
    name: GPU Usage
  - entity: sensor.unraid_gpu_cpu_temperature
    name: CPU Temperature
  - entity: sensor.unraid_gpu_power
    name: Power Draw
```

### Container Control Card

```yaml
type: entities
title: Docker Containers
entities:
  - entity: binary_sensor.unraid_container_plex
    name: Plex Status
  - entity: switch.unraid_container_plex
    name: Plex Control
  - type: divider
  - entity: binary_sensor.unraid_container_nginx
    name: Nginx Status
  - entity: switch.unraid_container_nginx
    name: Nginx Control
```

### UPS Status Card

```yaml
type: entities
title: UPS Status
entities:
  - entity: binary_sensor.unraid_ups_connected
    name: Connected
  - entity: sensor.unraid_ups_battery
    name: Battery Level
  - entity: sensor.unraid_ups_load
    name: Load
  - entity: sensor.unraid_ups_runtime
    name: Runtime
```

## Scripts

### Restart Container Script

```yaml
script:
  restart_unraid_container:
    alias: "Restart Unraid Container"
    description: "Restart a specific container"
    fields:
      container:
        description: "Container switch entity"
        example: "switch.unraid_container_plex"
    sequence:
      - service: switch.turn_off
        target:
          entity_id: "{{ container }}"
      - delay:
          seconds: 5
      - service: switch.turn_on
        target:
          entity_id: "{{ container }}"
```

### Safe Array Shutdown Script

```yaml
script:
  safe_array_shutdown:
    alias: "Safe Unraid Array Shutdown"
    description: "Safely stop all containers and array"
    sequence:
      # Stop all containers
      - service: switch.turn_off
        target:
          entity_id: all
        data:
          domain: switch
      # Wait for containers
      - delay:
          seconds: 30
      # Stop array
      - service: button.press
        target:
          entity_id: button.unraid_array_stop
      # Notify
      - service: notify.mobile_app
        data:
          title: "Unraid Shutdown"
          message: "Array shutdown initiated"
```

## Notifications

### Telegram Notification

```yaml
automation:
  - alias: "Unraid: Telegram Alerts"
    trigger:
      - platform: numeric_state
        entity_id: sensor.unraid_cpu_usage
        above: 90
    action:
      - service: notify.telegram
        data:
          message: "⚠️ Unraid CPU at {{ states('sensor.unraid_cpu_usage') }}%"
```

### Discord Notification

```yaml
automation:
  - alias: "Unraid: Discord Alerts"
    trigger:
      - platform: state
        entity_id: binary_sensor.unraid_parity_valid
        to: "on"
    action:
      - service: notify.discord
        data:
          message: "🚨 Unraid parity is invalid!"
```

## Tips

1. **Test Automations**: Always test automations in a safe environment before deploying
2. **Adjust Thresholds**: Customize alert thresholds based on your hardware
3. **Use Conditions**: Add conditions to prevent false alerts
4. **Graceful Shutdowns**: Always stop containers before stopping the array
5. **Monitor Logs**: Check Home Assistant logs for integration issues

