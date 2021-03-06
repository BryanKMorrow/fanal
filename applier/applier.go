package applier

import (
	"github.com/BryanKMorrow/fanal/analyzer"
	"github.com/BryanKMorrow/fanal/cache"
	"github.com/BryanKMorrow/fanal/types"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"golang.org/x/xerrors"
	"strings"
)

type Applier struct {
	cache cache.LocalArtifactCache
}

func NewApplier(c cache.LocalArtifactCache) Applier {
	return Applier{cache: c}
}

func (a Applier) ApplyLayers(imageID string, diffIDs []string, manifest v1.Manifest) (types.ArtifactDetail, error) {
	var layers []types.BlobInfo
	for _, diffID := range diffIDs {
		dif := strings.Split(diffID, "/")
		layer, _ := a.cache.GetBlob(dif[0])
		if layer.SchemaVersion == 0 {
			return types.ArtifactDetail{}, xerrors.Errorf("layer cache missing: %s", diffID)
		}
		layers = append(layers, layer)
	}

	mergedLayer := ApplyLayers(layers)
	if mergedLayer.OS == nil {
		return mergedLayer, analyzer.ErrUnknownOS // send back package and apps info regardless
	} else if mergedLayer.Packages == nil {
		return mergedLayer, analyzer.ErrNoPkgsDetected // send back package and apps info regardless
	}

	imageInfo, _ := a.cache.GetArtifact(imageID)
	mergedLayer.HistoryPackages = imageInfo.HistoryPackages

	mergedLayer.Manifest = manifest

	return mergedLayer, nil
}
