package main

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
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

func getReservedInstances(s *session.Session) ([]*ec2.ReservedInstances, error) {
	svc := ec2.New(s)
	filters := []*ec2.Filter{
		&ec2.Filter{
			Name:   aws.String("state"),
			Values: []*string{aws.String("active")}},
	}
	param := ec2.DescribeReservedInstancesInput{
		Filters: filters,
	}
	result, err := svc.DescribeReservedInstances(&param)
	if err != nil {
		return nil, err
	}
	return result.ReservedInstances, nil
}

func getInstances(s *session.Session) ([]*ec2.Instance, error) {
	svc := ec2.New(s)
	param := ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(&param)
	if err != nil {
		return nil, err
	}

	instances := make([]*ec2.Instance, 0)
	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			if i.Platform == nil {
				i.Platform = aws.String("Linux/UNIX")
			}
		}
		instances = append(instances, r.Instances...)
	}
	return instances, nil
}

func ToName(tags []*ec2.Tag) string {
	for _, t := range tags {
		if *t.Key == "Name" {
			return *t.Value
		}
	}
	return ""
}

func (cli *CLI) Run(args []string) int {

	s := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	instances, err := getInstances(s)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}
	ri_instances, err := getReservedInstances(s)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}
	// fmt.Println(instances)
	// fmt.Println(ri_instances)

	sim := &simurator.Simulator{
		Instances:         instances,
		ReservedInstances: ri_instances,
	}
	results, err := sim.Simulate()
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}
	fmt.Println("=== RI covered instances ===")
	for _, i := range results.MatchInstanceResults {
		fmt.Printf("%-20s %-12s %-10s %-20s %-s\n",
			*i.InstanceId,
			*i.InstanceType,
			*i.Platform,
			ToName(i.Tags),
			*i.State.Name)
	}
	fmt.Println()

	fmt.Println("=== RI *NOT* covered instances ===")
	for _, i := range results.UnmatchInstanceResults {
		fmt.Printf("%-20s %-12s %-10s %-20s %-s\n",
			*i.InstanceId,
			*i.InstanceType,
			*i.Platform,
			ToName(i.Tags),
			*i.State.Name)
	}
	fmt.Println()

	fmt.Println("=== Purchased but not applied RI ===")
	for _, ri := range results.UnmatchReservedInstanceResults {
		fmt.Printf("%20s %-12s %-10s %-12s %3d %v\n",
			"",
			*ri.InstanceType,
			*ri.ProductDescription,
			*ri.OfferingType,
			*ri.InstanceCount,
			*ri.End)
	}

	return ExitCodeOK
}
