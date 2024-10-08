package cache

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/charmbracelet/log"
	bolt "go.etcd.io/bbolt"
)

const (
	// Buckets
	PAGE_DATA    = "PageData"
	CALENDAR     = "Calendar"
	USER_DATA    = "UserData"
	LEADERBOARDS = "Leaderboards"
	PUZZLES      = "Puzzles"

	// Sub Buckets
	USER_INPUTS = "UserInputs"
	USER_PAGES  = "UserPages"

	// Other
	GENERIC_USER = "GenericUser"
)

var (
	UserCacheDir, _ = os.UserCacheDir()
	CacheDir        = path.Join(UserCacheDir, "aocgo")
	CacheFile       = path.Join(CacheDir, "%v.db")
	InputCacheDir   = path.Join(CacheDir, "inputs")
	GeneralCacheDB  = fmt.Sprintf(CacheFile, GENERIC_USER)
)

// Interface for storable resource
type Resource interface {
	GetID() string                // ID is used as key for storage
	GetBucketName() string        // Returns the name of the bucket the resource is stored in
	MarshalData() ([]byte, error) // Returns the resources data in a savable format
	SaveResource()
}

var masterDBM *DatabaseManager

// Create and initialize master database manager, taking in a valid AoC user session token
func StartupDBM(userSession string) error {
	masterDBM = &DatabaseManager{}
	return masterDBM.initializeDBM(userSession)
}

// Ensure Master DBM gets shutdown
func ShutdownDBM() {
	masterDBM.Shutdown()
}

// Database Manager
type DatabaseManager struct {
	sessionDB *bolt.DB
	// generalDB    *bolt.DB
	saveFilePath string
}

// Initializes the DBM
func (dbm *DatabaseManager) initializeDBM(userSession string) error {
	// log.Debug("---Initializing Database---")

	// Load save file path and ensure it exists
	dbm.saveFilePath = fmt.Sprintf(CacheFile, userSession)
	os.MkdirAll(path.Join(CacheDir), os.ModePerm)

	// log.Debugf("Trying to access save file path: %v", dbm.saveFilePath)

	// Open database. Read/Write for user, none for Group/Other, and none for Gretchen Weiners
	tempDB, err := bolt.Open(dbm.saveFilePath, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return err
	}
	dbm.sessionDB = tempDB
	// log.Debug("Session database opened")

	// log.Debugf("Trying to access save file path: %v", GeneralCacheDB)

	// // Open database. Read/Write for user, none for Group/Other, and none for Gretchen Weiners
	// tempDB, err = bolt.Open(GeneralCacheDB, 0600, &bolt.Options{Timeout: 10 * time.Second})
	// if err != nil {
	// 	return err
	// }
	// dbm.generalDB = tempDB
	// log.Debug("General database opened")
	//
	dbm.initializeBuckets()
	// log.Debug("Buckets initialized")
	return nil
}

// Ensure all buckets exist so they can assuredly be loaded later on
func (dbm *DatabaseManager) initializeBuckets() {
	dbm.sessionDB.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(PAGE_DATA))
		tx.CreateBucketIfNotExists([]byte(PUZZLES))
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
	// log.Debug("Database closed")
}

func SaveResource(r Resource) {
	if masterDBM == nil {
		return
	}
	// log.Debug("Saving resource", "bucket", r.GetBucketName(), "id", r.GetID(), "data", resourceData)
	masterDBM.sessionDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(r.GetBucketName()))
		resourceData, err := r.MarshalData()
		if err != nil {
			return err
		}
		bucket.Put([]byte(r.GetID()), resourceData)
		return nil
	})
}

// Save resource to database
func SaveGenericResource(bucketName, idToSave string, dataToSave []byte) {
	// log.Debug("Saving resource", "bucket", bucketName, "id", idToSave, "data", dataToSave)
	masterDBM.sessionDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		bucket.Put([]byte(idToSave), dataToSave)
		return nil
	})

}

// Load resource from database by ID
func LoadResource(bucketName, idToLoad string) []byte {
	if masterDBM == nil {
		return nil
	}

	var output []byte
	masterDBM.sessionDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		output = bucket.Get([]byte(idToLoad))
		return nil
	})
	// log.Debug("Loading resource", "bucket", bucketName, "id", idToLoad, "data", output)
	return output
}

// Clear database file for a certain user
func ClearUserDatabase(sessionToken string) {
	os.Remove(fmt.Sprintf(CacheFile, sessionToken))
}

func checkErr(err error) {
	if err != nil {
		log.Error("Database error!", "err", err)
	}
}
