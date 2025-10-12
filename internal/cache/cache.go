package cache

import (
	"fmt"
	"os"
	"path"
	"time"

	"go.dalton.dog/aocgo/internal/output"
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
	dbm := &DatabaseManager{}
	if err := dbm.initializeDBM(userSession); err != nil {
		return err
	}
	masterDBM = dbm
	return nil
}

// Ensure Master DBM gets shutdown
func ShutdownDBM() {
	if masterDBM == nil {
		return
	}
	masterDBM.Shutdown()
	masterDBM = nil
}

// Database Manager
type DatabaseManager struct {
	saveFilePath string
}

// Initializes the DBM
func (dbm *DatabaseManager) initializeDBM(userSession string) error {
	// log.Debug("---Initializing Database---")

	// Load save file path and ensure it exists
	dbm.saveFilePath = fmt.Sprintf(CacheFile, userSession)
	if err := os.MkdirAll(path.Join(CacheDir), os.ModePerm); err != nil {
		return err
	}

	// log.Debugf("Trying to access save file path: %v", dbm.saveFilePath)

	return dbm.initializeBuckets()
}

// Ensure all buckets exist so they can assuredly be loaded later on
func (dbm *DatabaseManager) initializeBuckets() error {
	return dbm.withDB(false, func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			buckets := [][]byte{
				[]byte(PAGE_DATA),
				[]byte(PUZZLES),
				[]byte(USER_INPUTS),
				[]byte(USER_DATA),
				[]byte(LEADERBOARDS),
			}

			for _, bucket := range buckets {
				if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
					return err
				}
			}

			return nil
		})
	})
}

// Ensure database is properly closed
func (dbm *DatabaseManager) Shutdown() {
	// No persistent handles remain; kept for API compatibility.
}

func SaveResource(r Resource) {
	if masterDBM == nil {
		return
	}
	// log.Debug("Saving resource", "bucket", r.GetBucketName(), "id", r.GetID(), "data", resourceData)
	err := masterDBM.withDB(false, func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(r.GetBucketName()))
			if bucket == nil {
				return fmt.Errorf("bucket %s does not exist", r.GetBucketName())
			}
			resourceData, err := r.MarshalData()
			if err != nil {
				return err
			}
			return bucket.Put([]byte(r.GetID()), resourceData)
		})
	})
	checkErr(err)
}

// Save resource to database
func SaveGenericResource(bucketName, idToSave string, dataToSave []byte) {
	if masterDBM == nil {
		return
	}
	// log.Debug("Saving resource", "bucket", bucketName, "id", idToSave, "data", dataToSave)
	err := masterDBM.withDB(false, func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketName))
			if bucket == nil {
				return fmt.Errorf("bucket %s does not exist", bucketName)
			}
			return bucket.Put([]byte(idToSave), dataToSave)
		})
	})
	checkErr(err)
}

// Load resource from database by ID
func LoadResource(bucketName, idToLoad string) []byte {
	if masterDBM == nil {
		return nil
	}

	var output []byte
	err := masterDBM.withDB(true, func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketName))
			if bucket == nil {
				return nil
			}

			value := bucket.Get([]byte(idToLoad))
			if value != nil {
				output = make([]byte, len(value))
				copy(output, value)
			}
			return nil
		})
	})
	checkErr(err)
	// log.Debug("Loading resource", "bucket", bucketName, "id", idToLoad, "data", output)
	return output
}

// Clear database file for a certain user
func ClearUserDatabase(sessionToken string) {
	os.Remove(fmt.Sprintf(CacheFile, sessionToken))
}

func checkErr(err error) {
	if err != nil {
		output.Error("Database error!", "err", err)
	}
}

// withDB opens the database, executes the provided function, and ensures the file handle is released.
func (dbm *DatabaseManager) withDB(readOnly bool, fn func(*bolt.DB) error) error {
	if dbm == nil || dbm.saveFilePath == "" {
		return fmt.Errorf("database manager not initialized")
	}

	options := &bolt.Options{
		Timeout:  10 * time.Second,
		ReadOnly: readOnly,
	}

	db, err := bolt.Open(dbm.saveFilePath, 0600, options)
	if err != nil {
		return err
	}
	defer db.Close()

	return fn(db)
}
