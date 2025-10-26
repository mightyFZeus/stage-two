package db

import (
	"context"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(connStr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*gorm.DB, error) {
	// Connect using MySQL driver
	gormDB, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxIdleTime(duration)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)

	// Run migrations or initialization SQL
	if err := migrate(gormDB); err != nil {
		return nil, err
	}

	// Check DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return gormDB, nil
}

func migrate(db *gorm.DB) error {
	// MySQL doesn’t support CREATE EXTENSION — skip it or add your own SQL setup
	statements := []string{
		// Example: set default charset or create initial table manually
		// `ALTER DATABASE stage_one CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;`,
	}

	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			log.Printf("Error executing statement: %s\n%v\n", stmt, err)
			return err
		}
	}

	return nil
}
