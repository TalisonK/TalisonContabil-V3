package database

import (
	"context"
	"fmt"

	"github.com/TalisonK/TalisonContabil/src/config"
	"github.com/TalisonK/TalisonContabil/src/domain"
	"github.com/TalisonK/TalisonContabil/src/util/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CloudCollections struct {
	User     *mongo.Collection
	Category *mongo.Collection
	Income   *mongo.Collection
	Expense  *mongo.Collection
	List     *mongo.Collection
}

var DBlocal *gorm.DB
var DBCloud CloudCollections

// OpenConnectionLocal starts a connection with the local database
func OpenConnectionLocal() error {

	conf := config.GetLocalDB()

	fmt.Println(conf)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Pass, conf.Host, conf.Port, conf.Database)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return fmt.Errorf(logging.FailedToOpenConnection("local", err, "database.OpenConnectionLocal"))
	}

	logging.OpenedConnection("local", "database.OpenConnectionLocal")
	conn.Logger = logger.Default.LogMode(logger.Info)

	conn.AutoMigrate(&domain.Category{})
	conn.AutoMigrate(&domain.User{})
	conn.AutoMigrate(&domain.Income{})
	conn.AutoMigrate(&domain.Expense{})
	conn.AutoMigrate(&domain.List{})

	DBlocal = conn

	return nil
}

// OpenConnectionCloud starts a connection with the cloud database
func OpenConnectionCloud() error {

	conf := config.GetCloudDB()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	dsn := fmt.Sprintf("%s://%s:%s@%s/?retryWrites=true&w=majority&appName=Base-contabil", conf.Host, conf.User, conf.Pass, conf.Database)

	opts := options.Client().ApplyURI(dsn).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		logging.FailedToOpenConnection("cloud", err, "database.OpenConnectionCloud")
		return err
	}

	// enviando um ping para confirmar conexão
	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err()

	if err != nil {
		logging.FailedToConnectToDB("cloud", err, "database.OpenConnectionCloud")
		return err
	}

	DBCloud.User = client.Database("contabil").Collection("user")
	DBCloud.Category = client.Database("contabil").Collection("category")
	DBCloud.Income = client.Database("contabil").Collection("income")
	DBCloud.Expense = client.Database("contabil").Collection("expense")
	DBCloud.List = client.Database("contabil").Collection("list")

	return nil
}

// CheckLocalDB checks if the local database is connected
func checkLocalDB() bool {

	if DBlocal == nil {
		logging.FailedToPingDB("local", nil, "database.checkLocalDB")
		return false
	}

	section, err := DBlocal.DB()

	if err != nil {
		logging.FailedToConnectToDB("local", err, "database.checkLocalDB")
		return false
	}

	err = section.Ping()

	if err != nil {
		logging.FailedToPingDB("local", err, "database.checkLocalDB")
		return false
	}
	return true

}

// CheckCloudDB checks if the cloud database is connected
func checkCloudDB() bool {
	if DBCloud.Expense == nil {
		logging.FailedToPingDB("cloud", nil, "database.checkCloudDB")
		return false
	}

	err := DBCloud.Expense.FindOne(context.TODO(), bson.D{}).Err()

	if err != nil {
		logging.FailedToPingDB("cloud", err, "database.checkCloudDB")
		return false
	}
	return true
}

// CloseConnections closes the connections with the databases
func CloseConnections() {
	db, _ := DBlocal.DB()

	err := db.Close()

	if err != nil {
		logging.FailedToCloseConnection("local", err, "database.CloseConnections")
	}
}

func CheckDBStatus() (bool, bool) {
	// Check database status
	statusDbLocal := checkLocalDB()
	statusDbCloud := checkCloudDB()

	return statusDbLocal, statusDbCloud
}
