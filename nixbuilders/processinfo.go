package nixbuilders

import (
	"time"

	"github.com/alecthomas/units"
	"github.com/shirou/gopsutil/process"
)

type ProcessInfo struct {
	CpuPercent    float64
	MemoryPercent float32
	MemoryBytes   units.MetricBytes
	RawMemInfo    *process.MemoryInfoStat
	CreateTime    time.Time
	Errors        []error
}

func NewInfo(item ActiveUser) ProcessInfo {

	errors := []error{}

	percent, err := item.RootProcess.CPUPercent()
	if err != nil {
		errors = append(errors, err)
		percent = 0
	}

	mem, err := item.RootProcess.MemoryPercent()
	if err != nil {
		errors = append(errors, err)
		mem = 0
	}

	memSize, err := item.RootProcess.MemoryInfo()
	var memVms units.MetricBytes
	if err == nil && memSize != nil {
		memVms = units.MetricBytes(int64(memSize.VMS))
	} else if err != nil {
		memVms = 0
		errors = append(errors, err)
	} else {
		memVms = 0
	}

	createTime, err := item.RootProcess.CreateTime()
	if err != nil {
		errors = append(errors, err)
		createTime = 0
	}

	return ProcessInfo{
		CpuPercent:    percent,
		MemoryPercent: mem,
		MemoryBytes:   memVms,
		RawMemInfo:    memSize,
		CreateTime:    time.UnixMilli(createTime),
		Errors:        errors,
	}
}
