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

package executeTaskTable

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
	protocol "github.com/erda-project/erda/modules/openapi/component-protocol"
	"github.com/erda-project/erda/modules/openapi/component-protocol/pkg/component_key"
)

type ExecuteTaskTable struct {
	CtxBdl     protocol.ContextBundle
	Type       string                 `json:"type"`
	Props      map[string]interface{} `json:"props"`
	Operations map[string]interface{} `json:"operations"`
	State      State                  `json:"state"`
	Data       map[string]interface{} `json:"data"`
}

type State struct {
	Total      int64  `json:"total"`
	PageSize   int64  `json:"pageSize"`
	PageNo     int64  `json:"pageNo"`
	PipelineID uint64 `json:"pipelineId"`
	StepID     uint64 `json:"stepId"`
	Name       string `json:"name"`
	Unfold     bool   `json:"unfold"`
}

type operationData struct {
	Meta meta `json:"meta"`
}

type meta struct {
	Target RowData `json:"target"`
}

const (
	DefaultPageSize = 1000
	DefaultPageNo   = 1
)

type Operation struct {
	Key           string      `json:"key"`
	Reload        bool        `json:"reload"`
	FillMeta      string      `json:"fillMeta"`
	Meta          interface{} `json:"meta"`
	ClickableKeys interface{} `json:"clickableKeys"`
}

type props struct {
	Key            string                   `json:"key"`
	Label          string                   `json:"label"`
	Component      string                   `json:"component"`
	Required       bool                     `json:"required"`
	Rules          []map[string]interface{} `json:"rules"`
	ComponentProps map[string]interface{}   `json:"componentProps,omitempty"`
}

type inParams struct {
	ProjectID int64 `json:"projectId"`
}

type columns struct {
	Title     string `json:"title"`
	DataIndex string `json:"dataIndex"`
	Width     int    `json:"width,omitempty"`
	Ellipsis  bool   `json:"ellipsis"`
	Fixed     string `json:"fixed"`
}

type dataOperation struct {
	Key         string                 `json:"key"`
	Reload      bool                   `json:"reload"`
	Text        string                 `json:"text"`
	Disabled    bool                   `json:"disabled"`
	DisabledTip string                 `json:"disabledTip,omitempty"`
	Confirm     string                 `json:"confirm,omitempty"`
	Meta        interface{}            `json:"meta,omitempty"`
	Command     map[string]interface{} `json:"command,omitempty"`
}

type RowData struct {
	Name              string `json:"name"`
	SnippetPipelineID uint64 `json:"snippetPipelineID"`
}

type AutoTestRunStep struct {
	ApiSpec     map[string]interface{} `json:"apiSpec"`
	WaitTime    int64                  `json:"waitTime"`
	Commands    []string               `json:"commands"`
	Image       string                 `json:"image"`
	WaitTimeSec int64                  `json:"waitTimeSec"`
}

func (a *ExecuteTaskTable) Import(c *apistructs.Component) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, a); err != nil {
		return err
	}
	return nil
}

func (a *ExecuteTaskTable) Export(c *apistructs.Component, gs *apistructs.GlobalStateData) error {
	// set component data
	b, err := json.Marshal(a)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, c); err != nil {
		return err
	}
	return nil
}

func (a *ExecuteTaskTable) Render(ctx context.Context, c *apistructs.Component, scenario apistructs.ComponentProtocolScenario, event apistructs.ComponentEvent, gs *apistructs.GlobalStateData) error {
	// import component data
	if err := a.Import(c); err != nil {
		logrus.Errorf("import component failed, err:%v", err)
		return err
	}

	a.CtxBdl = ctx.Value(protocol.GlobalInnerKeyCtxBundle.String()).(protocol.ContextBundle)

	if a.CtxBdl.InParams == nil {
		return fmt.Errorf("params is empty")
	}
	inParamsBytes, err := json.Marshal(a.CtxBdl.InParams)
	if err != nil {
		return fmt.Errorf("failed to marshal inParams, inParams:%+v, err:%v", a.CtxBdl.InParams, err)
	}
	var inParams inParams
	if err := json.Unmarshal(inParamsBytes, &inParams); err != nil {
		return err
	}

	defer func() {
		fail := a.marshal(c)
		if err == nil && fail != nil {
			err = fail
		}
		// export rendered component data
		c.Operations = a.Operations
		c.Props = getProps()
	}()

	// listen on operation
	switch event.Operation {
	case apistructs.ExecuteChangePageNoOperationKey, apistructs.RenderingOperation, apistructs.InitializeOperation:
		a.State.Unfold = false
		if err := a.handlerListOperation(a.CtxBdl, c, inParams, event); err != nil {
			return err
		}
	case apistructs.ExecuteClickRowNoOperationKey:
		if err := a.handlerClickRowOperation(a.CtxBdl, c, inParams, event); err != nil {
			return err
		}
	}
	return nil
}

func getOperations(clickableKeys []uint64) map[string]interface{} {
	return map[string]interface{}{
		"changePageNo": Operation{
			Key:    "changePageNo",
			Reload: true,
		},
		"clickRow": Operation{
			Key:           "clickRow",
			Reload:        true,
			FillMeta:      "target",
			Meta:          nil,
			ClickableKeys: clickableKeys,
		},
	}
}

