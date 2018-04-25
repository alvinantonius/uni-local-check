package checker

import (
	"log"

	"github.com/alvinantonius/acl-check/conn/cassandra"
)

func (c *Campaign) GetToken() []byte {
	var tokens []byte
	err := cassandra.GetSession().
		Query(`SELECT tokens FROM campaigns WHERE id = ?`, c.ID).
		Scan(&tokens)
	if err != nil {
		log.Println(err)
	}
	return tokens
}
