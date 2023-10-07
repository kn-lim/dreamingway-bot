package pixelmon

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/kn-lim/dreamingway-bot/internal/mcstatus"
)

func GetStatus(url string) (bool, int, error) {
	return mcstatus.GetMCStatus(url)
}

func StartInstance(instanceID string) error {
	// Setup AWS session
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	// Get EC2 instance
	client, instance, err := getInstance(cfg, instanceID)
	if err != nil {
		return err
	}

	// Start EC2 instance
	if instance.State.Name == types.InstanceStateNameStopped {
		input := &ec2.StartInstancesInput{
			InstanceIds: []string{*instance.InstanceId},
		}

		_, err := client.StartInstances(context.TODO(), input)
		if err != nil {
			log.Printf("Error! Failed to start instance: %v", err)
			return err
		}
	}

	return nil
}

func StartService(instanceID string, zoneID string, url string) error {
	// Setup AWS session
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	// Wait till EC2 instance is running
	var instance types.Instance
	for {
		// Get EC2 instance
		_, i, err := getInstance(cfg, instanceID)
		if err != nil {
			return err
		}

		if i.State.Name == types.InstanceStateNameRunning {
			instance = i
			break
		}

		log.Println("Waiting for instance...")
		time.Sleep(delay * time.Second)
	}

	log.Println("Instance is running")

	// Create Route53 DNS entry
	if err := createInstanceDNSEntry(cfg, instance, zoneID, url); err != nil {
		return err
	}

	// Delay to ensure command can be sent to EC2 instance
	time.Sleep(delay * time.Second)

	// Send start command to EC2 instance
	client := ssm.NewFromConfig(cfg)
	documentName := "AWS-RunShellScript"
	params := map[string][]string{
		"commands": {"cd /opt/pixelmon/ && tmux new-session -d -s minecraft './start.sh'"},
	}
	input := &ssm.SendCommandInput{
		InstanceIds:  []string{instanceID},
		DocumentName: &documentName,
		Parameters:   params,
	}
	_, err = client.SendCommand(context.TODO(), input)
	if err != nil {
		log.Printf("Error! Couldn't send command: %v", err)
		return err
	}

	log.Println("Sent command to EC2 instance")

	// Check if service is online
	for {
		status, _, err := GetStatus(url)
		if err != nil {
			return err
		}

		if status {
			break
		}

		log.Println("Waiting for service...")
		time.Sleep(delay * time.Second)
	}

	log.Println("Service is running")

	return nil
}

func StopInstance(instanceID string) error {
	// Setup AWS session
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	// Get EC2 instance
	client, instance, err := getInstance(cfg, instanceID)
	if err != nil {
		return err
	}

	// Stop EC2 instance
	if instance.State.Name == types.InstanceStateNameRunning {
		input := &ec2.StopInstancesInput{
			InstanceIds: []string{*instance.InstanceId},
		}

		_, err := client.StopInstances(context.TODO(), input)
		if err != nil {
			log.Printf("Error! Failed to stop instance: %v", err)
			return err
		}
	}

	log.Println("Instance is stopped")

	return nil
}

func StopService(instanceID string, zoneID string, url string, rcon string) error {
	// Setup AWS session
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	// Get EC2 instance
	_, instance, err := getInstance(cfg, instanceID)
	if err != nil {
		return err
	}

	// Remove Route53 DNS entry
	if err := removeInstanceDNSEntry(cfg, instance, zoneID, url); err != nil {
		return err
	}

	// Send stop command to EC2 instance
	client := ssm.NewFromConfig(cfg)
	documentName := "AWS-RunShellScript"
	params := map[string][]string{
		"commands": {"mcrcon -H localhost -p " + rcon + " \"stop\""},
	}
	input := &ssm.SendCommandInput{
		InstanceIds:  []string{instanceID},
		DocumentName: &documentName,
		Parameters:   params,
	}
	_, err = client.SendCommand(context.TODO(), input)
	if err != nil {
		log.Printf("Error! Couldn't send command: %v", err)
		return err
	}

	log.Println("Sent command to EC2 instance")

	// Check if service is offline
	for {
		log.Println("Waiting for service...")
		time.Sleep(delay * time.Second)

		status, _, err := GetStatus(url)
		if err != nil {
			return err
		}

		if !status {
			break
		}
	}

	log.Println("Service is stopped")

	return nil
}

func SayMessage(instanceID string, rcon string, username string, message string) error {
	// Setup AWS session
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	// Send say command to Pixelmon EC2 instance
	client := ssm.NewFromConfig(cfg)
	documentName := "AWS-RunShellScript"
	params := map[string][]string{
		"commands": {"mcrcon -H localhost -p " + rcon + " \"say " + username + ": " + message + "\""},
	}
	input := &ssm.SendCommandInput{
		InstanceIds:  []string{instanceID},
		DocumentName: &documentName,
		Parameters:   params,
	}
	_, err = client.SendCommand(context.TODO(), input)
	if err != nil {
		log.Printf("Error! Couldn't send command: %v", err)
		return err
	}

	return nil
}

func AddToWhitelist(instanceID string, rcon string, username string) error {
	// Setup AWS session
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	// Send whitelist command to Pixelmon EC2 instance
	client := ssm.NewFromConfig(cfg)
	documentName := "AWS-RunShellScript"
	params := map[string][]string{
		"commands": {"mcrcon -H localhost -p " + rcon + " \"whitelist add " + username + "\""},
	}
	input := &ssm.SendCommandInput{
		InstanceIds:  []string{instanceID},
		DocumentName: &documentName,
		Parameters:   params,
	}
	_, err = client.SendCommand(context.TODO(), input)
	if err != nil {
		log.Printf("Error! Couldn't send command: %v", err)
		return err
	}

	return nil
}
