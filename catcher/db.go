package catcher

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/coopernurse/gorp"
	_ "github.com/ziutek/mymysql/godrv"
)

func initDb(config *Configuration) (*gorp.DbMap, error) {
	dsn := fmt.Sprintf("%s/%s/%s",
		config.Database.Name,
		config.Database.User,
		"")

	db, err := sql.Open("mymysql", dsn)
	if err != nil {
		return nil, err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(CaughtRequest{}, "requests").SetKeys(true, "ID")
	dbmap.AddTableWithName(Header{}, "request_headers").SetKeys(true, "ID")

	return dbmap, err
}

func (c *Catcher) persistRequest(request *CaughtRequest) (err error) {
	if err := c.db.Insert(request); err != nil {
		return err
	}

	for key, vals := range request.Headers {
		val := strings.Join(vals, " ")
		header := &Header{
			RequestID: request.ID,
			Key:       key,
			Value:     val,
		}
		err = c.db.Insert(header)
		if err != nil {
			return err
		}
	}
	return err
}
