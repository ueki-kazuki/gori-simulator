package simurator

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Simulator struct {
	Instances         []types.Instance
	ReservedInstances []types.ReservedInstances
}

type SimulatorResult struct {
	MatchInstanceResults           []types.Instance
	UnmatchInstanceResults         []types.Instance
	UnmatchReservedInstanceResults []types.ReservedInstances
}

func (sim *Simulator) Simulate() (SimulatorResult, error) {
	results := SimulatorResult{}

	for _, i := range sim.Instances {
		if sim.is_match(i) {
			results.MatchInstanceResults = append(results.MatchInstanceResults, i)
		} else {
			results.UnmatchInstanceResults = append(results.UnmatchInstanceResults, i)
		}
	}

	for _, ri := range sim.ReservedInstances {
		if *ri.InstanceCount != 0 {
			results.UnmatchReservedInstanceResults = append(results.UnmatchReservedInstanceResults, ri)
		}
	}

	return results, nil
}

func (sim *Simulator) is_match(i types.Instance) bool {
	if i.State.Name != types.InstanceStateNameRunning {
		return false
	}
	for _, ri := range sim.ReservedInstances {
		if *ri.InstanceCount == 0 {
			continue
		}
		if i.InstanceType == ri.InstanceType {
			if string(i.Platform) == string(ri.ProductDescription) {
				*ri.InstanceCount -= 1
				return true
			}
		}
	}
	return false
}
