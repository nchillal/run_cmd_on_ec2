package run_cmd_on_ec2

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
)

// RunCommand runs command given on AWS EC2
func RunCommand(awsProfile string, awsRegion string, instanceId string, cmd string) (string, error) {
	// Load AWS configuration from default environment variables, shared config, or AWS profile
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(awsProfile),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		fmt.Println("Failed to load AWS configuration: ", err)
	}

	// Create an SSM client
	ssmClient := ssm.NewFromConfig(cfg)

	// Specify the EC2 instance ID where you want to run the command
	instanceID := instanceId

	// Specify the command you want to run on the EC2 instance
	command := cmd

	// Build the RunCommand input parameters
	input := &ssm.SendCommandInput{
		InstanceIds:  []string{instanceID},
		DocumentName: aws.String("AWS-RunShellScript"),
		Parameters: map[string][]string{
			"commands": {command},
		},
	}

	// Execute the RunCommand API
	resp, err := ssmClient.SendCommand(context.TODO(), input)
	if err != nil {
		fmt.Println("Failed to execute RunCommand:", err)
	}

	// Get the command execution details
	executionID := resp.Command.CommandId

	commandOutput, err := ssm.NewCommandExecutedWaiter(ssmClient).WaitForOutput(context.TODO(), &ssm.GetCommandInvocationInput{
		CommandId:  executionID,
		InstanceId: &instanceID,
	}, 2*time.Minute)

	if err != nil {
		fmt.Println(err)
	}
	return *commandOutput.StandardOutputContent, err
}
