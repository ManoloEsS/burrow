package app

import (
	"context"

	"github.com/ManoloEsS/burrow/internal/cli"
	"google.golang.org/grpc/balancer/grpclb/state"
)

type App struct {
	UI     cli.UI
	State  *state.State
	ctx    context.Context
	cancel context.CancelFunc
}
