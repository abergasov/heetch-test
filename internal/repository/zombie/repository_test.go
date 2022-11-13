package zombie_test

import (
	"context"
	"testing"
	"time"
	"zombie_locator/internal/repository/zombie"
	"zombie_locator/internal/storage/db"

	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
)

const (
	postgresqlURL = "postgresql://test:secret@localhost:5432/zombies?sslmode=disable"
	tile38URL     = "localhost:9851"
)

func TestZombie_LocatedAndCapturedZombie(t *testing.T) {
	connect, err := db.NewTile38Connection(tile38URL)
	require.NoError(t, err)
	pgConnect, err := db.NewPostgresConnection(postgresqlURL)
	require.NoError(t, err)

	repo := zombie.NewZombieRepository(pgConnect, connect)

	zombieID := uuid.New()

	// locate zombie
	err = repo.LocatedZombie(context.Background(), zombieID, 48.85905, 2.294533, time.Now().Format(time.RFC3339))
	require.NoError(t, err)

	// check zombie is in hunting list
	// get zombie list
	checkZombie(t, repo, zombieID, 48.872544, 2.332298, 5, true)

	// capture zombie
	err = repo.CapturedZombie(context.Background(), zombieID, time.Now().Format(time.RFC3339))
	require.NoError(t, err)

	// check zombie is not in hunting list
	checkZombie(t, repo, zombieID, 48.872544, 2.332298, 5, false)
}

func checkZombie(t *testing.T, repo zombie.Zombier, zombieID uuid.UUID, lat, lon, limitKm float64, shouldExist bool) {
	list, err := repo.LocateZombieList(context.Background(), lat, lon, limitKm)
	require.NoError(t, err)
	require.True(t, len(list) > 0)
	found := false
	for _, l := range list {
		if l.ZombieId == zombieID {
			found = true
		}
	}
	if shouldExist {
		require.True(t, found)
	} else {
		require.False(t, found)
	}
}
