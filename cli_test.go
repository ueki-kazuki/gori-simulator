package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
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
