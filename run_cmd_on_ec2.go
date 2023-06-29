package run_command_on_ec2

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
)

func RunCommand(instance_id string, cmd string) (string, error) {
	// Load AWS configuration from default environment variables, shared config, or AWS profile
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Failed to load AWS configuration:", err)
	}

	// Create an SSM client
	ssm_client := ssm.NewFromConfig(cfg)

	// Specify the EC2 instance ID where you want to run the command
	instanceID := instance_id

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
	resp, err := ssm_client.SendCommand(context.TODO(), input)
	if err != nil {
		fmt.Println("Failed to execute RunCommand:", err)
	}

	// Get the command execution details
	executionID := resp.Command.CommandId

	commandOutput, err := ssm.NewCommandExecutedWaiter(ssm_client).WaitForOutput(context.TODO(), &ssm.GetCommandInvocationInput{
		CommandId:  executionID,
		InstanceId: &instanceID,
	}, 2*time.Minute)

	if err != nil {
		fmt.Println(err)
	}
	return *commandOutput.StandardOutputContent, err
}
