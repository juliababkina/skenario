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

	"github.com/knative/pkg/logging"
	"go.uber.org/zap"

	"skenario/pkg/simulator"

	"github.com/knative/serving/pkg/autoscaler"
)

const (
	testNamespace = "simulator-namespace"
	testName      = "revisionService"
)

type KnativeAutoscalerConfig struct {
	TickInterval           time.Duration
	StableWindow           time.Duration
	PanicWindow            time.Duration
	PanicThreshold         float64
	ScaleToZeroGracePeriod time.Duration
	TargetConcurrency      float64
	MaxScaleUpRate         float64
}

type KnativeAutoscalerModel interface {
	Model
}

type knativeAutoscaler struct {
	env      simulator.Environment
	tickTock AutoscalerTicktockStock
}

func (kas *knativeAutoscaler) Env() simulator.Environment {
	return kas.env
}

func NewKnativeAutoscaler(env simulator.Environment, startAt time.Time, cluster ClusterModel, config KnativeAutoscalerConfig) KnativeAutoscalerModel {
	logger := logging.FromContext(env.Context())

	kpa := newKpa(logger, config)

	autoscalerEntity := simulator.NewEntity("Autoscaler", "Autoscaler")

	kas := &knativeAutoscaler{
		env:      env,
		tickTock: NewAutoscalerTicktockStock(env, autoscalerEntity, kpa, cluster),
	}

	for theTime := startAt.Add(config.TickInterval).Add(1 * time.Nanosecond); theTime.Before(env.HaltTime()); theTime = theTime.Add(config.TickInterval) {
		kas.env.AddToSchedule(simulator.NewMovement(
			"autoscaler_tick",
			theTime,
			kas.tickTock,
			kas.tickTock,
		))
	}

	return kas
}

func newKpa(logger *zap.SugaredLogger, kconfig KnativeAutoscalerConfig) *autoscaler.Autoscaler {
	deciderSpec := autoscaler.DeciderSpec{
		ServiceName:       testName,
		TickInterval:      kconfig.TickInterval,
		MaxScaleUpRate:    kconfig.MaxScaleUpRate,
		TargetConcurrency: kconfig.TargetConcurrency,
		PanicThreshold:    kconfig.PanicThreshold,
		MetricSpec: autoscaler.MetricSpec{
			StableWindow: kconfig.StableWindow,
			PanicWindow:  kconfig.PanicWindow,
		},
	}

	statsReporter, err := autoscaler.NewStatsReporter(testNamespace, testName, "config-1", "revision-1")
	if err != nil {
		logger.Fatalf("could not create stats reporter: %s", err.Error())
	}

	as, err := autoscaler.New(
		testNamespace,
		testName,
		&phonyMetricsClient{},
		&phonyCounter{},
		deciderSpec,
		statsReporter,
	)
	if err != nil {
		panic(err.Error())
	}

	return as
}

type phonyMetricsClient struct{}

func (pmc *phonyMetricsClient) StableAndPanicConcurrency(key string) (float64, float64, error) {
	return 1.0, 2.0, nil
}

type phonyCounter struct {}

func (*phonyCounter) ReadyCount() (int, error) {
	return 10, nil
}


