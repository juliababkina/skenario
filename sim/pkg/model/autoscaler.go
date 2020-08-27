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
	"time"

	"skenario/pkg/simulator"

	"github.com/josephburnett/sk-plugin/pkg/skplug"
	"github.com/josephburnett/sk-plugin/pkg/skplug/proto"
)

type AutoscalerConfig struct {
	TickInterval time.Duration
	HpaEnabled   bool
	HpaYaml      string
}

type AutoscalerModel interface {
	Model
}

type autoscaler struct {
	env      simulator.Environment
	tickTock AutoscalerTicktockStock
}

func (a *autoscaler) Env() simulator.Environment {
	return a.env
}

type stubCluster struct{}

// TODO: actually list running pods.
func (c *stubCluster) ListPods() ([]*skplug.Pod, error) {
	return nil, nil
}

func NewAutoscaler(env simulator.Environment, startAt time.Time, cluster ClusterModel, config AutoscalerConfig) AutoscalerModel {

	autoscalerEntity := simulator.NewEntity("Autoscaler", "Autoscaler")
	if config.HpaEnabled {
		err := env.Plugin().Event(startAt.UnixNano(), proto.EventType_CREATE, &skplug.Autoscaler{
			// TODO: select type and plugin based on the scenario.
			Type: "hpa.v2beta2.autoscaling.k8s.io",
			Yaml: config.HpaYaml,
		})
		if err != nil {
			panic(err)
		}
	}

	// Create the first pods since HPA can't scale from zero.
	cm := cluster.(*clusterModel)
	for i := 0; i < int(cm.config.InitialNumberOfReplicas); i++ {
		replica := NewReplicaEntity(env, &cm.replicaSource.(*replicaSource).failedSink).(simulator.Entity)
		cm.replicasLaunching.Add(replica)
		env.AddToSchedule(simulator.NewMovement(
			"start_initial_replica",
			startAt.Add(time.Nanosecond),
			cm.replicasLaunching,
			cm.replicasActive,
			&replica,
		))
	}
	as := &autoscaler{
		env:      env,
		tickTock: NewAutoscalerTicktockStock(env, autoscalerEntity, cluster),
	}

	for theTime := startAt.Add(config.TickInterval).Add(1 * time.Nanosecond); theTime.Before(env.HaltTime()); theTime = theTime.Add(config.TickInterval) {
		for _, autoscaler := range as.tickTock.EntitiesInStock() {
			as.env.AddToSchedule(simulator.NewMovement(
				"autoscaler_tick",
				theTime,
				as.tickTock,
				as.tickTock,
				autoscaler,
			))
		}
	}

	return as
}
