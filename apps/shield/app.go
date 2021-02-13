package shield

import (
	"context"
	"log"
	"time"

	"github.com/phaesoo/shield/configs"
	"github.com/phaesoo/shield/internal/mq"
	"github.com/phaesoo/shield/internal/services/keyauth"
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

func (app *Shield) setupServices() {
	mc := app.config.Mysql
	db, err := db.NewDB("mysql", db.DSN(mc.User, mc.Password, mc.Database, mc.Host, mc.Port))
	if err != nil {
		panic(err)
	}

	rc := app.config.Redis
	pool := memdb.NewPool(rc.Address(), rc.Database)
	store := store.NewStore(pool)
	_ = mq.NewRedisMQ(pool)

	// Register routes
	keyauthService := keyauth.NewService(store, db)
	if err := keyauthService.Initialize(context.Background()); err != nil {
		panic(err)
	}
	keyauth := keyauth.NewServer(keyauthService)
	keyauth.RegisterRoutes(app.server.Router())
}

func (app *Shield) processEvents(handler func(ctx context.Context) error) {
	log.Print("Start to process events")
	for {
		handler(context.Background())

		time.Sleep(1 * time.Second)
	}
}

// Listen starts server on the address
func (app *Shield) Listen() error {
	return app.server.Listen()
}

// Shutdown server
func (app *Shield) Shutdown() error {
	if err := app.server.Shutdown(); err != nil {
		return err
	}
	return nil
}
