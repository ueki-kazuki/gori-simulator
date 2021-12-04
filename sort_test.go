package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestMultiSorter_Len(t *testing.T) {
	type fields struct {
		instances []*ec2.Instance
		less      []lessFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			fields: fields{
				instances: []*ec2.Instance{
					{
						InstanceId: aws.String("i-000000000001"),
					},
					{
						InstanceId: aws.String("i-000000000002"),
					},
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MultiSorter{
				instances: tt.fields.instances,
				less:      tt.fields.less,
			}
			if got := ms.Len(); got != tt.want {
				t.Errorf("MultiSorter.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiSorter_Swap(t *testing.T) {
	type fields struct {
		instances []*ec2.Instance
		less      []lessFunc
	}
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			fields: fields{
				instances: []*ec2.Instance{
					{
						InstanceId: aws.String("i-000000000001"),
					},
					{
						InstanceId: aws.String("i-000000000002"),
					},
				},
			},
			args: args{
				i: 0,
				j: 1,
			},
			want: []string{"i-000000000002", "i-000000000001"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MultiSorter{
				instances: tt.fields.instances,
				less:      tt.fields.less,
			}
			ms.Swap(tt.args.i, tt.args.j)
			if got := ms.instances[0].InstanceId; *got != tt.want[0] {
				t.Errorf("MultiSorter.Swap() instances[0].InstanceId = %v, want %v", *got, tt.want[0])
			}
			if got := ms.instances[1].InstanceId; *got != tt.want[1] {
				t.Errorf("MultiSorter.Swap() instances[1].InstanceId = %v, want %v", *got, tt.want[1])
			}
		})
	}
}

func TestOrderBy(t *testing.T) {
	type args struct {
		lessFunc  []lessFunc
		instances []*ec2.Instance
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			args: args{
				lessFunc: []lessFunc{
					func(p1, p2 *ec2.Instance) bool {
						return *p1.InstanceId < *p2.InstanceId
					},
				},
				instances: []*ec2.Instance{
					{
						InstanceId: aws.String("i-000000000001"),
					},
					{
						InstanceId: aws.String("i-000000000002"),
					},
				},
			},
			want: []string{"i-000000000001", "i-000000000002"},
		},
		{
			args: args{
				lessFunc: []lessFunc{
					func(p1, p2 *ec2.Instance) bool {
						return *p1.InstanceId > *p2.InstanceId
					},
				},
				instances: []*ec2.Instance{
					{
						InstanceId: aws.String("i-000000000001"),
					},
					{
						InstanceId: aws.String("i-000000000002"),
					},
				},
			},
			want: []string{"i-000000000002", "i-000000000001"},
		},
		{
			args: args{
				lessFunc: []lessFunc{
					func(p1, p2 *ec2.Instance) bool {
						return *p1.InstanceType < *p2.InstanceType
					},
					func(p1, p2 *ec2.Instance) bool {
						return *p1.InstanceId > *p2.InstanceId
					},
				},
				instances: []*ec2.Instance{
					{
						InstanceId:   aws.String("i-000000000001"),
						InstanceType: aws.String("t3.small"),
					},
					{
						InstanceId:   aws.String("i-000000000002"),
						InstanceType: aws.String("t3.small"),
					},
				},
			},
			want: []string{"i-000000000002", "i-000000000001"},
		},
		{
			args: args{
				lessFunc: []lessFunc{
					func(p1, p2 *ec2.Instance) bool {
						return *p1.InstanceType < *p2.InstanceType
					},
					func(p1, p2 *ec2.Instance) bool {
						return *p1.InstanceId < *p2.InstanceId
					},
				},
				instances: []*ec2.Instance{
					{
						InstanceId:   aws.String("i-000000000001"),
						InstanceType: aws.String("t3.small"),
					},
					{
						InstanceId:   aws.String("i-000000000002"),
						InstanceType: aws.String("t3.small"),
					},
					{
						InstanceId:   aws.String("i-000000000003"),
						InstanceType: aws.String("t1.small"),
					},
				},
			},
			want: []string{"i-000000000003", "i-000000000001", "i-000000000002"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OrderBy(tt.args.lessFunc...).Sort(tt.args.instances)
			if got := tt.args.instances[0].InstanceId; *got != tt.want[0] {
				t.Errorf("OrderBy() instances[0].InstanceId = %v, want %v", *got, tt.want[0])
			}
			if got := tt.args.instances[1].InstanceId; *got != tt.want[1] {
				t.Errorf("OrderBy() instances[1].InstanceId = %v, want %v", *got, tt.want[1])
			}
		})
	}
}
