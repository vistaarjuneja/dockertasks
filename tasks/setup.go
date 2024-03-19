package tasks

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"unicode"

	"github.com/harness/lite-engine/api"
	"github.com/harness/lite-engine/engine"
	"github.com/harness/lite-engine/engine/spec"
)

var (
	ctxBg = context.Background()
)

type setupRequest struct {
	ID               string            `json:"id"` // stage runtime ID
	PoolID           string            `json:"pool_id"`
	Tags             map[string]string `json:"tags"`
	CorrelationID    string            `json:"correlation_id"`
	LogKey           string            `json:"log_key"`
	InfraType        string            `json:"infra_type"`
	api.SetupRequest `json:"setup_request"`
}

// setupRequest(id) creates a Request object with the given id.
// It sets the network as the same ID (stage runtime ID which is unique)
func SetupRequest(id string) setupRequest {
	fmt.Printf("in setup request, id is: %s", id)
	return setupRequest{
		ID: id,
		SetupRequest: api.SetupRequest{
			Network: spec.Network{
				ID: sanitize(id),
			},
			Volumes: []*spec.Volume{
				{
					HostPath: &spec.VolumeHostPath{
						ID:     "harness",
						Path:   generatePath(id),
						Create: true,
						Remove: true,
					},
				},
			},
		},
	}
}

func generatePath(id string) string {
	return fmt.Sprintf("/tmp/harness/%s", sanitize(id))
}

func sanitize(id string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return '_'
	}, id)
}

func HandleSetup(s setupRequest) error {
	if s.MountDockerSocket == nil || *s.MountDockerSocket { // required to support m1 where docker isn't installed.
		s.Volumes = append(s.Volumes, getDockerSockVolume())
	}
	cfg := &spec.PipelineConfig{
		Envs:    s.Envs,
		Network: s.Network,
		Platform: spec.Platform{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		Volumes:           s.Volumes,
		Files:             s.Files,
		EnableDockerSetup: s.MountDockerSocket,
		TTY:               s.TTY,
	}
	if err := engine.SetupPipeline(ctxBg, engine.Opts{}, cfg); err != nil {
		return err
	}
	return nil
}

func getDockerSockVolume() *spec.Volume {
	path := engine.DockerSockUnixPath
	if runtime.GOOS == "windows" {
		path = engine.DockerSockWinPath
	}
	return &spec.Volume{
		HostPath: &spec.VolumeHostPath{
			Name: engine.DockerSockVolName,
			Path: path,
			ID:   "docker",
		},
	}
}

func getDockerSockVolumeMount() *spec.VolumeMount {
	path := engine.DockerSockUnixPath
	if runtime.GOOS == "windows" {
		path = engine.DockerSockWinPath
	}
	return &spec.VolumeMount{
		Name: engine.DockerSockVolName,
		Path: path,
	}
}
