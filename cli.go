package main

import (
	"fmt"
	"io"
	"log"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	simurator "github.com/ueki-kazuki/gori-simulator/simulator"
)

const (
	ExitCodeOK int = iota
	ExitCodeError
)

type CLI struct {
	outStream io.Writer
	errStream io.Writer
}

// for mock testing
type Ec2Client interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	DescribeReservedInstances(ctx context.Context, params *ec2.DescribeReservedInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeReservedInstancesOutput, error)
}

func getReservedInstances(client Ec2Client) ([]types.ReservedInstances, error) {
	param := ec2.DescribeReservedInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("state"),
				Values: []string{"active"},
			},
		},
	}
	result, err := client.DescribeReservedInstances(context.TODO(), &param)
	if err != nil {
		return nil, err
	}
	return result.ReservedInstances, nil
}

func getInstances(client Ec2Client) ([]types.Instance, error) {
	param := ec2.DescribeInstancesInput{}
	result, err := client.DescribeInstances(context.TODO(), &param)
	if err != nil {
		return nil, err
	}

	instances := make([]types.Instance, 0)
	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			// プラットフォームが未定義なら "Linux/UNIX" とみなす
			if i.Platform == "" {
				i.Platform = "Linux/UNIX"
			}
		}
		instances = append(instances, r.Instances...)
	}
	return instances, nil
}

func ToName(tags []types.Tag) string {
	for _, t := range tags {
		if *t.Key == "Name" {
			return *t.Value
		}
	}
	return ""
}

func (cli *CLI) Run(args []string) int {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return ExitCodeError
	}

	client := ec2.NewFromConfig(cfg)

	instances, err := getInstances(client)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}
	ri_instances, err := getReservedInstances(client)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	sim := &simurator.Simulator{
		Instances:         instances,
		ReservedInstances: ri_instances,
	}
	results, err := sim.Simulate()
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	platform := func(p1, p2 types.Instance) bool {
		if p1.Platform != p2.Platform {
			return p1.Platform == ""
		} else {
			return true
		}
	}

	instancetype := func(p1, p2 types.Instance) bool {
		return p1.InstanceType < p2.InstanceType
	}

	name := func(p1, p2 types.Instance) bool {
		return ToName(p1.Tags) < ToName(p2.Tags)
	}

	state := func(p1, p2 types.Instance) bool {
		return *p1.State.Code < *p2.State.Code
	}

	fmt.Println("=== RI covered instances ===")
	OrderBy(state, platform, instancetype, name).Sort(results.MatchInstanceResults)
	for _, i := range results.MatchInstanceResults {
		fmt.Printf("%-20s %-12s %-10s %-20s %-s\n",
			*i.InstanceId,
			i.InstanceType,
			i.Platform,
			ToName(i.Tags),
			i.State.Name)
	}
	fmt.Println()

	fmt.Println("=== RI *NOT* covered instances ===")
	OrderBy(state, platform, instancetype, name).Sort(results.UnmatchInstanceResults)
	for _, i := range results.UnmatchInstanceResults {
		fmt.Printf("%-20s %-12s %-10s %-20s %-s\n",
			*i.InstanceId,
			i.InstanceType,
			i.Platform,
			ToName(i.Tags),
			i.State.Name)
	}
	fmt.Println()

	fmt.Println("=== Purchased but not applied RI ===")
	for _, ri := range results.UnmatchReservedInstanceResults {
		fmt.Printf("%20s %-12s %-10s %-12s %3d %v\n",
			"",
			ri.InstanceType,
			ri.ProductDescription,
			ri.OfferingType,
			*ri.InstanceCount,
			ri.End)
	}

	return ExitCodeOK
}
