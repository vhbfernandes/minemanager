package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getSvc() *ec2.EC2 {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return ec2.New(sess)

}

func describeInstances() (*ec2.DescribeInstancesOutput, error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				//todo make instance tag generic
				Name: aws.String("tag:instance"),
				Values: []*string{
					aws.String("factorio"),
				},
			},
		},
	}
	return getSvc().DescribeInstances(params)
}

func getEc2InstanceId() string {
	res, err := describeInstances()
	if err != nil {
		fmt.Printf("error getting instance id: %s", err.Error())
	}
	var instid string
	for _, i := range res.Reservations[0].Instances {
		fmt.Printf("instance: %s", *i.InstanceId)
		instid = *i.InstanceId
	}
	return instid
}

func stopInstance() {
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(getEc2InstanceId()),
		},
		DryRun: aws.Bool(false),
	}
	result, err := getSvc().StopInstances(input)
	if err != nil {
		fmt.Printf("error stopping instance: %s", err.Error())
	}
	fmt.Printf("stop instance output: %s", result.String())
}

func startInstance() {
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(getEc2InstanceId()),
		},
		DryRun: aws.Bool(false),
	}
	result, err := getSvc().StartInstances(input)
	if err != nil {
		fmt.Printf("error starting instance: %s", err.Error())
	}
	fmt.Printf("start instance result: %s", result.String())
}

func getInstanceIP() string {
	res, err := describeInstances()
	if err != nil {
		fmt.Printf("error getting instance ip: %s", err.Error())
	}
	fmt.Printf("instance ip: %s", *res.Reservations[0].Instances[0].PublicIpAddress)
	return *res.Reservations[0].Instances[0].PublicIpAddress
}
