package shared


import (
	"database/sql"
	"gopkg.in/gorp.v1"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path/filepath"
	"os"
)


func Init(){
	// initialize the DbMap
	dbmap := initDb()
	defer dbmap.Db.Close()
}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	dir, e:= os.Getwd()
	if e != nil {
		panic(e)
	}
	db, err := sql.Open("sqlite3", filepath.Join(dir, "post_db.bin"))
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	// dbmap.AddTableWithName(User{}, "user").SetKeys(true, "Id")
	UserSetup(dbmap)

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
