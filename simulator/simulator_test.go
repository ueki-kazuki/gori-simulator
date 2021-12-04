package simurator

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestSimulator_is_match(t *testing.T) {
	type fields struct {
		Instances         []types.Instance
		ReservedInstances []types.ReservedInstances
	}
	type args struct {
		i types.Instance
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Stopped instance is NOT match",
			fields: fields{ReservedInstances: []types.ReservedInstances{{}}},
			args: args{
				i: types.Instance{
					State: &types.InstanceState{Name: types.InstanceStateNameStopped},
				},
			},
			want: false,
		},
		{
			name: "Match",
			fields: fields{ReservedInstances: []types.ReservedInstances{{
				InstanceCount:      aws.Int32(1),
				InstanceType:       "t3.medium",
				ProductDescription: types.RIProductDescription("Linux/UNIX"),
			}}},
			args: args{
				i: types.Instance{
					State:        &types.InstanceState{Name: types.InstanceStateNameRunning},
					InstanceType: "t3.medium",
					Platform:     types.PlatformValues("Linux/UNIX"),
				},
			},
			want: true,
		},
		{
			name: "Num of RI is ZERO",
			fields: fields{ReservedInstances: []types.ReservedInstances{{
				InstanceCount:      aws.Int32(0),
				InstanceType:       "t3.medium",
				ProductDescription: types.RIProductDescription("Linux/UNIX"),
			}}},
			args: args{
				i: types.Instance{
					State:        &types.InstanceState{Name: types.InstanceStateNameRunning},
					InstanceType: "t3.medium",
					Platform:     types.PlatformValues("Linux/UNIX"),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim := &Simulator{
				Instances:         tt.fields.Instances,
				ReservedInstances: tt.fields.ReservedInstances,
			}
			if got := sim.is_match(tt.args.i); got != tt.want {
				t.Errorf("Simulator.is_match() = %v, want %v", got, tt.want)
			}
		})
	}
}
