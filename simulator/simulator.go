package simurator

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type Simulator struct {
	Instances         []*ec2.Instance
	ReservedInstances []*ec2.ReservedInstances
}

type SimulatorResult struct {
	MatchInstanceResults           []*ec2.Instance
	UnmatchInstanceResults         []*ec2.Instance
	UnmatchReservedInstanceResults []*ec2.ReservedInstances
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

func (sim *Simulator) is_match(i *ec2.Instance) bool {
	if *i.State.Name != ec2.InstanceStateNameRunning {
		return false
	}
	for _, ri := range sim.ReservedInstances {
		if *ri.InstanceCount == 0 {
			continue
		}
		if *i.InstanceType == *ri.InstanceType {
			if strings.EqualFold(*i.Platform, *ri.ProductDescription) {
				*ri.InstanceCount -= 1
				return true
			}
		}
	}
	return false
}
