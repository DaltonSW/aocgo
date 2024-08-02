package cache

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	bolt "go.etcd.io/bbolt"
)

var UserCacheDir, _ = os.UserCacheDir()
var CacheDir = path.Join(UserCacheDir, "aocgo")
var CacheFile = path.Join(CacheDir, "%v.db")
var InputCacheDir = path.Join(CacheDir, "inputs")

var GeneralCacheDB = fmt.Sprintf(CacheFile, GENERIC_USER)

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
	PAGE_DATA    = "PageData"
	CALENDAR     = "Calendar"
	USER_DATA    = "UserData"
	LEADERBOARDS = "Leaderboards"

	// Sub Buckets
	USER_INPUTS = "UserInputs"
	USER_PAGES  = "UserPages"

	// Other
	GENERIC_USER = "GenericUser"
)

// Create and initialize master database manager
func Startup(userSession string) error {
	masterDBM = &DatabaseManager{}
	return masterDBM.Initialize(userSession)
}

// Ensure Master DBM gets shutdown
func Shutdown() {
	masterDBM.Shutdown()
}

// Database Manager
type DatabaseManager struct {
	sessionDB *bolt.DB
	// generalDB    *bolt.DB
	saveFilePath string
}

// Initializes the DBM
func (dbm *DatabaseManager) Initialize(userSession string) error {
	log.Debug("---Initializing Database---")

	// Load save file path and ensure it exists
	dbm.saveFilePath = fmt.Sprintf(CacheFile, userSession)
	os.MkdirAll(path.Join(CacheDir), os.ModePerm)

	log.Debugf("Trying to access save file path: %v", dbm.saveFilePath)

	// Open database. Read/Write for user, none for Group/Other, and none for Gretchen Weiners
	tempDB, err := bolt.Open(dbm.saveFilePath, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return err
	}
	dbm.sessionDB = tempDB
	log.Debug("Session database opened")

	log.Debugf("Trying to access save file path: %v", GeneralCacheDB)

	// // Open database. Read/Write for user, none for Group/Other, and none for Gretchen Weiners
	// tempDB, err = bolt.Open(GeneralCacheDB, 0600, &bolt.Options{Timeout: 10 * time.Second})
	// if err != nil {
	// 	return err
	// }
	// dbm.generalDB = tempDB
	// log.Debug("General database opened")
	//
	dbm.initializeBuckets()
	log.Debug("Buckets initialized")
	return nil
}

// Ensure all buckets exist so they can assuredly be loaded later on
func (dbm *DatabaseManager) initializeBuckets() {
	dbm.sessionDB.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(PAGE_DATA))
		tx.CreateBucketIfNotExists([]byte(USER_INPUTS))
		tx.CreateBucketIfNotExists([]byte(USER_DATA))
		tx.CreateBucketIfNotExists([]byte(LEADERBOARDS))
		return nil
	})

	// dbm.generalDB.Update(func(tx *bolt.Tx) error {
	// 	tx.CreateBucketIfNotExists([]byte(LEADERBOARDS))
	// 	tx.CreateBucketIfNotExists([]byte(USER_DATA))
	// 	return nil
	// })
}

// Ensure database is properly closed
func (dbm *DatabaseManager) Shutdown() {
	masterDBM.sessionDB.Close()
	// masterDBM.generalDB.Close()
	log.Debug("Database closed")
}

// Save resource to database
func SaveResource(bucketName, idToSave string, dataToSave []byte) {
	log.Debug("Saving resource", "bucket", bucketName, "id", idToSave, "data", dataToSave)
	masterDBM.sessionDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		bucket.Put([]byte(idToSave), dataToSave)
		return nil
	})

}

// Load resource from database by ID
func LoadResource(bucketName, idToLoad string) []byte {
	log.Debug("Loading resource", "bucket", bucketName, "id", idToLoad)
	var output []byte
	masterDBM.sessionDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		output = bucket.Get([]byte(idToLoad))
		return nil
	})
	return output
}

func SaveSubResource(parentBucket, childBucket, idToSave string, dataToSave []byte) {
	log.Debug("Saving subresource", "parent", parentBucket, "child", childBucket, "ID", idToSave)
	// var db *bolt.DB
	// if parentBucket == GENERIC_USER {
	// 	db = masterDBM.generalDB
	// } else {
	// 	db = masterDBM.sessionDB
	// }
	masterDBM.sessionDB.Update(func(tx *bolt.Tx) error {
		parent, err := tx.CreateBucketIfNotExists([]byte(parentBucket))
		if err != nil {
			return err
		}

		child, err := parent.CreateBucketIfNotExists([]byte(childBucket))
		if err != nil {
			return err
		}
		child.Put([]byte(idToSave), dataToSave)
		log.Debug("Successfully saved subresource", "data", dataToSave)
		return nil
	})

}

func LoadSubResource(parentBucket, childBucket, idToLoad string) []byte {
	log.Debug("Loading subresource", "parent", parentBucket, "child", childBucket, "ID", idToLoad)
	// var db *bolt.DB
	// if parentBucket == GENERIC_USER {
	// 	db = masterDBM.generalDB
	// } else {
	// 	db = masterDBM.sessionDB
	// }
	var resource []byte
	masterDBM.sessionDB.View(func(tx *bolt.Tx) error {
		parent := tx.Bucket([]byte(parentBucket))
		if parent == nil {
			log.Debug("Parent bucket doesn't exist")
			return errors.New("No bucket with that name")
		}

		bucket := tx.Bucket([]byte(childBucket))
		log.Debug("Child bucket", "child", bucket)
		resource = bucket.Get([]byte(idToLoad))
		log.Debug("Resource", "res", resource)
		return nil
	})
	return resource
}

func InitUser(userSession string) {
	masterDBM.sessionDB.Update(func(tx *bolt.Tx) error {
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
