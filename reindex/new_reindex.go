package reindex

import (
	"github.com/nnqq/scr-proto/codegen/go/org"
	"github.com/rs/zerolog"
)

func NewReindex(logger zerolog.Logger, orgClient org.OrgClient) Reindex {
	return Reindex{
		logger:    logger,
		orgClient: orgClient,
	}
}
