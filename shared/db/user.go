package shared

import "gopkg.in/gorp.v1"

type User struct {
	Id int64
	Name string
	Password string
	Email string
	Created int64
	Modified int64
}

func UserSetup(dbmap *gorp.DbMap ){
	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(User{}, "user").SetKeys(true, "Id")
}
