package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() *pgxpool.Pool {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Println("[Info] DATABASE_URL not set; using local fallback connection string.")
	}

	ctx := context.Background()

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("[Error] Invalid DATABASE_URL: %v", err)
	}

	targetDB := cfg.ConnConfig.Database

	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err == nil {
		err = db.Ping(ctx)
	}

	if err == nil {
		if err := EnsureSchema(ctx, db); err != nil {
			db.Close()
			log.Fatalf("[Error] Unable to initialize database schema: %v", err)
		}
		log.Println("[Info] Connected to PostgreSQL successfully!")
		return db
	}

	db.Close()

	maintenanceCfg := *cfg
	maintenanceCfg.ConnConfig.Database = "postgres"

	maintenanceDB, err := pgxpool.NewWithConfig(ctx, &maintenanceCfg)
	if err != nil {
		log.Fatalf("[Error] Unable to connect to PostgreSQL server: %v", err)
	}

	if err := maintenanceDB.Ping(ctx); err != nil {
		maintenanceDB.Close()
		log.Fatalf("[Error] PostgreSQL server is not reachable: %v", err)
	}

	_, err = maintenanceDB.Exec(
		ctx,
		fmt.Sprintf(`CREATE DATABASE "%s"`, targetDB),
	)

	if err != nil {
		maintenanceDB.Close()
		log.Fatalf("[Error] Unable to create database %q: %v", targetDB, err)
	}

	log.Printf("[Info] Database %q created successfully!", targetDB)

	maintenanceDB.Close()

	db, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("[Error] Unable to connect to newly created database: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		log.Fatalf("[Error] Unable to ping newly created database: %v", err)
	}

	if err := EnsureSchema(ctx, db); err != nil {
		db.Close()
		log.Fatalf("[Error] Unable to initialize database schema after creating database: %v", err)
	}

	log.Println("[Info] Connected to PostgreSQL successfully!")

	return db
}
