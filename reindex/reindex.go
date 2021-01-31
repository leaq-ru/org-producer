package reindex

import (
	"github.com/nnqq/scr-proto/codegen/go/org"
	"github.com/rs/zerolog"
)

type Reindex struct {
	logger    zerolog.Logger
	orgClient org.OrgClient
}
