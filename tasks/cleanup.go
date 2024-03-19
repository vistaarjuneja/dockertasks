package tasks

import (
	"fmt"

	"github.com/harness/lite-engine/api"
	"github.com/harness/lite-engine/engine"
	"github.com/harness/lite-engine/engine/spec"
)

var (
	// this label is used to identify steps associated with a pipeline
	// It's used only internally to successfully destroy containers.
	internalStageLabel = "internal_stage_label"
)

type destroyRequest struct {
	PipelineConfig spec.PipelineConfig `json:"pipeline_config"`
	api.DestroyRequest
}

// destroyRequest(id) creates a DestroyRequest object with the given id.
func DestroyRequest(stepIDs []string, stageID string) destroyRequest {
	fmt.Printf("in destroy request, id is: %s", stageID)
	return destroyRequest{
		DestroyRequest: api.DestroyRequest{
			StageRuntimeID: stageID,
		},
		PipelineConfig: spec.PipelineConfig{
			Network: spec.Network{
				ID: sanitize(stageID),
			},
		},
	}
}

func HandleDestroy(s destroyRequest) error {
	return engine.DestroyPipeline(ctxBg, engine.Opts{}, &s.PipelineConfig, internalStageLabel, s.StageRuntimeID)
}
