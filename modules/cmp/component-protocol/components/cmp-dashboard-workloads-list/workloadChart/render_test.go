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

package workloadChart

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/i18n"
	"github.com/erda-project/erda/modules/cmp/component-protocol/cputil"
)

func TestComponentWorkloadChart_GenComponentState(t *testing.T) {
	component := &cptype.Component{
		State: map[string]interface{}{
			"values": Values{
				DeploymentsCount: Count{
					Active: 1,
					Error:  1,
				},
				DaemonSetCount: Count{
					Active: 1,
					Error:  1,
				},
				StatefulSetCount: Count{
					Active: 1,
					Error:  1,
				},
				JobCount: Count{
					Active:    1,
					Succeeded: 1,
					Failed:    1,
				},
				CronJobCount: Count{
					Active: 1,
				},
			},
		},
	}
	src, err := json.Marshal(component.State)
	if err != nil {
		t.Errorf("test failed, %v", err)
	}

	f := &ComponentWorkloadChart{}
	if err := f.GenComponentState(component); err != nil {
		t.Errorf("test failed, %v", err)
	}

	dst, err := json.Marshal(f.State)
	if err != nil {
		t.Errorf("test failed, %v", err)
	}

	if string(src) != string(dst) {
		t.Error("test failed, generate result is unexpected")
	}
}

type MockTran struct {
	i18n.Translator
}

func (m *MockTran) Text(lang i18n.LanguageCodes, key string) string {
	return ""
}

func (m *MockTran) Sprintf(lang i18n.LanguageCodes, key string, args ...interface{}) string {
	return ""
}

func TestComponentWorkloadChart_SetComponentValue(t *testing.T) {
	sdk := cptype.SDK{Tran: &MockTran{}}
	ctx := context.WithValue(context.Background(), cptype.GlobalInnerKeyCtxSDK, &sdk)
	component := &ComponentWorkloadChart{}
	if err := component.SetComponentValue(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestComponentWorkloadChart_Transfer(t *testing.T) {
	tmp := 1
	component := &ComponentWorkloadChart{
		State: State{
			Values: Values{
				DeploymentsCount: Count{
					Active:    1,
					Error:     1,
					Succeeded: 1,
					Failed:    1,
				},
				DaemonSetCount: Count{
					Active:    2,
					Error:     2,
					Succeeded: 2,
					Failed:    2,
				},
				StatefulSetCount: Count{
					Active:    3,
					Error:     3,
					Succeeded: 3,
					Failed:    3,
				},
				JobCount: Count{
					Active:    4,
					Error:     4,
					Succeeded: 4,
					Failed:    4,
				},
				CronJobCount: Count{
					Active:    5,
					Error:     5,
					Succeeded: 5,
					Failed:    5,
				},
			},
		},
		Props: Props{
			Option: Option{
				Tooltip: Tooltip{
					Trigger: "test",
					AxisPointer: AxisPointer{
						Type: "testType",
					},
				},
				Color: []string{"test"},
				Legend: Legend{
					Data: []string{"test"},
				},
				XAxis: Axis{
					Type: "testType",
					Data: []string{"test"},
				},
				YAxis: Axis{
					Type: "testType",
					Data: []string{"test"},
				},
				Series: []Series{
					{
						Name:     "testName",
						Type:     "testType",
						Stack:    "testStack",
						BarWidth: "teatBarWidth",
						Data: []*int{
							&tmp,
						},
					},
				},
			},
		},
	}
	c := &cptype.Component{}
	component.Transfer(c)
	ok, err := cputil.IsJsonEqual(c.State, component.State)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("test failed, json is not equal")
	}
}
