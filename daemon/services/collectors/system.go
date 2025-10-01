package collectors

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ruaandeysel/unraid-management-agent/daemon/domain"
	"github.com/ruaandeysel/unraid-management-agent/daemon/dto"
	"github.com/ruaandeysel/unraid-management-agent/daemon/lib"
	"github.com/ruaandeysel/unraid-management-agent/daemon/logger"
)

type SystemCollector struct {
	ctx *domain.Context
}

func NewSystemCollector(ctx *domain.Context) *SystemCollector {
	return &SystemCollector{ctx: ctx}
}

func (c *SystemCollector) Start(interval time.Duration) {
	logger.Info("Starting system collector (interval: %v)", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.Collect()
	}
}

func (c *SystemCollector) Collect() {
	if c.ctx.MockMode {
		logger.Debug("Mock mode: system collection skipped")
		return
	}

	logger.Debug("Collecting system data...")

	// Collect system info
	systemInfo, err := c.collectSystemInfo()
	if err != nil {
		logger.Error("Failed to collect system info: %v", err)
		return
	}

	// Publish event
	c.ctx.Hub.Pub(systemInfo, "system_update")
	logger.Debug("Published system_update event")
}

func (c *SystemCollector) collectSystemInfo() (*dto.SystemInfo, error) {
	info := &dto.SystemInfo{}

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		logger.Warning("Failed to get hostname", "error", err)
		info.Hostname = "unknown"
	} else {
		info.Hostname = hostname
	}

	// Get uptime
	uptime, err := c.getUptime()
	if err != nil {
		logger.Warning("Failed to get uptime", "error", err)
	} else {
		info.Uptime = uptime
	}

	// Get CPU info
	cpuPercent, err := c.getCPUInfo()
	if err != nil {
		logger.Warning("Failed to get CPU info", "error", err)
	} else {
		info.CPUUsage = cpuPercent
	}

	// Get memory info
	memUsed, memTotal, memFree, err := c.getMemoryInfo()
	if err != nil {
		logger.Warning("Failed to get memory info", "error", err)
	} else {
		info.RAMUsed = memUsed
		info.RAMTotal = memTotal
		info.RAMFree = memFree
		if memTotal > 0 {
			info.RAMUsage = float64(memUsed) / float64(memTotal) * 100
		}
	}

	// Get temperatures
	temperatures, err := c.getTemperatures()
	if err != nil {
		logger.Warning("Failed to get temperatures", "error", err)
	} else {
		// Extract CPU and motherboard temps if available
		for name, temp := range temperatures {
			if strings.Contains(strings.ToLower(name), "cpu") || strings.Contains(strings.ToLower(name), "core") {
				if info.CPUTemp == 0 || temp > info.CPUTemp {
					info.CPUTemp = temp
				}
			} else if strings.Contains(strings.ToLower(name), "motherboard") || strings.Contains(strings.ToLower(name), "mb") {
				info.MotherboardTemp = temp
			}
		}
	}

	// Get fan speeds
	fans, err := c.getFans()
	if err != nil {
		logger.Warning("Failed to get fan speeds", "error", err)
	} else {
		info.Fans = fans
	}

	// Set timestamp
	info.Timestamp = time.Now()

	return info, nil
}

func (c *SystemCollector) getUptime() (int64, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0, fmt.Errorf("invalid uptime format")
	}

	uptime, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, err
	}

	return int64(uptime), nil
}

func (c *SystemCollector) getCPUInfo() (float64, error) {
	// Get CPU usage by reading /proc/stat
	cpuPercent, err := c.calculateCPUPercent()
	if err != nil {
		logger.Warning("Failed to calculate CPU percent", "error", err)
		return 0, err
	}

	return cpuPercent, nil
}

func (c *SystemCollector) calculateCPUPercent() (float64, error) {
	// Read first snapshot
	stat1, err := c.readCPUStat()
	if err != nil {
		return 0, err
	}

	// Wait a short time
	time.Sleep(100 * time.Millisecond)

	// Read second snapshot
	stat2, err := c.readCPUStat()
	if err != nil {
		return 0, err
	}

	// Calculate usage
	total1 := stat1["user"] + stat1["nice"] + stat1["system"] + stat1["idle"] + stat1["iowait"] + stat1["irq"] + stat1["softirq"] + stat1["steal"]
	total2 := stat2["user"] + stat2["nice"] + stat2["system"] + stat2["idle"] + stat2["iowait"] + stat2["irq"] + stat2["softirq"] + stat2["steal"]

	idle1 := stat1["idle"] + stat1["iowait"]
	idle2 := stat2["idle"] + stat2["iowait"]

	totalDelta := total2 - total1
	idleDelta := idle2 - idle1

	if totalDelta == 0 {
		return 0, nil
	}

	usage := (float64(totalDelta-idleDelta) / float64(totalDelta)) * 100
	return usage, nil
}

func (c *SystemCollector) readCPUStat() (map[string]uint64, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) < 9 {
				return nil, fmt.Errorf("invalid cpu stat format")
			}

			stat := make(map[string]uint64)
			stat["user"], _ = strconv.ParseUint(fields[1], 10, 64)
			stat["nice"], _ = strconv.ParseUint(fields[2], 10, 64)
			stat["system"], _ = strconv.ParseUint(fields[3], 10, 64)
			stat["idle"], _ = strconv.ParseUint(fields[4], 10, 64)
			stat["iowait"], _ = strconv.ParseUint(fields[5], 10, 64)
			stat["irq"], _ = strconv.ParseUint(fields[6], 10, 64)
			stat["softirq"], _ = strconv.ParseUint(fields[7], 10, 64)
			stat["steal"], _ = strconv.ParseUint(fields[8], 10, 64)

			return stat, nil
		}
	}

	return nil, fmt.Errorf("cpu line not found in /proc/stat")
}

