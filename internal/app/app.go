package app

import (
	"context"
	"log"
	"net/http"
	"tt/config"
	"tt/internal/api"
	"tt/internal/db"
	"tt/internal/handlers"
	"tt/internal/repository"
	"tt/internal/services"
)

func RunApp(ctx context.Context, cfg *config.Config) error {
	mongoClient, err := db.ConnectMongoDB(ctx, cfg)
	if err != nil {
		return err
	}
	defer func() {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	redisClient, err := db.ConnectRedisClient(ctx, cfg)
	if err != nil {
		return err
	}
	defer func() {
		err := redisClient.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	psqlPool, err := db.ConnectPsqlDB(ctx, cfg)
	if err != nil {
		return err
	}
	defer psqlPool.Close()

	_ = db.Migrate(cfg)

	authService := services.NewAuthService(repository.NewAuthRepository(psqlPool), redisClient)
	dynamicDataService := services.NewDynamicDataService(
		repository.NewDynamicDataRepository(
			mongoClient.Database(cfg.MongoDatabase),
		))

	authHandler := handlers.NewAuthHandler(authService)
	dynamicHandler := handlers.NewDynamicDataHandler(dynamicDataService)

	router := api.SetupRouter(authService, handlers.Handlers{
		Auth:        authHandler,
		DynamicData: dynamicHandler,
	})

	server := &http.Server{
		Addr:    ":" + cfg.ServerAddress,
		Handler: router,
	}

	log.Printf("Starting server on %s\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
