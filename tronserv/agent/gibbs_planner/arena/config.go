package main

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/skelterjohn/tronimoes/tronserv/agent/gibbs_planner"
)

// ArenaConfig is the on-disk YAML shape for arena.
type ArenaConfig struct {
	Players []PlayerConfig `yaml:"players"`
}

// PlayerSlot is one player in list order; name is the agent / player name.
type PlayerConfig struct {
	Name   string               `yaml:"name"`
	Config gibbs_planner.Config `yaml:",inline"`
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
