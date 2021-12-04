package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestToName(t *testing.T) {
	type args struct {
		tags []*ec2.Tag
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				tags: []*ec2.Tag{
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
				tags: []*ec2.Tag{
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
