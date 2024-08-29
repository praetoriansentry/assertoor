package generateeoatransactions

import (
	"fmt"
)

type Config struct {
	ClientPattern      string `yaml:"clientPattern" json:"clientPattern"`
	RPCMethod          string `yaml:"method" json:"method"`
	Params             []any  `yaml:"params" json:"params"`
	ExpectError        bool   `yaml:"expectError" json:"expectError"`
	ExpectResponseCode int    `yaml:"expectResponseCode" json:"expectResponseCode"`
	ResponsePattern    string `yaml:"responsePattern" json:"responsePattern"`
}

func DefaultConfig() Config {
	return Config{
		ClientPattern:      "",
		RPCMethod:          "",
		Params:             []any{},
		ExpectError:        false,
		ExpectResponseCode: 0,
		ResponsePattern:    "",
	}
}

func (c *Config) Validate() error {
	if c.RPCMethod == "" {
		return fmt.Errorf("no RPC Method specified")
	}
	return nil
}
