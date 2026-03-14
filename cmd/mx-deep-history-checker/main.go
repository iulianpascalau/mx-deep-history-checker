package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/iulianpascalau/mx-deep-history-checker/factory"
	"github.com/iulianpascalau/mx-deep-history-checker/internal/config"
	"github.com/iulianpascalau/mx-deep-history-checker/internal/reporter"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/urfave/cli"
)

// appVersion should be populated at build time using ldflags
// Usage examples:
// Linux/macOS:
//
//	go build -v -ldflags="-X main.appVersion=$(git describe --all | cut -c7-32)
var appVersion = "undefined"
var log = logger.GetOrCreate("checker")

var (
	proxyHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
VERSION:
   {{.Version}}
   {{end}}
`

	nodeDir = cli.StringFlag{
		Name:  "node-dir",
		Usage: "The root path to the node data directory",
		Value: "",
	}
	startEpoch = cli.Uint64Flag{
		Name:  "start-epoch",
		Usage: "The starting epoch number to check (inclusive)",
		Value: 0,
	}
	endEpoch = cli.Uint64Flag{
		Name:  "end-epoch",
		Usage: "The ending epoch number to check (inclusive). If omitted, goes to the highest one.",
		Value: math.MaxUint64,
	}
	checkStatic = cli.BoolFlag{
		Name:  "check-static",
		Usage: "Check the Static directory databases",
	}
	parallelEpochs = cli.UintFlag{
		Name:  "parallel-epochs",
		Usage: "The number of epochs to process in parallel",
		Value: 4,
	}
	shard = cli.StringFlag{
		Name:  "shard",
		Usage: "The shard to be checked. Example: Shard_0, Shard_1, Shard_2, Shard_metachain",
		Value: "Shard_0",
	}
)

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = proxyHelpTemplate
	app.Name = "Deep History Checker"
	app.Version = fmt.Sprintf("%s/%s/%s-%s", appVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	app.Usage = "This is the entry point for starting a checker tool that can quickly analyze the DB structure for Deep History"
	app.Flags = []cli.Flag{
		nodeDir,
		startEpoch,
		endEpoch,
		checkStatic,
		parallelEpochs,
		shard,
	}
	app.Authors = []cli.Author{
		{
			Name:  "Iulian Pascalau",
			Email: "iulian.pascalau@gmail.com",
		},
	}

	app.Action = run

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	log.Info("Starting deep history checker", "version", appVersion, "pid", os.Getpid())
}

func run(ctx *cli.Context) error {
	cfg := config.Config{
		NodeDir:                     ctx.GlobalString(nodeDir.Name),
		StartEpoch:                  ctx.GlobalUint64(startEpoch.Name),
		EndEpoch:                    ctx.GlobalUint64(endEpoch.Name),
		CheckStatic:                 ctx.GlobalBool(checkStatic.Name),
		ParallelEpochs:              ctx.GlobalUint(parallelEpochs.Name),
		Shard:                       ctx.GlobalString(shard.Name),
		MandatoryEpochDirs:          getMandatoryEpochDirs(),
		MandatoryStaticDirsForShard: getMandatoryStaticDirsForShard(),
		MandatoryStaticDirsForMeta:  getMandatoryStaticDirsForMeta(),
	}

	rep := reporter.NewReporter()

	ctxRun, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		cancel()
	}()

	err := factory.DeepHistoryCheck(ctxRun, rep, &cfg)
	if err != nil {
		return err
	}

	rep.PrintSummary()

	return nil
}

func getMandatoryEpochDirs() []string {
	return []string{
		"AccountsTrie",
		"BlockHeaders",
		"BootstrapData",
		"DbLookupExtensions/MiniblocksMetadata",
		"DbLookupExtensions_ResultsHashesByTx",
		"Logs",
		"MetaBlock",
		"MiniBlocks",
		"PeerAccountsTrie",
		"Receipts",
		"RewardTransactions",
		"ScheduledSCRs",
		"Transactions",
		"UnsignedTransactions",
	}
}

func getMandatoryStaticDirsForShard() []string {
	return []string{
		"DbLookupExtensions_EpochByHash",
		"DbLookupExtensions_ESDTSupplies",
		"DbLookupExtensions_MiniblockHashByTxHash",
		"DbLookupExtensions_RoundHash",
		"MetaHdrHashNonce",
		"ShardHdrHashNonce0",
		"StatusMetricsStorageDB",
	}
}

func getMandatoryStaticDirsForMeta() []string {
	return []string{
		"DbLookupExtensions_EpochByHash",
		"DbLookupExtensions_ESDTSupplies",
		"DbLookupExtensions_MiniblockHashByTxHash",
		"DbLookupExtensions_RoundHash",
		"MetaHdrHashNonce",
		"ShardHdrHashNonce0",
		"ShardHdrHashNonce1",
		"ShardHdrHashNonce2",
		"StatusMetricsStorageDB",
	}
}
