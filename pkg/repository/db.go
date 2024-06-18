package repository

import (
    "context"
    "fmt"
    "time"

	"github.com/rs/zerolog"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBConfig holds the MongoDB connection configuration.
type MongoDBConfig struct {
    URI      string        `envconfig:"DATABASE_URI" default:"mongodb://localhost:27017"`
    Database string        `envconfig:"DATABASE_NAME" default:"hodlhero"`
    Timeout  time.Duration `envconfig:"DATABASE_TIMEOUT" default:"10s"`
}

type MongoDBRepository struct {
    client *mongo.Client
    log    *zerolog.Logger
    cfg    *MongoDBConfig
}

func NewMongoDBRepository(log *zerolog.Logger, cfg *MongoDBConfig) (*MongoDBRepository, error) {
    ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
    defer cancel()

    clientOptions := options.Client().ApplyURI(cfg.URI)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    if err = client.Ping(ctx, nil); err != nil {
        return nil, fmt.Errorf("MongoDB ping failed: %w", err)
    }

    log.Info().Msg("Connected to MongoDB")
    return &MongoDBRepository{
        client: client,
        log:    log,
        cfg:    cfg,
    }, nil
}

func (r *MongoDBRepository) Disconnect(ctx context.Context) error {
    if err := r.client.Disconnect(ctx); err != nil {
        return fmt.Errorf("MongoDB disconnect failed: %w", err)
    }
    r.log.Info().Msg("Disconnected from MongoDB")
    return nil
}