func getProps() map[string]interface{} {
	return map[string]interface{}{
		"rowKey": "key",
		"scroll": map[string]interface{}{"x": 1200},
		"columns": []columns{
			{
				Title:     "步骤名称",
				DataIndex: "name",
				Width:     200,
				Ellipsis:  true,
				Fixed:     "left",
			},
			{
				Title:     "步骤类型",
				DataIndex: "type",
				Width:     85,
				Ellipsis:  true,
			},
			{
				Title:     "步骤",
				DataIndex: "step",
				Width:     85,
				Ellipsis:  true,
			},
			{
				Title:     "子任务数",
				DataIndex: "tasksNum",
				Width:     85,
				Ellipsis:  true,
			},
			{
				Title:     "接口路径",
				DataIndex: "path",
			},
			{
				Title:     "状态",
				DataIndex: "status",
				Width:     120,
				Ellipsis:  true,
			},
			{
				Title:     "操作",
				DataIndex: "operate",
				Width:     120,
				Ellipsis:  true,
				Fixed:     "right",
			},
		},
	}
}

func transformStepType(str apistructs.StepAPIType) string {
	switch str {
	case apistructs.StepTypeWait:
		return "等待"
	case apistructs.StepTypeAPI:
		return "接口"
	case apistructs.StepTypeScene:
		return "场景"
	case apistructs.StepTypeConfigSheet:
		return "配置单"
	case apistructs.StepTypeCustomScript:
		return "自定义"
	}
	return string(str)
}

func getStatus(req apistructs.PipelineStatus) map[string]interface{} {
	res := map[string]interface{}{"renderType": "textWithBadge", "value": req.ToDesc()}
	if req.IsSuccessStatus() {
		res["status"] = "success"
	}
	if req.IsFailedStatus() {
		res["status"] = "error"
	}
	if req.IsReconcilerRunningStatus() {
		res["status"] = "processing"
	}
	if req.IsBeforePressRunButton() {
		res["status"] = "default"
	}
	return res
}

