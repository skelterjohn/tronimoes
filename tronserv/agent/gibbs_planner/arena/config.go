package main

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// ArenaConfig is the on-disk YAML shape for arena.
type ArenaConfig struct {
	Players []PlayerSlot `yaml:"players"`
}

// PlayerSlot is one player in list order; name is the agent / player name.
type PlayerSlot struct {
	Name          string `yaml:"name"`
	PlannerConfig `yaml:",inline"`
}

// PlannerConfig is tuning for an agent; pass-through to the agent binary over time.
type PlannerConfig struct {
	Which                 string  `yaml:"which"` // agent -which, e.g. gibbs, random
	MaxInferenceTimeMs    int     `yaml:"max_inference_time_ms"`
	MaxSimulationTimeMs   int     `yaml:"max_simulation_time_ms"`
	MaxSimulationDepth    int     `yaml:"max_simulation_depth"`
	MaxSimulationsPerMove int     `yaml:"max_simulations_per_move"`
	ValueDecay            float64 `yaml:"value_decay"`
}

// LoadArenaConfig reads and validates arena YAML from path.
func LoadArenaConfig(path string) (ArenaConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return ArenaConfig{}, err
	}
	var c ArenaConfig
	if err := yaml.Unmarshal(b, &c); err != nil {
		return ArenaConfig{}, err
	}
	if len(c.Players) < 1 {
		return ArenaConfig{}, errors.New("players: need at least one entry")
	}
	seen := make(map[string]struct{}, len(c.Players))
	for _, p := range c.Players {
		if p.Name == "" {
			return ArenaConfig{}, errors.New("players: empty name")
		}
		if _, dup := seen[p.Name]; dup {
			return ArenaConfig{}, errors.New("players: duplicate name " + p.Name)
		}
		seen[p.Name] = struct{}{}
	}
	return c, nil
}
