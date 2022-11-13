package http_test

import (
	"fmt"
	"net/http"
	"testing"
	appServer "zombie_locator/internal/http"
	"zombie_locator/internal/logger"
	"zombie_locator/internal/repository/zombie"
	"zombie_locator/internal/service/locator"

	"github.com/golang/mock/gomock"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

func TestServer_ZombieLocationsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	locatorService := locator.NewMockLocator(ctrl)
	appLog, err := logger.NewAppLogger()
	require.NoError(t, err)

	httpAddr := fmt.Sprintf(":%d", freeport.GetPort())
	appHTTPServer := appServer.NewServer(appLog, httpAddr, locatorService)
	go func() {
		require.NoError(t, appHTTPServer.Run())
	}()
	t.Cleanup(func() {
		require.NoError(t, appHTTPServer.Shutdown())
	})
	locatorService.EXPECT().Locate(gomock.Any(), float64(1), float64(2), float64(3)).Return([]zombie.Location{}, nil)
	requestEndpoint(t, httpAddr, 1, 2, 3)
}

func requestEndpoint(t *testing.T, host string, lat, lot, limit float64) {
	url := fmt.Sprintf("http://%s/zombies?lat=%f&lon=%f&limit=%f", host, lat, lot, limit)
	resp, err := http.Get(url)
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
