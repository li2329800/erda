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

package gittar

import (
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/openapi/api/apis"
)

var GITTAR_CHECK_RUN_CREATE = apis.ApiSpec{
	Path:         "/api/gittar/<org>/<repo>/check-runs",
	BackendPath:  "/<org>/<repo>/check-runs",
	Host:         "gittar.marathon.l4lb.thisdcos.directory:5566",
	Scheme:       "http",
	Method:       "POST",
	CheckLogin:   true,
	IsOpenAPI:    true,
	RequestType:  apistructs.CheckRun{},
	ResponseType: apistructs.CreateCheckRunResponse{},
	Doc:          `summary: 创建check run纪录`,
}
