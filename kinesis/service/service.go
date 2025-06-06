package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	kinesisModel "vijju/kinesis/model"
	userLogsModel "vijju/user-logs/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"gorm.io/gorm"
)

type KinesisService struct {
	Db            *gorm.DB
	KinesisClient *kinesis.Client
}

func NewKinesisService(db *gorm.DB) *KinesisService {
	return &KinesisService{Db: db, KinesisClient: NewKinesisClient()}
}

func NewKinesisClient() *kinesis.Client {
	// Initialize Kinesis client (LocalStack)
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localstack:4566"}, nil
			})),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}))
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}
	kinesisClient := kinesis.NewFromConfig(cfg)
	return kinesisClient
}

func (kinesisService *KinesisService) PutData(event *kinesisModel.UserEvent, key int) {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
	} else {
		_, err = kinesisService.KinesisClient.PutRecord(context.Background(), &kinesis.PutRecordInput{
			StreamName:   aws.String("user-logs"),
			Data:         eventBytes,
			PartitionKey: aws.String(fmt.Sprintf("%d", key)),
		})
		if err != nil {
			log.Printf("Failed to publish to Kinesis: %v", err)
		}
	}
}

func (kinesisService *KinesisService) ReadDataFromKinesis() {
	// Start Kinesis consumer in a goroutine
	go func() {
		for {
			// Get shards
			describeOutput, err := kinesisService.KinesisClient.DescribeStream(context.Background(), &kinesis.DescribeStreamInput{
				StreamName: aws.String("user-logs"),
			})
			if err != nil {
				log.Printf("Failed to describe stream: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// Process each shard
			for _, shard := range describeOutput.StreamDescription.Shards {
				iteratorOutput, err := kinesisService.KinesisClient.GetShardIterator(context.Background(), &kinesis.GetShardIteratorInput{
					StreamName:        aws.String("user-logs"),
					ShardId:           shard.ShardId,
					ShardIteratorType: types.ShardIteratorTypeTrimHorizon,
				})
				if err != nil {
					log.Printf("Failed to get shard iterator for %s: %v", *shard.ShardId, err)
					continue
				}

				iterator := iteratorOutput.ShardIterator
				for {
					recordsOutput, err := kinesisService.KinesisClient.GetRecords(context.Background(), &kinesis.GetRecordsInput{
						ShardIterator: iterator,
					})
					if err != nil {
						log.Printf("Failed to get records: %v", err)
						break
					}

					for _, record := range recordsOutput.Records {
						var event kinesisModel.UserEvent
						if err := json.Unmarshal(record.Data, &event); err != nil {
							log.Printf("Failed to unmarshal event: %v", err)
							continue
						}
						// Insert into logs table
						logEntry := userLogsModel.Log{
							UserID:    event.UserID,
							Action:    event.Action,
							Timestamp: time.Now(),
						}
						if err := kinesisService.Db.Create(&logEntry).Error; err != nil {
							log.Printf("Failed to save log: %v", err)
						}
						log.Printf("Consumed: UserID=%d, Action=%s", event.UserID, event.Action)
					}

					iterator = recordsOutput.NextShardIterator
					if iterator == nil {
						break // Shard is closed
					}
					time.Sleep(1 * time.Second)
				}
			}
			time.Sleep(5 * time.Second) // Retry for new shards
		}
	}()
}
