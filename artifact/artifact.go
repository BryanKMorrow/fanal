package artifact

import (
	"context"

	"github.com/BryanKMorrow/fanal/types"
)

type Artifact interface {
	Inspect(ctx context.Context) (reference types.ArtifactReference, err error)
}