func (c *SystemCollector) getMemoryInfo() (uint64, uint64, uint64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, 0, err
	}
	defer file.Close()

	var memTotal, memFree, memBuffers, memCached uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		value, _ := strconv.ParseUint(fields[1], 10, 64)
		value *= 1024 // Convert from KB to bytes

		switch key {
		case "MemTotal":
			memTotal = value
		case "MemFree":
			memFree = value
		case "Buffers":
			memBuffers = value
		case "Cached":
			memCached = value
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, 0, err
	}

	// Calculate used memory (excluding buffers and cache)
	memUsed := memTotal - memFree - memBuffers - memCached
	// Calculate actual free (including buffers and cache)
	memActualFree := memFree + memBuffers + memCached

	return memUsed, memTotal, memActualFree, nil
}

func (c *SystemCollector) getTemperatures() (map[string]float64, error) {
	temperatures := make(map[string]float64)

	// Try using sensors command first
	output, err := lib.ExecCommandOutput("sensors", "-u")
	if err == nil {
		temperatures = c.parseSensorsOutput(output)
		if len(temperatures) > 0 {
			return temperatures, nil
		}
	}

	// Fallback: try reading from /sys/class/hwmon
	temperatures, err = c.readHwmonTemperatures()
	if err != nil {
		return nil, err
	}

	return temperatures, nil
}

func (c *SystemCollector) parseSensorsOutput(output string) map[string]float64 {
	temperatures := make(map[string]float64)
	lines := strings.Split(output, "\n")

	var currentChip string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// New chip/adapter
		if !strings.Contains(line, ":") && !strings.HasPrefix(line, " ") {
			currentChip = line
			continue
		}

		// Temperature input line
		if strings.Contains(line, "_input:") && currentChip != "" {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				valueStr := strings.TrimSpace(parts[1])
				if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
					// Create a friendly name
					name := fmt.Sprintf("%s_%s", currentChip, key)
					name = strings.ReplaceAll(name, " ", "_")
					temperatures[name] = value / 1000.0 // Convert from millidegrees
				}
			}
		}
	}

	return temperatures
}

func (c *SystemCollector) readHwmonTemperatures() (map[string]float64, error) {
	temperatures := make(map[string]float64)

	// Read from /sys/class/hwmon/hwmon*/temp*_input
	for i := 0; i < 10; i++ {
		for j := 1; j < 20; j++ {
			path := fmt.Sprintf("/sys/class/hwmon/hwmon%d/temp%d_input", i, j)
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			value, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
			if err != nil {
				continue
			}

			// Try to get label
			labelPath := fmt.Sprintf("/sys/class/hwmon/hwmon%d/temp%d_label", i, j)
			labelData, err := os.ReadFile(labelPath)
			label := fmt.Sprintf("hwmon%d_temp%d", i, j)
			if err == nil {
				label = strings.TrimSpace(string(labelData))
			}

			temperatures[label] = value / 1000.0 // Convert from millidegrees
		}
	}

	if len(temperatures) == 0 {
		return nil, fmt.Errorf("no temperature sensors found")
	}

	return temperatures, nil
}

func (c *SystemCollector) getFans() ([]dto.FanInfo, error) {
	fanMap := make(map[string]int)

	// Try using sensors command first
	output, err := lib.ExecCommandOutput("sensors", "-u")
	if err == nil {
		fanMap = c.parseFanSpeeds(output)
	}

	// If no fans found, try fallback
	if len(fanMap) == 0 {
		fanMap, err = c.readHwmonFanSpeeds()
		if err != nil {
			return nil, err
		}
	}

	// Convert map to slice
	fans := make([]dto.FanInfo, 0, len(fanMap))
	for name, rpm := range fanMap {
		fans = append(fans, dto.FanInfo{
			Name: name,
			RPM:  rpm,
		})
	}

	return fans, nil
}

func (c *SystemCollector) parseFanSpeeds(output string) map[string]int {
	fanSpeeds := make(map[string]int)
	lines := strings.Split(output, "\n")

	var currentChip string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// New chip/adapter
		if !strings.Contains(line, ":") && !strings.HasPrefix(line, " ") {
			currentChip = line
			continue
		}

		// Fan input line
		if strings.Contains(line, "fan") && strings.Contains(line, "_input:") && currentChip != "" {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				valueStr := strings.TrimSpace(parts[1])
				if value, err := strconv.Atoi(valueStr); err == nil {
					name := fmt.Sprintf("%s_%s", currentChip, key)
					name = strings.ReplaceAll(name, " ", "_")
					fanSpeeds[name] = value
				}
			}
		}
	}

	return fanSpeeds
}

func (c *SystemCollector) readHwmonFanSpeeds() (map[string]int, error) {
	fanSpeeds := make(map[string]int)

	// Read from /sys/class/hwmon/hwmon*/fan*_input
	for i := 0; i < 10; i++ {
		for j := 1; j < 20; j++ {
			path := fmt.Sprintf("/sys/class/hwmon/hwmon%d/fan%d_input", i, j)
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			value, err := strconv.Atoi(strings.TrimSpace(string(data)))
			if err != nil {
				continue
			}

			// Try to get label
			labelPath := fmt.Sprintf("/sys/class/hwmon/hwmon%d/fan%d_label", i, j)
			labelData, err := os.ReadFile(labelPath)
			label := fmt.Sprintf("hwmon%d_fan%d", i, j)
			if err == nil {
				label = strings.TrimSpace(string(labelData))
			}

			fanSpeeds[label] = value
		}
	}

	if len(fanSpeeds) == 0 {
		return nil, fmt.Errorf("no fan sensors found")
	}

	return fanSpeeds, nil
}
