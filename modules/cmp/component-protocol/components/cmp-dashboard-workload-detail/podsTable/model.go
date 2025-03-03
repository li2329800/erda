// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package podsTable

import (
	"context"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda/bundle"
	"github.com/erda-project/erda/modules/cmp"
	"github.com/erda-project/erda/modules/openapi/component-protocol/components/base"
)

type ComponentPodsTable struct {
	base.DefaultProvider
	sdk    *cptype.SDK
	bdl    *bundle.Bundle
	ctx    context.Context
	server cmp.SteveServer

	Type       string                 `json:"type,omitempty"`
	State      State                  `json:"state,omitempty"`
	Data       Data                   `json:"data,omitempty"`
	Props      Props                  `json:"props,omitempty"`
	Operations map[string]interface{} `json:"operations,omitempty"`
}

type State struct {
	ClusterName       string `json:"clusterName,omitempty"`
	WorkloadID        string `json:"workloadId,omitempty"`
	PageNo            int    `json:"pageNo"`
	PageSize          int    `json:"pageSize"`
	Sorter            Sorter `json:"sorterData,omitempty"`
	Total             int    `json:"total"`
	PodsTableURLQuery string `json:"podsTableURLQuery,omitempty"`
}

type Sorter struct {
	Field string `json:"field,omitempty"`
	Order string `json:"order,omitempty"`
}

type Data struct {
	List []Item `json:"list,omitempty"`
}

type Item struct {
	ID                string  `json:"id,omitempty"`
	Status            Status  `json:"status,omitempty"`
	Name              Link    `json:"name,omitempty"`
	Namespace         string  `json:"namespace,omitempty"`
	IP                string  `json:"ip,omitempty"`
	Age               string  `json:"age,omitempty"`
	CPURequests       string  `json:"cpuRequests,omitempty"`
	CPURequestsNum    int64   `json:"CPURequestsNum,omitempty"`
	CPUPercent        Percent `json:"cpuPercent,omitempty"`
	CPULimits         string  `json:"cpuLimits,omitempty"`
	CPULimitsNum      int64   `json:"CPULimitsNum,omitempty"`
	MemoryRequests    string  `json:"memoryRequests,omitempty"`
	MemoryRequestsNum int64   `json:"MemoryRequestsNum,omitempty"`
	MemoryPercent     Percent `json:"memoryPercent,omitempty"`
	MemoryLimits      string  `json:"memoryLimits,omitempty"`
	MemoryLimitsNum   int64   `json:"MemoryLimitsNum,omitempty"`
	Ready             string  `json:"ready,omitempty"`
	NodeName          string  `json:"nodeName,omitempty"`
}

type Status struct {
	RenderType string      `json:"renderType,omitempty"`
	Size       string      `json:"size,omitempty"`
	Value      StatusValue `json:"value,omitempty"`
}

type StatusValue struct {
	Label string `json:"label,omitempty"`
	Color string `json:"color,omitempty"`
}

type StyleConfig struct {
	Color string `json:"color,omitempty"`
}

type Link struct {
	RenderType string                 `json:"renderType,omitempty"`
	Value      string                 `json:"value,omitempty"`
	Operations map[string]interface{} `json:"operations,omitempty"`
}

type LinkOperation struct {
	Command Command `json:"command,omitempty"`
	Reload  bool    `json:"reload"`
}

type Command struct {
	Key     string       `json:"key,omitempty"`
	Target  string       `json:"target,omitempty"`
	State   CommandState `json:"state,omitempty"`
	JumpOut bool         `json:"jumpOut"`
}

type CommandState struct {
	Params map[string]string `json:"params,omitempty"`
	Query  map[string]string `json:"query,omitempty"`
}

type Percent struct {
	RenderType string `json:"renderType,omitempty"`
	Value      string `json:"value,omitempty"`
	Tip        string `json:"tip,omitempty"`
	Status     string `json:"status,omitempty"`
}

type Props struct {
	RequestIgnore   []string               `json:"requestIgnore,omitempty"`
	PageSizeOptions []string               `json:"pageSizeOptions,omitempty"`
	Columns         []Column               `json:"columns,omitempty"`
	RowKey          string                 `json:"rowKey,omitempty"`
	Operations      map[string]interface{} `json:"operations,omitempty"`
	SortDirections  []string               `json:"sortDirections,omitempty"`
}

type Column struct {
	DataIndex string `json:"dataIndex,omitempty"`
	Title     string `json:"title,omitempty"`
	Width     int    `json:"width"`
	Sorter    bool   `json:"sorter"`
}

type Operation struct {
	Key    string `json:"key,omitempty"`
	Reload bool   `json:"reload"`
}
