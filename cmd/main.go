// The zundoko-client command is a client of Zundoko Server.
package main

import (
	"github.com/kaitoy/zundoko-go-client/pkg/client"
	"github.com/kaitoy/zundoko-go-client/pkg/logging"
	"github.com/kaitoy/zundoko-go-client/pkg/runner"
	"go.uber.org/zap/zapcore"
)

func main() {
	logging.Init(zapcore.InfoLevel)
	defer logging.GetLogger().Sync()

	cl := client.NewClient("http://localhost:8080")
	if err := runner.NewRunner(cl).Run(1000); err != nil {
		logging.GetLogger().Errorw("An error occurred.", "err", err)
	}
}
