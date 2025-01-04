package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	// Connect to database
	var err error

	//getting the database URL from the environment variable
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

}

func SequenceGeneratorCreators() {
	// Ensure sequence exists in the database
	DB.Exec("CREATE SEQUENCE IF NOT EXISTS user_otp_seq START 1")
	DB.Exec("CREATE SEQUENCE IF NOT EXISTS trip_fee_mapping_seq START 1")
	//DB.Exec("INSERT INTO vehicle_types (name)  VALUES ('TAMPOO'), ('BIKE'), ('SCOOTER'), ('CAR') ON CONFLICT DO NOTHING")
}

func InitializeValuesInDb() {
	populateVehicleType()
	populateRideType()
}

func populateVehicleType() {
	// Check if the "vehicle_types" table exists
	var tableExists bool
	if err := DB.Raw(`
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_name = 'vehicle_types'
        );
    `).Scan(&tableExists).Error; err != nil {
		log.Fatal("Failed to check table existence:", err)
	}

	// Insert records if the table exists
	if tableExists {
		if err := DB.Exec(`
            INSERT INTO vehicle_types (type)
            VALUES ('TAMPOO'), ('BIKE'), ('SCOOTER'), ('CAR')
            ON CONFLICT DO NOTHING;
        `).Error; err != nil {
			log.Fatal("Failed to insert records:", err)
		} else {
			fmt.Println("Records inserted successfully!")
		}
	} else {
		fmt.Println("Table vehicle_types does not exist.")
	}
}

func populateRideType() {
	// Check if the "vehicle_types" table exists
	var tableExists bool
	if err := DB.Raw(`
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_name = 'ride_types'
        );
    `).Scan(&tableExists).Error; err != nil {
		log.Fatal("Failed to check table existence:", err)
	}

	// Insert records if the table exists
	if tableExists {
		if err := DB.Exec(`
            INSERT INTO ride_types (type)
            VALUES ('PRIVATE'), ('PUBLIC')
            ON CONFLICT DO NOTHING;
        `).Error; err != nil {
			log.Fatal("Failed to insert records:", err)
		} else {
			fmt.Println("Records inserted successfully!")
		}
	} else {
		fmt.Println("Table vehicle_types does not exist.")
	}
}
