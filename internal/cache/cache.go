package cache

import (
	"errors"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	bolt "go.etcd.io/bbolt"
)

var UserCacheDir, _ = os.UserCacheDir()
var CacheDir = path.Join(UserCacheDir, "aocutil")
var CacheFile = path.Join(CacheDir, "aocutil.db")
var InputCacheDir = path.Join(CacheDir, "inputs")

func InitCache() {
	os.MkdirAll(CacheDir, 0600)
}

func LoadUserInput(year int, day int, userSession string) []byte {
	log.Infof("Loading user puzzle input for Day %v (%v) for user %v", day, year, userSession)

	fileDir := path.Join(InputCacheDir, userSession, strconv.Itoa(year))
	filePath := path.Join(fileDir, strconv.Itoa(day)+".input")

	var file *os.File
	var err error

	file, err = os.Open(filePath)
	if err != nil {
		return []byte{}
	}
	defer file.Close()

	data, _ := io.ReadAll(file)
	log.Infof("Success")
	return data
}

func SaveUserInput(year int, day int, userSession string, input []byte) error {
	log.Infof("Saving user puzzle input for Day %v (%v) for user %v", day, year, userSession)

	fileDir := path.Join(InputCacheDir, userSession, strconv.Itoa(year))
	filePath := path.Join(fileDir, strconv.Itoa(day)+".input")

	var file *os.File
	var err error

	os.MkdirAll(fileDir, 0600)
	file, err = os.Open(filePath)
	if errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(filePath)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	defer file.Close()

	file.Write(input)
	log.Infof("Success")
	return nil
}

// Interface for storable resource
type Resource interface {
	GetID() string                // ID is used as key for storage
	GetBucketName() string        // Returns the name of the bucket the resource is stored in
	MarshalData() ([]byte, error) // Returns the resources data in a savable format
}

var masterDBM *DatabaseManager

const ( // Buckets
	PAGE_DATA = "Page Data"
	CALENDAR  = "Calendar"
	USER_DATA = "User Data"

	// Sub Buckets
	USER_INPUTS = "User Inputs"
	USER_PAGES  = "User Pages"

	// Other
	GENERIC_USER = "Generic User"

	// PAGE_DATA = "Page Data"
	// PAGE_DATA = "Page Data"
)

// Create and initialize master database manager
func Startup() {
	masterDBM = &DatabaseManager{}
	masterDBM.Initialize()
}

// Ensure Master DBM gets shutdown
func Shutdown() {
	masterDBM.Shutdown()
}

// Database Manager
type DatabaseManager struct {
	database     *bolt.DB
	saveFilePath string
}

// Initializes the DBM
func (dbm *DatabaseManager) Initialize() error {
	log.Debug("---Initializing Database---")

	// Load save file path and ensure it exists
	dbm.saveFilePath = CacheFile
	os.MkdirAll(path.Join(CacheDir), os.ModePerm)

	log.Printf("Trying to access save file path: %v\n", dbm.saveFilePath)

	// Open database. Read/Write for user, Read for group & other
	tempDB, err := bolt.Open(dbm.saveFilePath, 0644, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return err
	}
	dbm.database = tempDB
	log.Debug("Database opened")

	dbm.initializeBuckets()
	log.Debug("Buckets initialized")
	return nil
}

// Ensure all buckets exist so they can assuredly be loaded later on
func (dbm *DatabaseManager) initializeBuckets() {
	dbm.database.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(PAGE_DATA))
		tx.CreateBucketIfNotExists([]byte(CALENDAR))
		tx.CreateBucketIfNotExists([]byte(USER_DATA))
		return nil
	})
}

// Ensure database is properly closed
func (dbm *DatabaseManager) Shutdown() {
	masterDBM.database.Close()
	log.Debug("Database closed")
}

// Save resource to database
func SaveResource(resource Resource) {
	log.Printf("Saving resource: %v\n", resource)
	// fmt.Println("Saving resource", resource)
	masterDBM.database.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(resource.GetBucketName()))
		resourceData, err := resource.MarshalData()
		checkErr(err)
		bucket.Put([]byte(resource.GetID()), resourceData)
		return nil
	})

}

// Load resource from database by ID
func LoadResource(bucketName, idToLoad string) []byte {
	log.Printf("Loading resource from %v: %v\n", bucketName, idToLoad)

	// fmt.Println("Loading resource", bucketName, idToLoad)
	var output []byte
	masterDBM.database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		output = bucket.Get([]byte(idToLoad))
		return nil
	})
	return output
}

func SaveSubResource(parentBucket, childBucket, idToSave string, dataToSave []byte) {
	masterDBM.database.Update(func(tx *bolt.Tx) error {
		parent, err := tx.CreateBucketIfNotExists([]byte(parentBucket))
		if err != nil {
			return err
		}

		bucket, err := parent.CreateBucketIfNotExists([]byte(childBucket))
		if err != nil {
			return err
		}
		bucket.Put([]byte(idToSave), dataToSave)
		return nil
	})

}

func LoadSubResource(parentBucket, childBucket, idToLoad string) []byte {
	var resource []byte
	masterDBM.database.Update(func(tx *bolt.Tx) error {
		parent := tx.Bucket([]byte(parentBucket))
		if parent == nil {
			return errors.New("No bucket with that name")
		}

		bucket := tx.Bucket([]byte(childBucket))
		resource = bucket.Get([]byte(idToLoad))

		return nil
	})
	return resource
}

func InitUser(userSession string) {
	masterDBM.database.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(userSession))
		if err != nil {
			return err
		}

		bucket.CreateBucket([]byte(USER_PAGES))
		bucket.CreateBucket([]byte(USER_INPUTS))
		return nil
	})
}

func checkErr(err error) {
	if err != nil {
		log.Error("Database error!", "err", err)
	}
}
