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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"

	"knative-simulator/pkg/simulator"
)

type Replica interface {
	Activate()
	Deactivate()
	SendRequest(entity simulator.Entity)
}

type ReplicaEntity interface {
	simulator.Entity
	Replica
}

type replicaEntity struct {
	env                simulator.Environment
	kubernetesClient   kubernetes.Interface
	endpointsInformer  informers.EndpointsInformer
	endpointAddress    corev1.EndpointAddress
	requestsProcessing simulator.ThroughStock
	requestsComplete   simulator.SinkStock
}

func (re *replicaEntity) Activate() {
	endpoints, err := re.kubernetesClient.CoreV1().Endpoints("skenario").Get("Skenario Revision", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	endpoints.Subsets[0].Addresses = append(endpoints.Subsets[0].Addresses, re.endpointAddress)

	updatedEndpoints, err := re.kubernetesClient.CoreV1().Endpoints("skenario").Update(endpoints)
	if err != nil {
		panic(err.Error())
	}
	err = re.endpointsInformer.Informer().GetIndexer().Update(updatedEndpoints)
	if err != nil {
		panic(err.Error())
	}
}

func (re *replicaEntity) Deactivate() {
	endpoints, err := re.kubernetesClient.CoreV1().Endpoints("skenario").Get("Skenario Revision", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	addresses := endpoints.Subsets[0].Addresses

	for i, addr := range addresses {
		if addr == re.endpointAddress {
			// remove by swapping with last entry and then truncating
			addresses[i] = addresses[len(addresses)-1]
			addresses = addresses[:len(addresses)-1]

			break
		}
	}

	endpoints.Subsets[0].Addresses = addresses

	updatedEndpoints, err := re.kubernetesClient.CoreV1().Endpoints("skenario").Update(endpoints)
	if err != nil {
		panic(err.Error())
	}
	err = re.endpointsInformer.Informer().GetIndexer().Update(updatedEndpoints)
	if err != nil {
		panic(err.Error())
	}
}

func (re *replicaEntity) SendRequest(entity simulator.Entity) {
	re.requestsProcessing.Add(entity)

	re.env.AddToSchedule(simulator.NewMovement(
		"processing -> complete",
		re.env.CurrentMovementTime().Add(1*time.Second),
		re.requestsProcessing,
		re.requestsComplete,
	))
}

func (re *replicaEntity) Name() simulator.EntityName {
	return "Replica"
}

func (re *replicaEntity) Kind() simulator.EntityKind {
	return "Replica"
}

func NewReplicaEntity(env simulator.Environment, client kubernetes.Interface, endpointsInformer informers.EndpointsInformer, address string) ReplicaEntity {
	re := &replicaEntity{
		env:                env,
		kubernetesClient:   client,
		endpointsInformer:  endpointsInformer,
		requestsProcessing: simulator.NewThroughStock("RequestsProcessing", "Request"),
		requestsComplete: simulator.NewSinkStock("RequestsComplete", "Request"),
	}

	re.endpointAddress = corev1.EndpointAddress{
		IP:       address,
		Hostname: string(re.Name()),
	}

	return re
}
