package main

import (
	"github.com/alvinantonius/uni-local-check/checker"
	"github.com/alvinantonius/uni-local-check/conn/cassandra"
	"github.com/alvinantonius/uni-local-check/conn/elastic"
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
	checker.CheckORG()
}
