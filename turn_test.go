package edgeturn

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
)


func TestTurn(t *testing.T) {
	SetupTurn("54.179.156.111","1234567890","1234567890",65535,65535,65000)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}