package game

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"cloud.google.com/go/compute/metadata"
	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/skelterjohn/tronimoes/tronserv/clog"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type AgentSpawner interface {
	NewAgent(ctx context.Context, which string, code string, roundOut int) error
}

type LocalAgentSpawner struct {
}

func (s LocalAgentSpawner) NewAgent(ctx context.Context, which string, code string, roundOut int) error {
	ctx = clog.WithKeyword(ctx, "which", which)
	ctx = clog.WithKeyword(ctx, "code", code)
	clog.Info(ctx, "Spawning agent for game")
	exeDir := ""
	if exe, err := os.Executable(); err == nil {
		exeDir = filepath.Dir(exe)
	}
	// Spawn via PATH resolution (works when PATH includes the Go bin dir).
	agentExe := "agent"
	replicantExe := "replicant"
	// Use Background so the agent is not killed when the HTTP request context is cancelled.
	runCtx := context.WithoutCancel(ctx)
	cmd := exec.CommandContext(runCtx,
		replicantExe,
		fmt.Sprintf("%d", roundOut-1),
		agentExe,
		"--which", which,
		"--code", code,
		"--round-out", fmt.Sprintf("%d", roundOut),
		"--ready-with", fmt.Sprintf("%d", roundOut),
		"--dev",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = exeDir
	if exeDir != "" {
		clog.Info(ctx, "Running agent in", "exeDir", exeDir)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start agent: %w", err)
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			clog.Error(ctx, "Agent exited", err)
		}
	}()
	return nil
}

// GCRAgentSpawner runs agents as Google Cloud Run Jobs. The job must already
// exist in the project with the desired container image; this spawner only
// executes it with overridden args (e.g. --which, --code).
// When running inside Cloud Run, the default service account is used.
type GCRAgentSpawner struct {
	// Code is the code of the game to spawn agents for.
	Code string
	// JobsClient is the Cloud Run Jobs API client. If nil, NewAgent will create
	// one with default credentials (ADC).
	JobsClient *run.JobsClient
	// ProjectID is the GCP project (e.g. "my-project"). If empty and JobResourceName is empty, inferred from GCE metadata (same as this Cloud Run service).
	ProjectID string
	// Region is the Cloud Run region (e.g. "us-central1"). If empty and JobResourceName is empty, inferred from GCE metadata.
	Region string
	// JobName is the short job id (e.g. "tronagent"). Ignored if JobResourceName is set.
	JobName string
	// JobResourceName is the full resource name. If set, ProjectID/Region/JobName are ignored.
	// Format: projects/{project}/locations/{location}/jobs/{job}
	JobResourceName string
	// TronservAddr is the address of the tronimoes game server.
	TronservAddr string
}

// InferConfig fills in empty config from GCE metadata (project, region) and environment
// (GCR_AGENT_JOB, GCR_AGENT_CONTAINER). Call once at startup before using the spawner.
// Returns an error only when a required value is missing and metadata is unavailable.
func (s *GCRAgentSpawner) Initialize(ctx context.Context) error {
	if s.ProjectID == "" {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			return fmt.Errorf("project id not set and metadata unavailable: %w", err)
		}
		s.ProjectID = projectID
	}
	if s.Region == "" {
		regionPath, err := metadata.GetWithContext(ctx, "instance/region")
		if err != nil {
			return fmt.Errorf("region not set and metadata unavailable: %w", err)
		}
		s.Region = path.Base(regionPath)
	}
	if s.JobName == "" {
		s.JobName = "tronagent"
	}

	s.JobResourceName = fmt.Sprintf("projects/%s/locations/%s/jobs/%s", s.ProjectID, s.Region, s.JobName)

	if s.TronservAddr == "" {
		s.TronservAddr = "https://games.tronimoes.com"
	}

	if s.JobsClient == nil {
		creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
		if err != nil {
			return fmt.Errorf("application default credentials: %w", err)
		}
		s.JobsClient, err = run.NewJobsClient(ctx, option.WithTokenSource(creds.TokenSource))
		if err != nil {
			return fmt.Errorf("creating Cloud Run Jobs client: %w", err)
		}
	}
	return nil
}

func (s *GCRAgentSpawner) NewAgent(ctx context.Context, which string, code string, roundOut int) error {
	if s.Code != "" {
		code = s.Code
	}
	clog.Info(ctx, "Spawning agent via", "which", which, "code", code, "job", s.JobResourceName)
	req := &runpb.RunJobRequest{
		Name: s.JobResourceName,
		Overrides: &runpb.RunJobRequest_Overrides{
			ContainerOverrides: []*runpb.RunJobRequest_Overrides_ContainerOverride{
				{
					Args: []string{
						fmt.Sprintf("%d", roundOut-1),
						"/app/agent",
						"--addr", s.TronservAddr,
						"--which", which,
						"--code", code,
						"--round-out", fmt.Sprintf("%d", roundOut),
						"--ready-with", fmt.Sprintf("%d", roundOut),
						"--gce",
					},
				},
			},
		},
	}

	op, runErr := s.JobsClient.RunJob(ctx, req)
	if runErr != nil {
		return fmt.Errorf("run job %q: %w", s.JobResourceName, runErr)
	}
	clog.Info(ctx, "Job submitted", "job", op.Name(), "code", code)
	// Optionally wait for execution: op.Wait(ctx). For fire-and-forget we don't wait.
	return nil
}
