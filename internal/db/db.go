package db

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Timer struct {
	StartTime time.Time
	StopTime  time.Time
	Total     time.Duration
}

var (
	db   *gorm.DB
	onec sync.Once
)

func GetDB() *gorm.DB {
	onec.Do(func() {
		initializePostgress()
	})
	return db
}

func initializePostgress() {
	err := godotenv.Load("../internal/db/dbPostgress.env")
	if err != nil {
		fmt.Println("Err loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	db = connection
}

func Start(db *gorm.DB, start time.Time) error {
	if (!db.Migrator().HasTable(&Timer{})) {
		err := db.Migrator().CreateTable(&Timer{})
		if err != nil {
			return err
		}
	}

	timer := Timer{}
	tx := db.Last(&timer)
	if tx.Error != nil {
		res := db.Create(&Timer{StartTime: start})
		if res.Error != nil {
			return res.Error
		}
	}

	dateLastRecord := strings.Split(timer.StartTime.String(), " ")[0]
	dateNow := strings.Split(time.Now().String(), " ")[0]

	fmt.Println("dateLastRecord : ", dateLastRecord)
	fmt.Println("dateNow : ", dateNow)

	if dateLastRecord == dateNow {
		timer.StartTime = time.Now()
		tx = db.Model(&Timer{}).Where("stop_time = ? ", timer.StopTime).Updates(timer)
		if tx.Error != nil {
			return tx.Error
		}
	} else {
		res := db.Create(&Timer{StartTime: start})
		if res.Error != nil {
			return res.Error
		}
	}

	return nil
}

func Stop(db *gorm.DB) error {
	timer := Timer{}
	tx := db.Last(&timer)
	if tx.Error != nil {
		return tx.Error
	}
	totalTime := timer.Total
	timer.StopTime = time.Now()
	timer.Total = timer.StopTime.Sub(timer.StartTime)
	timer.Total += totalTime
	tx = db.Model(&Timer{}).Where("start_time = ? ", timer.StartTime).Updates(timer)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func TotalTime(db *gorm.DB) (error, string) {
	timers := []Timer{}
	tx := db.Find(&timers)
	if tx.Error != nil {
		return tx.Error, "There is no record to show ."
	}

	var result string
	result = "---- Date ---------  Hours  ------  Miniutes  ----\n"
	for i := len(timers) - 1; i > 0; i-- {
		date := strings.Split(timers[i].StartTime.String(), " ")
		dayWorkTime := fmt.Sprintf(" [*] %s ---  %.2f  -------  %.2f\n", date[0], timers[i].Total.Hours(), timers[i].Total.Minutes())
		result += dayWorkTime
	}
	return nil, result
}

func Clear(db *gorm.DB) error {
	tx := db.Where("total != ?", 0).Delete(&Timer{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
