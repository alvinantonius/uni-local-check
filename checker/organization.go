package checker

import (
	"log"

	"github.com/alvinantonius/acl-check/conn/cassandra"
)

var orgMap map[string]Organization

func init() {
	orgMap = make(map[string]Organization)
}

func GetOrg(orgID string) Organization {

	org, ok := orgMap[orgID]
	if ok {
		return org
	}

	var name, externalID string
	var aclToken, groupToken []byte
	var buyers, publishers []string
	var isDMP, isBuyer, isPublisher bool

	err := cassandra.GetSession().
		Query(`SELECT name, eyeota_id, acl_token, dmp_buyers, dmp_publishers, group_tokens, isdmp, isbuyer, ispublisher FROM organizations WHERE id = ?`, orgID).
		Scan(&name, &externalID, &aclToken, &buyers, &publishers, &groupToken, &isDMP, &isBuyer, &isPublisher)
	if err != nil {
		log.Println(err)
	}

	org = Organization{
		ID:          orgID,
		Name:        name,
		ExternalID:  externalID,
		ACLToken:    aclToken,
		GroupToken:  groupToken,
		Buyers:      buyers,
		Publisher:   publishers,
		IsDMP:       isDMP,
		IsBuyer:     isBuyer,
		IsPublisher: isPublisher,
	}

	orgMap[orgID] = org

	return org
}

func (o *Organization) HasBuyer(id string) bool {
	return checkSlice(o.Buyers, id)
}

func (o *Organization) HasPublisher(id string) bool {
	return checkSlice(o.Publisher, id)
}

func checkSlice(toCheck []string, entry string) bool {
	for _, val := range toCheck {
		if val == entry {
			return true
		}
	}
	return false
}
