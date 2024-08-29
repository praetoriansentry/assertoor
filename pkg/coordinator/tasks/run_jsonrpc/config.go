package generateeoatransactions

import "fmt"

type Config struct {
	ClientPattern string `yaml:"clientPattern" json:"clientPattern"`
	RPCMethod     string `yaml:"method" json:"method"`
	Params        []any  `yaml:"params" json:"params"`
}

func DefaultConfig() Config {
	return Config{
		ClientPattern: "",
		RPCMethod:     "",
		Params:        []any{},
	}
}

func (c *Config) Validate() error {
	if c.RPCMethod == "" {
		return fmt.Errorf("no RPC Method specified")
	}
	return nil
}
