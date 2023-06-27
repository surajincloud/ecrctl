package awsutil

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

type RepositoriesTags struct {
	Key, Value string
}

type Repository struct {
	Name             string
	Arn              string
	Uri              string
	Age              string
	CreatedAt        *time.Time
	Tags             []RepositoriesTags
	TagMutability    types.ImageTagMutability
	EncryptionType   types.EncryptionType
	EncryptionKMSKey string
	ScanOnPush       bool
}

type Image struct {
	Digest                string
	Tag                   string
	ArtifactMediaType     string
	TagORDigest           string
	CriticalVulnerability int32
	HighVulnerability     int32
	MediumVulnerability   int32
	Age                   string
	Size                  string
	ScanStatus            types.ScanStatus
	ScanStatusDesc        string
	BasicScanFindings     []types.ImageScanFinding
	EnhancedScanFindings  []types.EnhancedImageScanFinding
}
