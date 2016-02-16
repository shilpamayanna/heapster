// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"time"

	source_api "k8s.io/heapster/sources/api"
)

// Stub out for testing
var timeSince = time.Since

var statMetrics = []SupportedStatMetric{
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "uptime",
			Description: "Number of milliseconds since the container was started",
			Type:        MetricCumulative,
			ValueType:   ValueInt64,
			Units:       UnitsMilliseconds,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return !spec.CreationTime.IsZero()
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			return []InternalPoint{{Value: timeSince(spec.CreationTime).Nanoseconds() / time.Millisecond.Nanoseconds()}}
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "cpu/usage",
			Description: "Cumulative CPU usage on all cores",
			Type:        MetricCumulative,
			ValueType:   ValueInt64,
			Units:       UnitsNanoseconds,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasCpu
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			return []InternalPoint{{Value: int64(stat.Cpu.Usage.Total)}}
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "cpu/limit",
			Description: "CPU hard limit in millicores.",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsCount,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasCpu && (spec.Cpu.Limit > 0)
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			// Normalize to a conversion factor of 1000.
			return []InternalPoint{{Value: int64(spec.Cpu.Limit*1000) / 1024}}
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "cpu/request",
			Description: "CPU request (the guaranteed amount of resources) in millicores. This metric is Kubernetes specific.",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsCount,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.CpuRequest > 0
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			return []InternalPoint{{Value: spec.CpuRequest}}
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "memory/usage",
			Description: "Total memory usage",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasMemory
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			return []InternalPoint{{Value: int64(stat.Memory.Usage)}}
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "memory/working_set",
			Description: "Total working set usage. Working set is the memory being used and not easily dropped by the kernel",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasMemory
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			return []InternalPoint{{Value: int64(stat.Memory.WorkingSet)}}
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "memory/limit",
			Description: "Memory hard limit in bytes.",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasMemory && (spec.Memory.Limit > 0)
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			return []InternalPoint{{Value: int64(spec.Memory.Limit)}}
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "memory/request",
			Description: "Memory request (the guaranteed amount of resources) in bytes. This metric is Kubernetes specific.",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.MemoryRequest > 0
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			return []InternalPoint{{Value: spec.MemoryRequest}}
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "memory/page_faults",
			Description: "Number of page faults",
			Type:        MetricCumulative,
			ValueType:   ValueInt64,
			Units:       UnitsCount,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasMemory
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			return []InternalPoint{{Value: int64(stat.Memory.ContainerData.Pgfault)}}
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "memory/major_page_faults",
			Description: "Number of major page faults",
			Type:        MetricCumulative,
			ValueType:   ValueInt64,
			Units:       UnitsCount,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasMemory
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			return []InternalPoint{{Value: int64(stat.Memory.ContainerData.Pgmajfault)}}
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "network/rx",
			Description: "Cumulative number of bytes received over the network",
			Type:        MetricCumulative,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasNetwork
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Network.Interfaces))
			for _, intf := range stat.Network.Interfaces {
				result = append(result, InternalPoint{
					Value: int64(intf.RxBytes),
					Labels: map[string]string{
						LabelResourceID.Key: intf.Name,
					},
				})
			}
			return result
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "network/rx_errors",
			Description: "Cumulative number of errors while receiving over the network",
			Type:        MetricCumulative,
			ValueType:   ValueInt64,
			Units:       UnitsCount,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasNetwork
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Network.Interfaces))
			for _, intf := range stat.Network.Interfaces {
				result = append(result, InternalPoint{
					Value: int64(intf.RxErrors),
					Labels: map[string]string{
						LabelResourceID.Key: intf.Name,
					},
				})
			}
			return result
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "network/tx",
			Description: "Cumulative number of bytes sent over the network",
			Type:        MetricCumulative,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasNetwork
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Network.Interfaces))
			for _, intf := range stat.Network.Interfaces {
				result = append(result, InternalPoint{
					Value: int64(intf.TxBytes),
					Labels: map[string]string{
						LabelResourceID.Key: intf.Name,
					},
				})
			}
			return result
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "network/tx_errors",
			Description: "Cumulative number of errors while sending over the network",
			Type:        MetricCumulative,
			ValueType:   ValueInt64,
			Units:       UnitsCount,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasNetwork
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Network.Interfaces))
			for _, intf := range stat.Network.Interfaces {
				result = append(result, InternalPoint{
					Value: int64(intf.TxErrors),
					Labels: map[string]string{
						LabelResourceID.Key: intf.Name,
					},
				})
			}
			return result
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/usage",
			Description: "Total number of bytes consumed on a filesystem",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.Usage),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/limit",
			Description: "The total size of filesystem in bytes",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.Limit),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/reads_completed",
			Description: "The total number of reads completed",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.ReadsCompleted),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/reads_merged",
			Description: "The total number of reads merged",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.ReadsMerged),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/sectors_read",
			Description: "The total number of sectors read",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.SectorsRead),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/read_time",
			Description: "The total number of milliseconds spent by all reads",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsMilliseconds,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.ReadTime),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/writes_completed",
			Description: "The total number of writes completed",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.WritesCompleted),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/writes_merged",
			Description: "The total number of writes merged",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.WritesMerged),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/sectors_written",
			Description: "The total number of sectors written",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.SectorsWritten),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/write_time",
			Description: "The total number of milliseconds spent by all writes",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsMilliseconds,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.WriteTime),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/io_in_progress",
			Description: "Number of I/Os in progress",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsBytes,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.IoInProgress),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/io_time",
			Description: "Number of milliseconds spent doing I/Os",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsMilliseconds,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.IoTime),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
		OnlyExportIfChanged: true,
	},
	{
		MetricDescriptor: MetricDescriptor{
			Name:        "filesystem/weighted_io_time",
			Description: "Number of weighted milliseconds spent doing I/Os",
			Type:        MetricGauge,
			ValueType:   ValueInt64,
			Units:       UnitsMilliseconds,
			Labels:      metricLabels,
		},
		HasValue: func(spec *source_api.ContainerSpec) bool {
			return spec.HasFilesystem
		},
		GetValue: func(spec *source_api.ContainerSpec, stat *source_api.ContainerStats) []InternalPoint {
			result := make([]InternalPoint, 0, len(stat.Filesystem))
			for _, fs := range stat.Filesystem {
				result = append(result, InternalPoint{
					Value: int64(fs.WeightedIoTime),
					Labels: map[string]string{
						LabelResourceID.Key: fs.Device,
					},
				})
			}
			return result
		},
		OnlyExportIfChanged: true,
	},

	// TODO(vmarmol): DiskIO stats if we find those useful and know how to export them in a user-friendly way.
}

// TODO: Add Status metrics - restarts, OOMs, etc.

func SupportedStatMetrics() []SupportedStatMetric {
	result := make([]SupportedStatMetric, len(statMetrics))
	copy(result, statMetrics)
	return result
}
