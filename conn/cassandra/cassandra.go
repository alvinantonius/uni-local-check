package cassandra

import (
	"log"

	"github.com/gocql/gocql"
)

var session *gocql.Session

func Connect(hosts []string, keyspace string) {
	var err error
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	session, err = cluster.CreateSession()
	if err != nil {
		log.Panic(err)
	}
}

func GetSession() *gocql.Session {
	return session
}
