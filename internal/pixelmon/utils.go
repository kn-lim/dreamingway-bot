package pixelmon

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

const (
	delay = 10 // Seconds
)

func getConfig() (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("PIXELMON_REGION")))
	if err != nil {
		log.Printf("Error! Couldn't create AWS config: %v", err)
		return aws.Config{}, err
	}

	return cfg, nil
}

func getInstance(cfg aws.Config, instanceID string) (*ec2.Client, ec2Types.Instance, error) {
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{
			instanceID,
		},
	}

	result, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		log.Printf("Error! Couldn't get the description of the instance: %v", err)
		return nil, ec2Types.Instance{}, err
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		log.Printf("Error! Couldn't find the instance [%s]: %v", instanceID, err)
		return nil, ec2Types.Instance{}, err
	}

	return client, result.Reservations[0].Instances[0], nil
}

func createInstanceDNSEntry(cfg aws.Config, instance ec2Types.Instance, zoneID string, url string) error {
	publicIP := instance.PublicIpAddress

	client := route53.NewFromConfig(cfg)

	log.Printf("Creating A record of %v to %v", *publicIP, url)

	_, err := client.ChangeResourceRecordSets(context.TODO(), &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &zoneID,
		ChangeBatch: &route53Types.ChangeBatch{
			Changes: []route53Types.Change{
				{
					Action: route53Types.ChangeActionUpsert,
					ResourceRecordSet: &route53Types.ResourceRecordSet{
						Name: &url,
						Type: route53Types.RRTypeA,
						TTL:  aws.Int64(300),
						ResourceRecords: []route53Types.ResourceRecord{
							{
								Value: publicIP,
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error! Couldn't create A record: %v", err)
		return err
	}

	log.Printf("Created A record of %v to %v", *publicIP, url)

	return nil
}

func removeInstanceDNSEntry(cfg aws.Config, instance ec2Types.Instance, zoneID string, url string) error {
	publicIP := instance.PublicIpAddress

	client := route53.NewFromConfig(cfg)

	log.Printf("Removing A record of %v to %v", *publicIP, url)

	_, err := client.ChangeResourceRecordSets(context.TODO(), &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &zoneID,
		ChangeBatch: &route53Types.ChangeBatch{
			Changes: []route53Types.Change{
				{
					Action: route53Types.ChangeActionDelete,
					ResourceRecordSet: &route53Types.ResourceRecordSet{
						Name: &url,
						Type: route53Types.RRTypeA,
						TTL:  aws.Int64(300),
						ResourceRecords: []route53Types.ResourceRecord{
							{
								Value: publicIP,
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error! Couldn't remove A record: %v", err)
		return err
	}

	log.Printf("Removed A record of %v to %v", *publicIP, url)

	return nil
}
