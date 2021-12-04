package simurator

import (
	"reflect"
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

func TestSimulator_Simulate(t *testing.T) {
	type fields struct {
		Instances         []types.Instance
		ReservedInstances []types.ReservedInstances
	}
	tests := []struct {
		name    string
		fields  fields
		want    SimulatorResult
		wantErr bool
	}{
		{
			fields: fields{
				Instances: []types.Instance{
					{
						State:        &types.InstanceState{Name: types.InstanceStateNameStopped},
						InstanceType: "t3.medium",
						Platform:     types.PlatformValues("Linux/UNIX"),
					},
					{
						State:        &types.InstanceState{Name: types.InstanceStateNameRunning},
						InstanceType: "t3.medium",
						Platform:     types.PlatformValues("Linux/UNIX"),
					},
					{
						State:        &types.InstanceState{Name: types.InstanceStateNameRunning},
						InstanceType: "t3.medium",
						Platform:     types.PlatformValuesWindows,
					},
				},
				ReservedInstances: []types.ReservedInstances{
					{
						InstanceCount:      aws.Int32(1),
						InstanceType:       "t3.medium",
						ProductDescription: types.RIProductDescription("Linux/UNIX"),
					},
					{
						InstanceCount:      aws.Int32(1),
						InstanceType:       "c5.xlarge",
						ProductDescription: types.RIProductDescription("Linux/UNIX"),
					},
				},
			},
			wantErr: false,
			want: SimulatorResult{
				MatchInstanceResults: []types.Instance{
					{
						State:        &types.InstanceState{Name: types.InstanceStateNameRunning},
						InstanceType: "t3.medium",
						Platform:     types.PlatformValues("Linux/UNIX"),
					},
				},
				UnmatchInstanceResults: []types.Instance{
					{
						State:        &types.InstanceState{Name: types.InstanceStateNameStopped},
						InstanceType: "t3.medium",
						Platform:     types.PlatformValues("Linux/UNIX"),
					},
					{
						State:        &types.InstanceState{Name: types.InstanceStateNameRunning},
						InstanceType: "t3.medium",
						Platform:     types.PlatformValuesWindows,
					},
				},
				UnmatchReservedInstanceResults: []types.ReservedInstances{
					{
						InstanceCount:      aws.Int32(1),
						InstanceType:       "c5.xlarge",
						ProductDescription: types.RIProductDescription("Linux/UNIX"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim := &Simulator{
				Instances:         tt.fields.Instances,
				ReservedInstances: tt.fields.ReservedInstances,
			}
			got, err := sim.Simulate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Simulator.Simulate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Simulator.Simulate() = %v, want %v", got, tt.want)
			}
		})
	}
}
