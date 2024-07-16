package cockroachdb

import (
	"log"
	"ospm/config"
	"ospm/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitialDB() {
	var err error

	DB, err = gorm.Open(postgres.Open(config.OSPM.RDMS.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	// Ping the database to ensure the connection is alive
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("failed to get generic database object: ", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("failed to connect to the database: ", err)
	}

	// Create the uuid-ossp extension
	err = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Fatal("failed to create uuid-ossp extension: ", err)
	}

	log.Println("database connection established successfully")

	// Run auto migration
	err = DB.AutoMigrate(
		&models.Organization{},
		&models.OrganizationDetails{},
		&models.OrganizationOwner{},
		&models.Subscriber{},
		&models.SubscriberDetails{},
		&models.Credentials{},
		&models.SubscriberGroup{},
		&models.Permission{},
		&models.ProductOffering{},
		&models.ProductOfferingSpecification{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
}
