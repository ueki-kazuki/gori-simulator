package main

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestToName(t *testing.T) {
	type args struct {
		tags []types.Tag
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("Server01"),
					},
					{
						Key:   aws.String("Group"),
						Value: aws.String("AWS Division"),
					},
				},
			},
			want: "Server01",
		},
		{
			args: args{
				tags: []types.Tag{
					{
						Key:   aws.String("Group"),
						Value: aws.String("AWS Division"),
					},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToName(tt.args.tags); got != tt.want {
				t.Errorf("ToName() = %v, want %v", got, tt.want)
			}
		})
	}
}

// see
// https://aws.github.io/aws-sdk-go-v2/docs/unit-testing/
// https://qiita.com/tenntenn/items/eac962a49c56b2b15ee8
type MockEc2Client struct {
	reservedInstances []types.ReservedInstances
	instances         []types.Instance
}

func (m MockEc2Client) DescribeReservedInstances(ctx context.Context, params *ec2.DescribeReservedInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeReservedInstancesOutput, error) {
	for _, ri := range m.reservedInstances {
		if ri.ProductDescription == "Plan9" {
			return nil, errors.New("invalid ProductDescription")
		}
	}
	return &ec2.DescribeReservedInstancesOutput{
		ReservedInstances: m.reservedInstances,
	}, nil
}

func (m MockEc2Client) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	for _, i := range m.instances {
		if i.Platform == "Plan9" {
			return nil, errors.New("invalid Platform")
		}
	}
	return &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: m.instances,
			},
		},
	}, nil
}

func Test_getReservedInstances(t *testing.T) {
	type args struct {
		client Ec2Client
	}
	tests := []struct {
		name    string
		args    args
		want    []types.ReservedInstances
		wantErr bool
	}{
		{
			args: args{
				client: MockEc2Client{
					reservedInstances: []types.ReservedInstances{
						{
							InstanceCount:      aws.Int32(1),
							InstanceType:       "t3.medium",
							ProductDescription: types.RIProductDescription("Linux/UNIX"),
						},
					},
				},
			},
			want: []types.ReservedInstances{
				{
					InstanceCount:      aws.Int32(1),
					InstanceType:       "t3.medium",
					ProductDescription: types.RIProductDescription("Linux/UNIX"),
				},
			},
		},
		{
			args: args{
				client: MockEc2Client{
					reservedInstances: []types.ReservedInstances{
						{
							InstanceCount:      aws.Int32(9999),
							InstanceType:       "t1.dummy",
							ProductDescription: types.RIProductDescription("Plan9"),
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getReservedInstances(tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("getReservedInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getReservedInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getInstances(t *testing.T) {
	type args struct {
		client Ec2Client
	}
	tests := []struct {
		name    string
		args    args
		want    []types.Instance
		wantErr bool
	}{
		{
			args: args{
				client: MockEc2Client{
					instances: []types.Instance{
						{
							InstanceType: "t3.medium",
							Platform:     "",
						},
						{
							InstanceType: "t3.medium",
							Platform:     "windows",
						},
					},
				},
			},
			want: []types.Instance{
				{
					InstanceType: "t3.medium",
					Platform:     "Linux/UNIX",
				},
				{
					InstanceType: "t3.medium",
					Platform:     "Windows",
				},
			},
		},
		{
			args: args{
				client: MockEc2Client{
					instances: []types.Instance{
						{
							InstanceType: "t1.dummy",
							Platform:     "Plan9",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getInstances(tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}
