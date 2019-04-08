package conf

import (
	"fmt"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/util"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"net"
	"strconv"
	"time"
)

const (
	ModeProd                 = "prod"
	ModeDev                  = "dev"
	DefaultServiceSphereName = "sphere-server"
)

var (
	SupportedNameSeparators  = []string{".", "#", "|"}
	SupportedConfigFileTypes = []string{".yaml"}
)

type Config struct {
	ServerConfig  *ServerConfig
	StratumConfig *StratumConfig
}

type ServerConfig struct {
	Name string
	Host string
	Port int
	Mode string
}

type GRPCConfig struct {
	Host string
	Port int
}

type StratumConfig struct {
	StratumSubscribeConfig
	StratumAuthorizeConfig
	StratumSetDifficultyConfig
	StratumNotifyConfig
	DefaultDifficultResults
}

type DefaultDifficultResults struct {
	results map[string]uint64
}

type StratumNotifyConfig struct {
	CoinType                    string
	PoolTag                     string
	PayoutAddress               string
	UsedTestNet                 bool
	NotifyLoopInterval          time.Duration
	DifficultyCheckLoopInterval time.Duration
	JobHashCheckLoopInterval    time.Duration
}

type StratumSetDifficultyConfig struct {
	DefaultDifficultyFilePath string
	MinDifficulty             uint64
	MaxDifficulty             uint64
	DefaultDifficulty         uint64
	DifficultyFactor          float64
}

type StratumSubscribeConfig struct {
	NotifyPlaceholder     string
	DifficultyPlaceholder string
	ExtraNonce1Length     int
	ExtraNonce2Length     int
}

type StratumAuthorizeConfig struct {
	NameSeparator string
	SkipVerify    bool
}

func NewConfig() *Config {
	return &Config{}
}

func (c *StratumConfig) IsValidAuthorizeNameSeparator() bool {
	for _, separator := range SupportedNameSeparators {
		if separator == c.NameSeparator {
			return true
		}
	}
	return false
}

func (c *ServerConfig) IsValidMode() bool {
	if c.Mode == ModeProd {
		return true
	}

	if c.Mode == ModeDev {
		return true
	}
	return false
}

func (c *ServerConfig) FormatHostPort() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func (c *StratumConfig) IsValidDefaultDifficultySupportedFileType() bool {
	if len(c.DefaultDifficultyFilePath) > 0 {
		fileType := path.Ext(c.DefaultDifficultyFilePath)
		for _, supportedType := range SupportedConfigFileTypes {
			if supportedType == fileType {
				return true
			}
		}
	}
	return false
}

func (c *StratumConfig) IsValidDefaultDifficultyFilePath() bool {
	return util.PathExist(c.DefaultDifficultyFilePath)
}

func (c *StratumConfig) LoadDefaultDifficultyResults() error {
	results := make(map[string]uint64)
	filePath := c.DefaultDifficultyFilePath

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileData, results)
	if err != nil {
		return err
	}

	c.DefaultDifficultResults.results = results
	return nil
}

func (r *DefaultDifficultResults) GetNotifyDefaultDifficulty(userAgent string) (uint64, error) {
	var ok bool
	var difficulty uint64
	if difficulty, ok = r.results[userAgent]; !ok {
		return 0, errors.New(fmt.Sprint("not found difficulty with user agent"))
	}
	return difficulty, nil
}

func (c *GRPCConfig) FormatHostPort() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func (c *Config) Init() error {
	if err := c.StratumConfig.LoadDefaultDifficultyResults(); err != nil {
		return err
	}

	return nil
}

func (c *Config) Validate() error {
	if !c.ServerConfig.IsValidMode() {
		return errors.New(fmt.Sprintf("invalid server mode :%s", c.ServerConfig.Mode))
	}

	if !c.StratumConfig.IsValidDefaultDifficultyFilePath() {
		return errors.New(fmt.Sprintf("invalid default difficulty file path : %s", c.StratumConfig.DefaultDifficultyFilePath))
	}

	if !c.StratumConfig.IsValidDefaultDifficultySupportedFileType() {
		return errors.New(fmt.Sprintf("invalid default difficulty file type : %s", c.StratumConfig.DefaultDifficultyFilePath))
	}

	if !c.StratumConfig.IsValidAuthorizeNameSeparator() {
		return errors.New(fmt.Sprintf("invalid name separator: %s", c.StratumConfig.NameSeparator))
	}

	if len(c.StratumConfig.PayoutAddress) == 0 {
		return errors.New("coin address is required.")
	}

	if len(c.StratumConfig.PoolTag) == 0 {
		return errors.New("coin tag is required.")
	}
	return nil
}
