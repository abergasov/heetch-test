package zombie

import (
	"context"
	"fmt"
	"time"
	"zombie_locator/internal/storage/db"

	"github.com/xjem/t38c"

	"github.com/google/uuid"
)

const (
	t38Key         = "zombies"
	capturedStatus = "status"
	locatedStatus  = "located"
)

type Zombie struct {
	dbConnect  db.Connector
	t38Connect *t38c.Client
}

func NewZombieRepository(dbConnect db.Connector, connect *t38c.Client) *Zombie {
	// TODO load from postgres to tile38 actual zombie list on server start
	return &Zombie{
		t38Connect: connect,
		dbConnect:  dbConnect,
	}
}

func (z *Zombie) CapturedZombie(ctx context.Context, zombieID uuid.UUID, updatedAt string) error {
	data, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return fmt.Errorf("unable to parse time: %w", err)
	}
	if _, err = z.dbConnect.Client().NamedExec(`
		INSERT INTO zombies(id, updated_at, status)
		VALUES(:id, :date, :status)
		ON CONFLICT (id) DO UPDATE SET updated_at = :date, status = :status;
	`, map[string]interface{}{
		"id":     zombieID,
		"date":   data,
		"status": capturedStatus,
	}); err != nil {
		return fmt.Errorf("unable to capture zombie: %w", err)
	}
	if err = z.t38Connect.Keys.Del(t38Key, zombieID.String()); err != nil {
		return fmt.Errorf("unable to delete zombie from tile38: %w", err)
	}
	return nil
}

func (z *Zombie) LocatedZombie(ctx context.Context, zombieID uuid.UUID, lat, lon float64, updatedAt string) error {
	data, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return fmt.Errorf("unable to parse time: %w", err)
	}
	if _, err = z.dbConnect.Client().NamedExec(`
		INSERT INTO zombies(id, updated_at, point, status)
		VALUES(:id, :date, point(:lat, :lon), :status)
		ON CONFLICT (id) DO UPDATE SET updated_at = :date, point = point(:lat, :lon), status = :status;
	`, map[string]interface{}{
		"id":     zombieID,
		"date":   data,
		"lat":    lat,
		"lon":    lon,
		"status": locatedStatus,
	}); err != nil {
		return fmt.Errorf("unable to locate zombie: %w", err)
	}
	if err = z.t38Connect.Keys.Set(t38Key, zombieID.String()).Point(lat, lon).Do(); err != nil {
		return fmt.Errorf("unable to save zombie to tile38: %w", err)
	}
	return nil
}

func (z *Zombie) LocateZombieList(ctx context.Context, lat, lon, limitKm float64) ([]Location, error) {
	nearbyRes, err := z.t38Connect.Search.Nearby(t38Key, lat, lon, limitKm*1000).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get nearby zombies: %w", err)
	}
	if nearbyRes.Count == 0 {
		return []Location{}, nil
	}
	result := make([]Location, 0, nearbyRes.Count)
	for i := range nearbyRes.Objects {
		zID, err := uuid.Parse(nearbyRes.Objects[i].ID)
		if err != nil {
			return nil, fmt.Errorf("unable to parse uuid: %w", err)
		}
		result = append(result, Location{
			ZombieId:  zID,
			Latitude:  nearbyRes.Objects[i].Object.Geometry.Point[0],
			Longitude: nearbyRes.Objects[i].Object.Geometry.Point[0],
		})
	}
	return result, nil
}
