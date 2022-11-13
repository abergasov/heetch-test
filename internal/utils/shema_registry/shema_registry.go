package shema_registry

import "errors"

var (
	UnsupportedEvent        = errors.New("unsupported event")
	UnsupportedEventVersion = errors.New("unsupported version event")
)

// simply encode|decode to base 64. update to usage protobuf.
type caster func(v int, data any) (any, bool)
type deCaster func(v int, data []byte) (any, error)

type SchemaRegistry interface {
	EncodeZombieLocationStreamEvent(version int, payload interface{}) ([]byte, error)
	DecodeZombieLocationStreamEvent(message []byte) (map[int]interface{}, error)

	EncodeZombieCapturedStreamEvent(version int, payload interface{}) ([]byte, error)
	DecodeZombieCapturedStreamEvent(message []byte) (map[int]interface{}, error)
}
