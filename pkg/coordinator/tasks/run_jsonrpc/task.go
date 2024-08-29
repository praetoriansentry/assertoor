package generateeoatransactions

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethpandaops/assertoor/pkg/coordinator/types"
	"github.com/ethpandaops/assertoor/pkg/coordinator/wallet"
	"github.com/sirupsen/logrus"
)

var (
	TaskName       = "run_jsonrpc"
	TaskDescriptor = &types.TaskDescriptor{
		Name:        TaskName,
		Description: "Make a JSON RPC request to an endpoint",
		Config:      DefaultConfig(),
		NewTask:     NewTask,
	}
)

type Task struct {
	ctx        *types.TaskContext
	options    *types.TaskOptions
	config     Config
	logger     logrus.FieldLogger
	txIndex    uint64
	wallet     *wallet.Wallet
	walletPool *wallet.WalletPool

	targetAddr      common.Address
	transactionData []byte
}

func NewTask(ctx *types.TaskContext, options *types.TaskOptions) (types.Task, error) {
	return &Task{
		ctx:     ctx,
		options: options,
		logger:  ctx.Logger.GetLogger(),
	}, nil
}

func (t *Task) Config() interface{} {
	return t.config
}

func (t *Task) Timeout() time.Duration {
	return t.options.Timeout.Duration
}

func (t *Task) LoadConfig() error {
	config := DefaultConfig()

	// parse static config
	if t.options.Config != nil {
		if err := t.options.Config.Unmarshal(&config); err != nil {
			return fmt.Errorf("error parsing task config for %v: %w", TaskName, err)
		}
	}

	// load dynamic vars
	err := t.ctx.Vars.ConsumeVars(&config, t.options.ConfigVars)
	if err != nil {
		return err
	}

	// validate config
	if valerr := config.Validate(); valerr != nil {
		return valerr
	}

	t.config = config

	return nil
}

func (t *Task) Execute(ctx context.Context) error {
	for _, client := range t.ctx.Scheduler.GetServices().ClientPool().GetClientsByNamePatterns(t.config.ClientPattern, "") {
		rpcLogger := t.logger.WithField("client", client.Config.Name).WithField("method", t.config.RPCMethod).WithField("title", t.options.Title)
		rpcLogger.Info("sending rpc request")

		innerRPCClient := client.ExecutionClient.GetRPCClient().GetEthClient().Client()

		var result any

		err := innerRPCClient.CallContext(ctx, &result, t.config.RPCMethod, t.config.Params...)
		if t.config.ExpectError && err == nil {
			return fmt.Errorf("an error was expected, but we received a valid response")
		}
		if !t.config.ExpectError && err != nil {
			return fmt.Errorf("an error was received, but we expected a valid response %w", err)
		}
		if t.config.ExpectError {
			result = err
		}
		resultBytes, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("unable to json encode the json rpc response %w", err)
		}

		if t.config.ExpectError && t.config.ExpectResponseCode != 0 {
			err = checkErrorResponse(resultBytes, t.config.ExpectResponseCode)
			if err != nil {
				return err
			}
		}

		if t.config.ResponsePattern != "" {
			err = checkResponsePattern(resultBytes, t.config.ResponsePattern)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkErrorResponse(response []byte, expectedCode int) error {
	type jsonError struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
	var je jsonError
	err := json.Unmarshal(response, &je)
	if err != nil {
		return fmt.Errorf("error body didn't contain a normal json error %w", err)
	}
	if je.Code != expectedCode {
		return fmt.Errorf("expected response code %d, got %d", expectedCode, je.Code)
	}
	return nil

}

func checkResponsePattern(response []byte, responsePattern string) error {
	re, err := regexp.Compile(responsePattern)
	if err != nil {
		return fmt.Errorf("failed to compile response pattern %s: %w", responsePattern, err)
	}
	if !re.Match(response) {
		return fmt.Errorf("the response pattern %s did not match the response: %s", responsePattern, string(response))
	}
	return nil

}
