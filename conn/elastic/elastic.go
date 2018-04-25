package elastic

import (
	"log"

	"gopkg.in/olivere/elastic.v3"
)

var client *elastic.Client

func Connect(hosts []string) {
	var err error
	client, err = elastic.NewClient(elastic.SetURL(hosts...))
	if err != nil {
		log.Panic(err)
	}
}

func GetClient() *elastic.Client {
	return client
}
