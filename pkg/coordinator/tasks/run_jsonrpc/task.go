package generateeoatransactions

import (
	"context"
	"fmt"
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
		rpcLogger := t.logger.WithField("client", client.Config.Name).WithField("method", t.config.RPCMethod)
		rpcLogger.Info("sending rpc request")

		innerRPCClient := client.ExecutionClient.GetRPCClient().GetEthClient().Client()

		result := new(interface{})
		err := innerRPCClient.CallContext(ctx, &result, t.config.RPCMethod, t.config.Params...)
		if err != nil {
			rpcLogger.Errorf("rpc request failed %v", err)
			return err
		}
		rpcLogger.Info("successfully sent rpc request")

	}
	return nil
}
