package shield

import (
	"context"
	"log"
	"time"

	"github.com/phaesoo/shield/configs"
	"github.com/phaesoo/shield/internal/mq"
	"github.com/phaesoo/shield/internal/store"
	"github.com/phaesoo/shield/pkg/db"
	"github.com/phaesoo/shield/pkg/memdb"
	"github.com/phaesoo/shield/pkg/server"
)

type Shield struct {
	server *server.Server
	config configs.Config
}

func NewApp(config configs.Config) *Shield {
	ac := config.App

	server := &Shield{
		server: server.NewServer(ac.Address(), server.ServerConfig{
			Profile: ac.Profile,
			Metrics: ac.Metrics,
		}),
		config: config,
	}

	server.setupServices()

	return server
}

func (m *Shield) setupServices() {
	mc := m.config.Mysql
	_, err := db.NewDB("mysql", db.DSN(mc.User, mc.Password, mc.Database, mc.Host, mc.Port))
	if err != nil {
		panic(err)
	}

	rc := m.config.Redis
	pool := memdb.NewPool(rc.Address(), rc.Database)
	_ = store.NewStore(pool)
	_ = mq.NewRedisMQ(pool)

	// Register routes
}

func (m *Shield) processEvents(handler func(ctx context.Context) error) {
	log.Print("Start to process events")
	for {
		handler(context.Background())

		time.Sleep(1 * time.Second)
	}
}

// Listen starts server on the address
func (m *Shield) Listen() error {
	return m.server.Listen()
}

// Shutdown server
func (m *Shield) Shutdown() error {
	if err := m.server.Shutdown(); err != nil {
		return err
	}
	return nil
}
