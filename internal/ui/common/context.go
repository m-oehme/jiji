package common

import (
	"log/slog"

	"github.com/m-oehme/jiji/internal/config"
)

// Context holds immutable infrastructure shared by all components.
// It is created once in app.New and passed to every component constructor.
type Context struct {
	Config *config.Config
	Logger *slog.Logger
}
