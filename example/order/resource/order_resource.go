/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package resource

import (
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"github.com/go-chassis/go-chassis/server/restful"
	"net/http"
)

type OrderResource struct {
}

func (r *OrderResource) Get(context *restful.Context) {
	orderID := context.ReadPathParameter("order_id")
	invoker := core.NewRestInvoker()
	req, err := rest.NewRequest("GET", "http://restaurant/v1/restaurant/"+orderID, nil)
	if err != nil {
		context.Write([]byte(err.Error()))
		context.WriteHeader(http.StatusInternalServerError)
	}
	resp, err := invoker.ContextDo(context.Ctx, req)
	if err != nil {
		context.Write([]byte(err.Error()))
		context.WriteHeader(http.StatusInternalServerError)
		return
	}
	restaurant := httputil.ReadBody(resp)
	context.Write(restaurant)
}

func (r *OrderResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/v1/order/{order_id}", ResourceFuncName: "Get"},
	}
}
