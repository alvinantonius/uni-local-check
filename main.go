package main

import (
	"github.com/alvinantonius/acl-check/checker"
	"github.com/alvinantonius/acl-check/conn/cassandra"
	"github.com/alvinantonius/acl-check/conn/elastic"
)

const (
	cassandraHost = "127.0.0.1"
	keyspace      = "eyeota_uni"
	elasticHost   = "http://127.0.0.1:9200"
)

func init() {
	elastic.Connect([]string{elasticHost})
	cassandra.Connect([]string{cassandraHost}, keyspace)
}

func main() {
	checker.CheckCampaignsOwnership()
}
