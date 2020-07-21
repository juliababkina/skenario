/*
 * Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.
 *
 * This program and the accompanying materials are made available under the terms
 * of the Apache License, Version 2.0 (the "License”); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at:
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package model

import (
	"skenario/pkg/simulator"
)

type TrafficSource interface {
	simulator.SourceStock
}

type trafficSource struct {
	env             simulator.Environment
	requestsRouting RequestsRoutingStock
	requestConfig   RequestConfig
}

func (ts *trafficSource) Name() simulator.StockName {
	return "TrafficSource"
}

func (ts *trafficSource) KindStocked() simulator.EntityKind {
	return "Request"
}

func (ts *trafficSource) Count() uint64 {
	return 0
}

func (ts *trafficSource) EntitiesInStock() []*simulator.Entity {
	return []*simulator.Entity{}
}

func (ts *trafficSource) Remove(entity *simulator.Entity) simulator.Entity {
	return NewRequestEntity(ts.env, ts.requestsRouting, ts.requestConfig)
}

func NewTrafficSource(env simulator.Environment, requestsRouting RequestsRoutingStock, requestConfig RequestConfig) TrafficSource {
	return &trafficSource{
		env:             env,
		requestsRouting: requestsRouting,
		requestConfig:   requestConfig,
	}
}
