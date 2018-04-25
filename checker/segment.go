package checker

import (
	"log"

	"github.com/alvinantonius/acl-check/conn/cassandra"
)

const (
	// GeoIPRoot is the root segment id for geo ip segments
	GeoIPRoot = "5c6fca6a-cc6f-4112-8c68-c4d94ce1f91e"
)

var segmentMap map[string]Segment

func init() {
	segmentMap = make(map[string]Segment)
}

func GetSegment(id string) Segment {
	seg, ok := segmentMap[id]
	if ok {
		return seg
	}

	var orgID string
	var parentID string
	var fullname string

	err := cassandra.GetSession().
		Query(`SELECT organization_id, parent_id, fullname FROM segments WHERE id = ?`, id).
		Scan(&orgID, &parentID, &fullname)
	if err != nil {
		log.Println(err)
	}

	seg = Segment{
		ID:             id,
		OrganizationID: orgID,
		ParentID:       parentID,
		FullName:       fullname,
	}

	segmentMap[id] = seg

	return seg
}
