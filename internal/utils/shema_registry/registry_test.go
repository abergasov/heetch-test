package shema_registry_test

import (
	"testing"
	"time"
	"zombie_locator/internal/entities"
	"zombie_locator/internal/utils/shema_registry"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRegistry_Encode_DecodeZombieCapturedStreamEvent(t *testing.T) {
	registry := shema_registry.NewRegistry([]int{1})
	payload := entities.ZombieCapturedV1{
		ZombieID:  uuid.New(),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	data, err := registry.EncodeZombieCapturedStreamEvent(1, payload)
	require.NoError(t, err)
	res, err := registry.DecodeZombieCapturedStreamEvent(data)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
	decodedPayload, ok := res[1].(*entities.ZombieCapturedV1)
	require.True(t, ok)
	require.Equal(t, payload, *decodedPayload)
}

func TestRegistry_Encode_DecodeZombieLocationStreamEvent(t *testing.T) {
	registry := shema_registry.NewRegistry([]int{1, 2})
	t.Run("v1", func(t *testing.T) {
		payload := entities.ZombieLocationV1{
			ZombieID:  uuid.New(),
			Latitude:  123.123,
			Longitude: 456.456,
			UpdatedAt: time.Now().Format(time.RFC3339),
		}
		data, err := registry.EncodeZombieLocationStreamEvent(1, payload)
		require.NoError(t, err)
		res, err := registry.DecodeZombieLocationStreamEvent(data)
		require.NoError(t, err)
		require.Equal(t, 1, len(res))
		decodedPayload, ok := res[1].(*entities.ZombieLocationV1)
		require.True(t, ok)
		require.Equal(t, payload, *decodedPayload)
	})
	t.Run("v2", func(t *testing.T) {
		payload := entities.ZombieLocationV2{
			ZombieID:  uuid.New(),
			Type:      "witch",
			Latitude:  789.789,
			Longitude: 456.456,
			UpdatedAt: time.Now().Format(time.RFC3339),
		}
		data, err := registry.EncodeZombieLocationStreamEvent(2, payload)
		require.NoError(t, err)
		res, err := registry.DecodeZombieLocationStreamEvent(data)
		require.NoError(t, err)
		require.Equal(t, 1, len(res))
		decodedPayload, ok := res[2].(*entities.ZombieLocationV2)
		require.True(t, ok)
		require.Equal(t, payload, *decodedPayload)
	})
}
