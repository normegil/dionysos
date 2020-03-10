package listener_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/dionysos/internal/http/listener"
	"github.com/normegil/postgres"
	"net/http"
	"testing"
)

func TestListener(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test")
	}

	lst := initTest(t)
	t.Run("", func(t *testing.T) {

	})
}

func initTest(t testing.TB) http.Handler {
	cfg, containerCfg := postgres.Test_Deploy(t)
	defer postgres.Test_RemoveContainer(t, containerCfg.Identifier)

	lst := listener.NewListener(listener.Configuration{
		APILogErrors: false,
		Database: postgres.Configuration{
			Address:  cfg.Address,
			Port:     cfg.Port,
			User:     cfg.User,
			Password: cfg.Password,
			Database: "dionysos-listener-test-" + uuid.New().String(),
		},
	})
	defer func() {
		if err := lst.Close(); nil != err {
			t.Fatal(fmt.Errorf("closing listener: %w", err))
		}
	}()

	handler, err := lst.Load()
	if err != nil {
		t.Fatal(fmt.Errorf("load handler: %w", err))
	}
	return handler
}
