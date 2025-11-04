package config

import (
	"fmt"
)

func MigrateDB() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	);
	`

	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		sku_name VARCHAR(100) NOT NULL,
		quantity INT DEFAULT 0
	);
	`

	createLocationsTable := `
	CREATE TABLE IF NOT EXISTS locations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		code VARCHAR(50) NOT NULL UNIQUE,
		name VARCHAR(100) NOT NULL,
		capacity INT NOT NULL
	);
	`

	createStockMovementsTable := `
	CREATE TABLE IF NOT EXISTS stock_movements (
		id INT AUTO_INCREMENT PRIMARY KEY,
		product_id INT NOT NULL,
		location_id INT NOT NULL,
		type ENUM('IN','OUT') NOT NULL,
		quantity INT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
		FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE
	);
	`

	queries := []string{
		createUsersTable,
		createProductsTable,
		createLocationsTable,
		createStockMovementsTable,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			fmt.Println("❌ Migration error:", err)
		}
	}

	fmt.Println("✅ Database migrated successfully")
}
