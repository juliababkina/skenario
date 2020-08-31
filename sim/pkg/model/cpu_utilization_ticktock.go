package model

import "skenario/pkg/simulator"

type CpuUtilizationTicktockStock interface {
	simulator.ThroughStock
}

type cpuUtilizationTicktockStock struct {
	env                  simulator.Environment
	replicasActive       *ReplicasActiveStock
	cpuUtilizationEntity simulator.Entity
}

func (cuts *cpuUtilizationTicktockStock) Name() simulator.StockName {
	return "Cpu utilization Ticktock"
}

func (cuts *cpuUtilizationTicktockStock) KindStocked() simulator.EntityKind {
	return "Cpu utilization"
}

func (cuts *cpuUtilizationTicktockStock) Count() uint64 {
	return 1
}

func (cuts *cpuUtilizationTicktockStock) EntitiesInStock() []*simulator.Entity {
	return []*simulator.Entity{&cuts.cpuUtilizationEntity}
}

func (cuts *cpuUtilizationTicktockStock) Remove(entity *simulator.Entity) simulator.Entity {
	return cuts.cpuUtilizationEntity
}

func (cuts *cpuUtilizationTicktockStock) Add(entity simulator.Entity) error {
	countActiveReplicas := 0.0
	totalCPUUtilization := 0.0 // total cpuUtilization for all active replicas in percentage

	for _, en := range (*cuts.replicasActive).EntitiesInStock() {
		replica := (*en).(*replicaEntity)
		totalCPUUtilization += replica.occupiedCPUCapacityMillisPerSecond * 100 / replica.totalCPUCapacityMillisPerSecond
		countActiveReplicas++
	}
	if countActiveReplicas > 0 {
		averageCPUUtilizationPerReplica := simulator.CPUUtilization{CPUUtilization: totalCPUUtilization / countActiveReplicas,
			CalculatedAt: cuts.env.CurrentMovementTime()}
		cuts.env.AppendCPUUtilization(&averageCPUUtilizationPerReplica)
	}
	return nil
}

func NewCpuUtilizationTicktockStock(env simulator.Environment, cpuUtilizationEntity simulator.Entity, replicasActive *ReplicasActiveStock) CpuUtilizationTicktockStock {
	return &cpuUtilizationTicktockStock{
		env:                  env,
		replicasActive:       replicasActive,
		cpuUtilizationEntity: cpuUtilizationEntity,
	}
}
