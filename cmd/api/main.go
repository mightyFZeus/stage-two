package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mightyzeus/stage-two/internal/db"
	"github.com/mightyzeus/stage-two/internal/env"
	"github.com/mightyzeus/stage-two/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Could not load .env file, falling back to defaults")
	}

	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiUrl: env.GetString("API_URL", "localhost:8080"),
		db: dbConfig{
			dbAddr: env.GetString("DB_ADDR", "admin:adminpassword@tcp(127.0.0.1:3307)/stage_two?charset=utf8mb4&parseTime=True&loc=Local"),

			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	gormDB, err := db.New(cfg.db.dbAddr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("failed to get underlying sql.DB:", err)
	}
	if err := store.AutoMigrate(gormDB); err != nil {
		log.Fatal("error running migrations", err)
	}

	defer sqlDB.Close()

	store := store.NewStorage(gormDB)

	app := &application{
		config: cfg,

		store: store,
	}

	mux := app.mount()
	log.Println("✅ Database connected successfully")
	log.Println("🚀 Server starting on", cfg.addr)
	log.Fatal(app.run(mux))

}
