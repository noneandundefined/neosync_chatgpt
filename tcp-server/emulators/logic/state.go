package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"neomatica/neosync-tcp/pkg/protocol/common"
)

type DeviceState struct {
	IMEI            string `json:"imei"`
	Password        string `json:"password"`
	FirmwareVersion uint16 `json:"firmware_version"`
	CfgVersion      uint8  `json:"cfg_version"`
	LastModTime     uint32 `json:"last_mod_time"`
	CfgHash         uint32 `json:"cfg_hash"`
}

type StateStore struct {
	Dir string
}

func (s *StateStore) statePath() string {
	return filepath.Join(s.Dir, "state.json")
}

func (s *StateStore) configPath() string {
	return filepath.Join(s.Dir, "config.bin")
}

func (s *StateStore) configHexPath() string {
	return filepath.Join(s.Dir, "config.hex")
}

func (s *StateStore) Load(defaults DeviceState) (*DeviceState, []byte, error) {
	if err := os.MkdirAll(s.Dir, 0o755); err != nil {
		return nil, nil, err
	}

	cfg, err := s.loadConfigBytes()
	if err != nil {
		return nil, nil, err
	}

	state := defaults
	if raw, err := os.ReadFile(s.statePath()); err == nil {
		if err := json.Unmarshal(raw, &state); err != nil {
			return nil, nil, fmt.Errorf("parse state.json: %w", err)
		}
	}

	if len(cfg) > 0 {
		if hash := common.GetCfgHash(cfg); hash != 0 {
			state.CfgHash = hash
		}
		if ts := common.GetCfgLastModTime(cfg); ts != 0 {
			state.LastModTime = ts
		}
	}

	if err := s.Save(&state, cfg); err != nil {
		return nil, nil, err
	}

	return &state, cfg, nil
}

func (s *StateStore) loadConfigBytes() ([]byte, error) {
	if raw, err := os.ReadFile(s.configPath()); err == nil && len(raw) > 0 {
		return raw, nil
	}

	if raw, err := os.ReadFile(s.configHexPath()); err == nil && len(raw) > 0 {
		text := string(raw)
		cfg, err := hex.DecodeString(compactHex(text))
		if err != nil {
			return nil, fmt.Errorf("decode config.hex: %w", err)
		}
		packet, err := ensureCfgPacket(cfg)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(s.configPath(), packet, 0o644); err != nil {
			return nil, err
		}
		return packet, nil
	}

	return nil, nil
}

func (s *StateStore) Save(state *DeviceState, cfg []byte) error {
	if state == nil {
		return fmt.Errorf("state is nil")
	}

	if len(cfg) > 0 {
		if hash := common.GetCfgHash(cfg); hash != 0 {
			state.CfgHash = hash
		}
		if ts := common.GetCfgLastModTime(cfg); ts != 0 {
			state.LastModTime = ts
		}
	}

	raw, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.statePath(), raw, 0o644); err != nil {
		return err
	}

	if len(cfg) > 0 {
		if err := os.WriteFile(s.configPath(), cfg, 0o644); err != nil {
			return err
		}
		_ = os.WriteFile(s.configHexPath(), []byte(hex.EncodeToString(cfg)), 0o644)
	}

	return nil
}

func compactHex(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') {
			out = append(out, c)
		}
	}
	return string(out)
}
