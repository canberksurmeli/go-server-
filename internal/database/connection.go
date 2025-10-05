package database

import (
	"context"
	"fmt"
	"message-provider-go/internal/config"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

var (
	instance *DB
	once     sync.Once
	ctx      = context.Background()
)

// Init initializes the database connection (call once at startup)
func Init() error {
	var err error
	once.Do(func() {
		cfg := config.Get()

		connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
		)

		fmt.Printf("Connecting to database: %s:%d/%s\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

		dbpool, poolErr := pgxpool.New(ctx, connString)
		if poolErr != nil {
			err = fmt.Errorf("failed to create connection pool: %w", poolErr)
			return
		}

		if pingErr := dbpool.Ping(ctx); pingErr != nil {
			err = fmt.Errorf("failed to ping database: %w", pingErr)
			return
		}

		instance = &DB{dbpool}
		fmt.Println("Connected to PostgreSQL database!")
	})
	return err
}

// Get returns the database instance
func Get() *DB {
	if instance == nil {
		fmt.Printf("Database not initialized. Call database.Init() first\n")
	}
	return instance
}

// Begin starts a new transaction
func (db *DB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.Pool.Begin(ctx)
}

// Close closes the database connection
func (db *DB) Close() {
	if instance != nil {
		instance.Pool.Close()
	}
}

func (db *DB) GetDB() *pgxpool.Pool {
	return db.Pool
}
