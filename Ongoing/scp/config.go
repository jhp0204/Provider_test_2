package scp

// 각 상품별 sdk 정보를 가져와서 진행, 이후 scp용 sdk 작성 후, import package 목록으로 수정, 반영 가능 
import (
	"time"
	"github.com/jhp0204/Provider_test1/scp-clilib"
)

// DefaultWaitForInterval is Interval for checking status in WaitForXXX method < 큰 수정 불필요 판단 
const DefaultWaitForInterval = 10

// Default timeout
const DefaultTimeout = 5 * time.Minute
const DefaultCreateTimeout = 1 * time.Hour
const DefaultUpdateTimeout = 10 * time.Minute
const DefaultStopTimeout = 5 * time.Minute

// SCP도 AccessKey, SecretKey는 동일 
type Config struct {
	AccessKey string
	SecretKey string
}

// 앞선 상품별 sdk를 연결, 역시 변경해야하는 부분, vpc 제외하고 전체 삭제 
// ncloud vpc sdk 부분 변경 필요
type ScpAPIClient struct {
	vpc           *vpc.APIClient
}

// ncloud sdk 부분 변경 필요
func (c *Config) Client() (*ScpAPIClient, error) {
	apiKey := &scp-clilib.APIKey{
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	}
// 앞선 상품별 sdk를 연결, 역시 변경해야하는 부분, vpc 제외하고 전체 삭제 
// ncloud vpc sdk 부분 변경 필요
	return &ScpAPIClient{
		vpc:           vpc.NewAPIClient(vpc.NewConfiguration(apiKey)),
	}, nil
}

// ncloud vpc sdk 부분 변경 필요
type ProviderConfig struct {
	Site       string
	SupportVPC bool
	RegionCode string
	RegionNo   string
	Client     *ScpAPIClient
}