func (a *ExecuteTaskTable) setData(pipeline *apistructs.PipelineDetailDTO) error {
	lists := []map[string]interface{}{}
	clickableKeys := []uint64{}
	num := (a.State.PageNo - 1) * (a.State.PageSize)
	ret := a.State.PageSize
	a.State.Total = 0
	stepIdx := 1
	for _, each := range pipeline.PipelineStages {

		a.State.Total += int64(len(each.PipelineTasks))
		if ret == 0 {
			continue
		}

		if int64(len(each.PipelineTasks)) <= num {
			num -= int64(len(each.PipelineTasks))
			continue
		}

		for _, task := range each.PipelineTasks {

			if num > 0 {
				num--
				continue
			}

			var item map[string]interface{}

			// not autotest task
			// snippet acton remove operation add taskNum
			if task.Labels == nil || len(task.Labels) == 0 {
				operations := map[string]interface{}{}
				var taskNum interface{}
				if task.Type != apistructs.ActionTypeSnippet {
					operations = map[string]interface{}{
						"checkDetail": dataOperation{
							Key:         "checkDetail",
							Text:        "查看结果",
							Reload:      false,
							Meta:        task.Result,
							DisabledTip: "禁用接口无法查看结果",
							Disabled:    task.Status.IsDisabledStatus(),
						},
						"checkLog": dataOperation{
							Key:    "checkLog",
							Reload: false,
							Text:   "日志",
							Meta: map[string]interface{}{
								"logId":      task.Extra.UUID,
								"pipelineId": a.State.PipelineID,
								"nodeId":     task.ID,
							},
							DisabledTip: "禁用接口无法查看日志",
							Disabled:    task.Status.IsDisabledStatus(),
						},
					}
					taskNum = "-"
				} else {
					clickableKeys = append(clickableKeys, task.ID)
					if task.SnippetPipelineDetail != nil {
						taskNum = task.SnippetPipelineDetail.DirectSnippetTasksNum
					}
				}
				item = map[string]interface{}{
					"id":                task.ID,
					"key":               component_key.GetKey(task.ID),
					"snippetPipelineID": task.SnippetPipelineID,
					"operate": map[string]interface{}{
						"renderType": "tableOperation",
						"operations": operations,
					},
					"tasksNum": taskNum,
					"name":     task.Name,
					"status":   getStatus(task.Status),
					"type":     transformStepType(apistructs.StepAPIType(task.Type)),
					"path":     "",
					"step":     stepIdx,
				}
			} else {
				// autotest task
				res := apistructs.AutoTestSceneStep{}
				value := AutoTestRunStep{
					ApiSpec: map[string]interface{}{},
				}
				if _, ok := task.Labels[apistructs.AutotestSceneStep]; ok {
					resByte, err := base64.StdEncoding.DecodeString(task.Labels[apistructs.AutotestSceneStep])
					if err != nil {
						logrus.Error("error to decode ", apistructs.AutotestSceneStep, ", err: ", err)
						return err
					}
					if err := json.Unmarshal(resByte, &res); err != nil {
						return err
					}
					if res.Type == apistructs.StepTypeAPI || res.Type == apistructs.StepTypeWait || res.Type == apistructs.StepTypeCustomScript {
						err := json.Unmarshal([]byte(res.Value), &value)
						if err != nil {
							return err
						}
					}
					if res.Type == apistructs.StepTypeWait {
						if value.WaitTime > 0 {
							value.WaitTimeSec = value.WaitTime
						}
						res.Name = transformStepType(res.Type) + strconv.FormatInt(value.WaitTimeSec, 10) + "s"
					}
				} else {
					res.Name = task.Name
					res.Type = apistructs.StepAPIType(task.Type)
				}

				// api or wait add operations
				operations := map[string]interface{}{}
				if res.Type == apistructs.StepTypeAPI || res.Type == apistructs.StepTypeWait || res.Type == apistructs.StepTypeCustomScript {
					operations = map[string]interface{}{
						"checkDetail": dataOperation{
							Key:         "checkDetail",
							Text:        "查看结果",
							Reload:      false,
							Meta:        task.Result,
							DisabledTip: "禁用接口无法查看结果",
							Disabled:    task.Status.IsDisabledStatus(),
						},
						"checkLog": dataOperation{
							Key:    "checkLog",
							Reload: false,
							Text:   "日志",
							Meta: map[string]interface{}{
								"logId":      task.Extra.UUID,
								"pipelineId": a.State.PipelineID,
								"nodeId":     task.ID,
							},
							DisabledTip: "禁用接口无法查看日志",
							Disabled:    task.Status.IsDisabledStatus(),
						},
					}
				}

				path := value.ApiSpec["url"]
				if path == nil {
					path = ""
				}
				item = map[string]interface{}{
					"id":                task.ID,
					"key":               component_key.GetKey(task.ID),
					"snippetPipelineID": task.SnippetPipelineID,
					"operate": map[string]interface{}{
						"renderType": "tableOperation",
						"operations": operations,
					},
					"tasksNum": "-",
					"name":     res.Name,
					"status":   getStatus(task.Status),
					"type":     transformStepType(res.Type),
					"path":     path,
					"step":     stepIdx,
				}

				// scene or configSheet add other info
				if task.SnippetPipelineID != nil &&
					(res.Type == apistructs.StepTypeScene || res.Type == apistructs.StepTypeConfigSheet) {
					clickableKeys = append(clickableKeys, task.ID)
					if task.SnippetPipelineDetail != nil {
						item["tasksNum"] = task.SnippetPipelineDetail.DirectSnippetTasksNum
					}
				}
			}

			lists = append(lists, item)

			ret--
			if ret == 0 {
				break
			}
		}
		stepIdx++
	}

	if a.State.Total <= (a.State.PageNo-1)*(a.State.PageSize) && a.State.Total > 0 {
		a.State.PageNo = DefaultPageNo
		return a.setData(pipeline)
	}
	a.Data["list"] = lists
	a.Operations = getOperations(clickableKeys)
	return nil
}

func (a *ExecuteTaskTable) marshal(c *apistructs.Component) error {
	stateValue, err := json.Marshal(a.State)
	if err != nil {
		return err
	}
	var state map[string]interface{}
	err = json.Unmarshal(stateValue, &state)
	if err != nil {
		return err
	}

	propValue, err := json.Marshal(a.Props)
	if err != nil {
		return err
	}
	var props interface{}
	err = json.Unmarshal(propValue, &props)
	if err != nil {
		return err
	}

	c.Props = props
	c.State = state
	c.Type = a.Type
	return nil
}

func (e *ExecuteTaskTable) handlerListOperation(bdl protocol.ContextBundle, c *apistructs.Component, inParams inParams, event apistructs.ComponentEvent) error {

	e.State.PageNo = DefaultPageNo
	e.State.PageSize = DefaultPageSize

	if e.State.PipelineID == 0 {
		c.Data = map[string]interface{}{}
		return nil
	}
	list, err := bdl.Bdl.GetPipeline(e.State.PipelineID)
	if err != nil {
		return err
	}
	err = e.setData(list)
	if err != nil {
		return err
	}
	c.Data = e.Data
	return nil
}

func (e *ExecuteTaskTable) handlerClickRowOperation(bdl protocol.ContextBundle, c *apistructs.Component, inParams inParams, event apistructs.ComponentEvent) error {

	res := operationData{}
	b, err := json.Marshal(event.OperationData)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	e.State.Name = res.Meta.Target.Name
	e.State.PipelineID = res.Meta.Target.SnippetPipelineID
	e.State.Unfold = true
	if res.Meta.Target.SnippetPipelineID == 0 {
		return nil
	}
	if err := e.handlerListOperation(e.CtxBdl, c, inParams, event); err != nil {
		return err
	}
	return nil
}

func RenderCreator() protocol.CompRender {
	return &ExecuteTaskTable{
		CtxBdl:     protocol.ContextBundle{},
		Props:      map[string]interface{}{},
		Operations: map[string]interface{}{},
		State:      State{},
		Data:       map[string]interface{}{},
	}
}
