//Package config of WarmAPI
package config

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
)

// DB represents the Mongo Database
var DB *mgo.Database

// DBName represents the name of the Database
var DBName = "warmDB"

// DBTestName represents the test database for running go tests.
var DBTestName = "warmDB-TEST"

/**
// Collections
var Requests *mgo.Collection
var Responses *mgo.Collection
*/

func init() {
	fmt.Println("************")
	fmt.Println("FIRING IT UP -> " + DBName)
	fmt.Println("************")
	// get a mongo session
	s, err := mgo.Dial("mongodb://127.0.0.1:27017/" + DBName)
	if err != nil {
		panic(err)
	}

	if err = s.Ping(); err != nil {
		panic(err)
	}

	DB = s.DB(DBName)
	//DB = s.DB(DBTestName)

	if DB.Name == DBTestName {
		DB.DropDatabase()
	}

	log.Println("[Database] You connected to your mongo database.")
}

/*
*	Public Functions
 */

// GetCollection function that returns a collection in the database,
// that corresponds to the name inside the database, passed in the string.
func GetCollection(name string) *mgo.Collection {
	return DB.C(name)
}
