package aws

import (
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/convox/logger"
	"github.com/convox/rack/api/structs"
)

var (
	cloudformationEventsTopic = os.Getenv("CLOUDFORMATION_EVENTS_TOPIC")
	customTopic               = os.Getenv("CUSTOM_TOPIC")
	notificationTopic         = os.Getenv("NOTIFICATION_TOPIC")
	sortableTime              = "20060102.150405.000000000"
)

// Logger is a package-wide logger
var Logger = logger.New("ns=provider.aws")

type AWSProvider struct {
	Region   string
	Endpoint string
	Access   string
	Secret   string
	Token    string

	BuildCluster      string
	Cluster           string
	Development       bool
	DockerImageAPI    string
	DynamoBuilds      string
	DynamoReleases    string
	EncryptionKey     string
	NotificationHost  string
	NotificationTopic string
	Password          string
	Rack              string
	RegistryHost      string
	SecurityGroup     string
	SettingsBucket    string
	Subnets           string
	SubnetsPrivate    string
	Vpc               string
	VpcCidr           string

	SkipCache bool
}

// NewProviderFromEnv returns a new AWS provider from env vars
func FromEnv() *AWSProvider {
	return &AWSProvider{
		Region:            os.Getenv("AWS_REGION"),
		Endpoint:          os.Getenv("AWS_ENDPOINT"),
		Access:            os.Getenv("AWS_ACCESS"),
		Secret:            os.Getenv("AWS_SECRET"),
		Token:             os.Getenv("AWS_TOKEN"),
		BuildCluster:      os.Getenv("BUILD_CLUSTER"),
		Cluster:           os.Getenv("CLUSTER"),
		Development:       os.Getenv("DEVELOPMENT") == "true",
		DockerImageAPI:    os.Getenv("DOCKER_IMAGE_API"),
		DynamoBuilds:      os.Getenv("DYNAMO_BUILDS"),
		DynamoReleases:    os.Getenv("DYNAMO_RELEASES"),
		EncryptionKey:     os.Getenv("ENCRYPTION_KEY"),
		NotificationHost:  os.Getenv("NOTIFICATION_HOST"),
		NotificationTopic: os.Getenv("NOTIFICATION_TOPIC"),
		Password:          os.Getenv("PASSWORD"),
		Rack:              os.Getenv("RACK"),
		RegistryHost:      os.Getenv("REGISTRY_HOST"),
		SecurityGroup:     os.Getenv("SECURITY_GROUP"),
		SettingsBucket:    os.Getenv("SETTINGS_BUCKET"),
		Subnets:           os.Getenv("SUBNETS"),
		SubnetsPrivate:    os.Getenv("SUBNETS_PRIVATE"),
		Vpc:               os.Getenv("VPC"),
		VpcCidr:           os.Getenv("VPCCIDR"),
	}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (p *AWSProvider) Initialize(opts structs.ProviderOptions) error {
	if opts.LogOutput != nil {
		Logger = logger.NewWriter("ns=provider.aws", opts.LogOutput)
	}

	return nil
}

/** services ****************************************************************************************/

func (p *AWSProvider) config() *aws.Config {
	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(p.Access, p.Secret, p.Token),
	}

	if p.Region != "" {
		config.Region = aws.String(p.Region)
	}

	if p.Endpoint != "" {
		config.Endpoint = aws.String(p.Endpoint)
	}

	if os.Getenv("DEBUG") == "true" || os.Getenv("CONVOX_DEBUG") == "true" {
		config.WithLogLevel(aws.LogDebugWithHTTPBody)
	}

	return config
}

func (p *AWSProvider) acm() *acm.ACM {
	return acm.New(session.New(), p.config())
}

func (p *AWSProvider) autoscaling() *autoscaling.AutoScaling {
	return autoscaling.New(session.New(), p.config())
}

func (p *AWSProvider) cloudformation() *cloudformation.CloudFormation {
	return cloudformation.New(session.New(), p.config())
}

func (p *AWSProvider) cloudwatch() *cloudwatch.CloudWatch {
	return cloudwatch.New(session.New(), p.config())
}

func (p *AWSProvider) cloudwatchlogs() *cloudwatchlogs.CloudWatchLogs {
	return cloudwatchlogs.New(session.New(), p.config())
}

func (p *AWSProvider) dynamodb() *dynamodb.DynamoDB {
	return dynamodb.New(session.New(), p.config())
}

func (p *AWSProvider) ec2() *ec2.EC2 {
	return ec2.New(session.New(), p.config())
}

func (p *AWSProvider) ecr() *ecr.ECR {
	return ecr.New(session.New(), p.config())
}

func (p *AWSProvider) ecs() *ecs.ECS {
	return ecs.New(session.New(), p.config())
}

func (p *AWSProvider) kms() *kms.KMS {
	return kms.New(session.New(), p.config())
}

func (p *AWSProvider) iam() *iam.IAM {
	return iam.New(session.New(), p.config())
}

func (p *AWSProvider) s3() *s3.S3 {
	return s3.New(session.New(), p.config().WithS3ForcePathStyle(true))
}

func (p *AWSProvider) sns() *sns.SNS {
	return sns.New(session.New(), p.config())
}

// IsTest returns true when we're in test mode
func (p *AWSProvider) IsTest() bool {
	return p.Region == "us-test-1"
}
