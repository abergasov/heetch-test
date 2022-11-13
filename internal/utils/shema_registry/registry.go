package shema_registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"zombie_locator/internal/entities"
)

type entityStreamEvent struct {
	Version int    `json:"v"`
	Data    string `json:"d"`
}

type Registry struct {
	supportedUserVersions map[int]struct{}
}

func NewRegistry(supportedUserVersions []int) *Registry {
	result := &Registry{
		supportedUserVersions: make(map[int]struct{}),
	}
	for _, version := range supportedUserVersions {
		result.supportedUserVersions[version] = struct{}{}
	}
	return result
}

func (r *Registry) encodeEvent(version int, payload interface{}, caster caster) ([]byte, error) {
	result := &entityStreamEvent{
		Version: version,
	}
	if _, ok := r.supportedUserVersions[version]; !ok {
		return nil, UnsupportedEventVersion
	}
	msg, ok := caster(version, payload)
	if !ok {
		return nil, UnsupportedEvent
	}
	data, err := r.encode(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode event: %w", err)
	}
	result.Data = data
	return json.Marshal(result)
}

func (r *Registry) decodeEvent(message []byte, deCaster deCaster) (map[int]interface{}, error) {
	var result entityStreamEvent
	if err := json.Unmarshal(message, &result); err != nil {
		return nil, err
	}
	if _, ok := r.supportedUserVersions[result.Version]; !ok {
		return nil, UnsupportedEventVersion
	}
	decoded, err := r.decode(result.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode event: %w", err)
	}
	payload, err := deCaster(result.Version, decoded)
	if err != nil {
		return nil, err
	}
	return map[int]interface{}{result.Version: payload}, nil
}

func (r *Registry) EncodeZombieLocationStreamEvent(version int, payload interface{}) ([]byte, error) {
	return r.encodeEvent(version, payload, func(v int, data any) (any, bool) {
		switch v {
		case 1:
			zLocation, ok := payload.(entities.ZombieLocationV1)
			return zLocation, ok
		case 2:
			zLocation, ok := payload.(entities.ZombieLocationV2)
			return zLocation, ok
			// for future enhancements...
		}
		return nil, false
	})
}

func (r *Registry) DecodeZombieLocationStreamEvent(message []byte) (map[int]interface{}, error) {
	return r.decodeEvent(message, func(v int, data []byte) (any, error) {
		switch v {
		case 1:
			var location entities.ZombieLocationV1
			err := json.Unmarshal(data, &location)
			return &location, err
		case 2:
			var location entities.ZombieLocationV2
			err := json.Unmarshal(data, &location)
			return &location, err
		}
		return nil, UnsupportedEvent
	})
}

func (r *Registry) EncodeZombieCapturedStreamEvent(version int, payload interface{}) ([]byte, error) {
	return r.encodeEvent(version, payload, func(v int, data any) (any, bool) {
		switch v {
		case 1:
			tsk, ok := payload.(entities.ZombieCapturedV1)
			return tsk, ok
		case 2:
			// for future upgrades...
			return nil, false
		}
		return nil, false
	})
}

func (r *Registry) DecodeZombieCapturedStreamEvent(message []byte) (map[int]interface{}, error) {
	return r.decodeEvent(message, func(v int, data []byte) (any, error) {
		switch v {
		case 1:
			var usr entities.ZombieCapturedV1
			err := json.Unmarshal(data, &usr)
			return &usr, err
		case 2:
			// for future upgrades...
			return nil, UnsupportedEvent
		}
		return nil, UnsupportedEvent
	})
}

func (r *Registry) encode(data any) (string, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func (r *Registry) decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
