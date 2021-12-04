package simurator

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestSimulator_is_match(t *testing.T) {
	type fields struct {
		Instances         []*ec2.Instance
		ReservedInstances []*ec2.ReservedInstances
	}
	type args struct {
		i *ec2.Instance
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Stopped instance is NOT match",
			fields: fields{ReservedInstances: []*ec2.ReservedInstances{{}}},
			args: args{
				i: &ec2.Instance{
					State: &ec2.InstanceState{Name: aws.String("stopped")},
				},
			},
			want: false,
		},
		{
			name: "Match",
			fields: fields{ReservedInstances: []*ec2.ReservedInstances{{
				InstanceCount:      aws.Int64(1),
				InstanceType:       aws.String("t3.medium"),
				ProductDescription: aws.String("Linux/UNIX"),
			}}},
			args: args{
				i: &ec2.Instance{
					State:        &ec2.InstanceState{Name: aws.String("running")},
					InstanceType: aws.String("t3.medium"),
					Platform:     aws.String("Linux/UNIX"),
				},
			},
			want: true,
		},
		{
			name: "Num of RI is ZERO",
			fields: fields{ReservedInstances: []*ec2.ReservedInstances{{
				InstanceCount:      aws.Int64(0),
				InstanceType:       aws.String("t3.medium"),
				ProductDescription: aws.String("Linux/UNIX"),
			}}},
			args: args{
				i: &ec2.Instance{
					State:        &ec2.InstanceState{Name: aws.String("running")},
					InstanceType: aws.String("t3.medium"),
					Platform:     aws.String("Linux/UNIX"),
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
