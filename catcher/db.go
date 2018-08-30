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

func (c *Catcher) deleteOldRequests() error {
	const q = `
		DELETE requests, request_headers
		FROM requests
		LEFT JOIN request_headers ON requests.id = request_headers.request_id
		WHERE requests.when < NOW() - INTERVAL 1 DAY;
	`
	_, err := c.db.Db.Exec(q)
	return err
}
