package checker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	es "github.com/alvinantonius/acl-check/conn/elastic"

	"gopkg.in/olivere/elastic.v3"
)

var systemContext = []byte{0x00, 0x00, 0x00, 0x00, 0x3f, 0xb1, 0xf8, 0x05, 0x1f, 0x59, 0x40, 0xd4, 0xb1, 0xfa, 0x8f, 0xd6, 0xeb, 0x01, 0xc8, 0x91, 0x00, 0x07}

const EyeotaOrgID = "78dd4a51-4d4e-431d-bb94-c422a2ede16f"
const ESIndex = "uni"

func CheckByMappings() {

	// Get all mappings
	mappings := GetAllMappings()

	fmt.Printf("Mapping ID,mapping OrgID,Org Name,Segment Name,Ignore,Parameter\n")

	// from all mappings, get the segment's org_id
	for _, mapping := range mappings {
		// for each segment's org_id, make sure mapping's org_id is on publisher list
		s := GetSegment(mapping.SegmentID)
		if s.OrganizationID == mapping.OrganizationID {
			continue
		}
		sOrg := GetOrg(s.OrganizationID)
		if sOrg.HasPublisher(mapping.OrganizationID) == false {
			morg := GetOrg(mapping.OrganizationID)
			if morg.IsDMP || morg.IsPublisher || morg.IsBuyer {
				fmt.Printf("%v, %v, %v, %v, %v, %v\n",
					mapping.ID, mapping.OrganizationID, morg.Name, mapping.SegmentName, mapping.Ignore, mapping.Parameter)
			}
		}
	}
}

func CheckByCampaigns() {
	// Get all campaigns
	campaigns := GetAllCampaigns()

	fmt.Printf("Campaign ID,Campaign OrgID,Buyer Name,Publisher Name,Campaign Name,Active,BuyerID,Piggyback\n")

	// from all campaigns, get the segment's org_id
	for _, c := range campaigns {
		// for each segment's org_id, make sure campaign's org_id is on buyer list
		for _, sID := range c.Segments {
			s := GetSegment(sID)
			if s.OrganizationID == c.OrganizationID {
				continue
			}
			sOrg := GetOrg(s.OrganizationID)
			if sOrg.HasBuyer(c.OrganizationID) == false {
				corg := GetOrg(c.OrganizationID)
				if corg.IsDMP || corg.IsPublisher || corg.IsBuyer {
					fmt.Printf("%v, %v, %v, %v, %v, %v, %v, %v\n",
						c.ID, c.OrganizationID, corg.Name, c.Name, sOrg.Name, c.Active, c.BuyerID, c.Piggyback)

				}
			}
		}
	}
}

func CampaignsNoPublisher() {
	campaigns := GetAllCampaigns()

	var data [][]string

	header := []string{"Campaign Org ID", "Campaign ID", "Campaign Name", "Deleted"}
	data = append(data, header)

	for _, campaign := range campaigns {
		if len(campaign.Segments) == 0 {
			data = append(data, []string{campaign.OrganizationID, campaign.ID, campaign.Name, strconv.FormatBool(campaign.Deleted)})
			continue
		}

		var publisherID []string
		var ignore bool
		for _, sID := range campaign.Segments {
			segment := GetSegment(sID)
			if segment.ParentID == GeoIPRoot {
				if len(campaign.Segments) == 1 {
					ignore = true
				}
				continue
			}

			if segment.OrganizationID == "" {
				continue
			}

			org := GetOrg(segment.OrganizationID)
			if org.ID == "" {
				continue
			}

			publisherID = append(publisherID, org.ID)
		}

		if ignore {
			continue
		}

		if len(publisherID) == 0 {
			data = append(data, []string{campaign.OrganizationID, campaign.ID, campaign.Name, strconv.FormatBool(campaign.Deleted)})
		}
	}

	ToCSV("./camp-no-publisher.csv", data)
}

func CampaignsWithCrossOrgSegments() {
	campaigns := GetAllCampaigns()

	var data [][]string

	header := []string{"Campaign Org ID", "Campaign ID", "Campaign Name", "Segments Full Name", "Segment's Org Name"}
	data = append(data, header)

	for _, campaign := range campaigns {
		lastSegmentOrgID := ""
		SegmentsOrgCount := 0
		for _, sID := range campaign.Segments {
			if sID == GeoIPRoot {
				continue
			}
			segment := GetSegment(sID)
			if segment.ParentID == GeoIPRoot {
				continue
			}

			if lastSegmentOrgID != segment.OrganizationID {
				SegmentsOrgCount++
				lastSegmentOrgID = segment.OrganizationID
			}
		}

		if SegmentsOrgCount > 1 {
			campaignPrinted := false
			for _, sID := range campaign.Segments {
				segment := GetSegment(sID)
				sOrg := GetOrg(segment.OrganizationID)

				var row []string

				if segment.ParentID == GeoIPRoot {
					continue
				}
				if campaignPrinted == false {
					row = []string{campaign.OrganizationID, campaign.ID, campaign.Name}
					campaignPrinted = true
				} else {
					row = []string{"", "", ""}

				}

				row = append(row, segment.FullName, sOrg.Name)
				data = append(data, row)
			}
		}
	}

	ToCSV("./result.csv", data)
}

// GetAllMappings is to get all undeleted mappings
func GetAllMappings() (maps []Mapping) {
	query := elastic.NewMatchQuery("deleted", false)
	res, err := es.GetClient().Search().
		Index(ESIndex).
		Query(query).
		Type("mappings").
		Do()
	if err != nil {
		log.Println(err)
	}

	total := res.Hits.TotalHits

	scroller := es.GetClient().Scroll().
		Index(ESIndex).
		Query(query).
		Type("mappings").
		Size(int(total))

	for {
		res, err := scroller.Do()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
		}
		for _, hit := range res.Hits.Hits {
			var m Mapping
			json.Unmarshal(*hit.Source, &m)
			maps = append(maps, m)
		}
	}

	return maps
}

// GetAllCampaigns is to get all undeleted mappings
func GetAllCampaigns() (campaigns []Campaign) {
	// query := elastic.NewMatchQuery("deleted", false)
	res, err := es.GetClient().Search().
		Index(ESIndex).
		// Query(query).
		Type("campaigns").
		Do()
	if err != nil {
		log.Println(err)
	}

	total := res.Hits.TotalHits

	scroller := es.GetClient().Scroll().
		Index(ESIndex).
		// Query(query).
		Type("campaigns").
		Size(int(total))

	for {
		res, err := scroller.Do()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
		}
		for _, hit := range res.Hits.Hits {
			var c Campaign
			json.Unmarshal(*hit.Source, &c)
			campaigns = append(campaigns, c)
		}
	}

	return campaigns
}

func UniqueOrgHasCampaign() {
	campaigns := GetAllCampaigns()
	type orgInfo struct {
		Name         string
		TotalCamp    int64
		ActiveCamp   int64
		InactiveCamp int64
		DeletedCamp  int64
	}
	orgList := make(map[string]orgInfo)

	for _, campaign := range campaigns {
		var active, deleted, inactive int64
		if campaign.Deleted {
			deleted = 1
		} else {
			if campaign.Active {
				active = 1
			} else {
				inactive = 1
			}
		}

		if o, ok := orgList[campaign.OrganizationID]; ok {
			o.TotalCamp = o.TotalCamp + 1
			o.ActiveCamp = o.ActiveCamp + active
			o.DeletedCamp = o.DeletedCamp + deleted
			o.InactiveCamp = o.InactiveCamp + inactive

			orgList[campaign.OrganizationID] = o
			continue
		}

		org := GetOrg(campaign.OrganizationID)
		if org.ID == "" {
			continue
		}
		orgList[org.ID] = orgInfo{
			Name:         org.Name,
			TotalCamp:    1,
			ActiveCamp:   active,
			InactiveCamp: inactive,
			DeletedCamp:  deleted,
		}
	}

	var data [][]string
	header := []string{"Org ID", "Org Name", "Total", "Active", "Inactive", "Deleted"}
	data = append(data, header)
	for id, org := range orgList {
		data = append(data, []string{id, org.Name, strconv.FormatInt(org.TotalCamp, 10), strconv.FormatInt(org.ActiveCamp, 10), strconv.FormatInt(org.InactiveCamp, 10), strconv.FormatInt(org.DeletedCamp, 10)})
	}

	ToCSV("all-campaigns-old-org.csv", data)
}

func CheckCampaignsOwnership() {
	fmt.Println("Check Campaigns Ownership")
	campaigns := GetAllCampaigns()

	/*
		Check:
		- campaigns without segment
		- campaigns that has multiple segments (exclude geoip) is
		  own by Eyeota and Eyeota should be under `dmp_buyers` list in
		  those segments owner organizations
		  	+ need to check group token
		-
	*/

	fmt.Println("Total Campaigns :", len(campaigns))

	for _, c := range campaigns {
		// check platform_id shouldn't be null
		if c.PlatformID == "" {
			fmt.Println("No platform:", c.ID)
			// continue
		}

		if c.PlatformID == c.OrganizationID {
			fmt.Println("Org and platform id is similar. CampaignID :", c.ID)
		}

		// check if platform is buyer or not?
		platform := GetOrg(c.PlatformID)
		if !platform.IsBuyer {
			fmt.Println("platform-id is not buyer org. CID :", c.ID, "OrgName :", platform.Name, "deleted :", c.Deleted)
		}

		// check campaigns should be owned by one of segments org
		var targetOrgID string
		mapOrgID := make(map[string]bool)
		for _, sID := range c.Segments {
			s := GetSegment(sID)

			// get all org from segments that is not geoip segment
			if s.ParentID != GeoIPRoot {
				mapOrgID[s.OrganizationID] = true
			}
		}

		if len(mapOrgID) == 1 {
			for orgID := range mapOrgID {
				targetOrgID = orgID
			}
		} else {
			targetOrgID = EyeotaOrgID
		}

		if targetOrgID == "" {
			fmt.Println("can't find campaign owner:", c.ID)
		}
		if targetOrgID != c.OrganizationID {
			fmt.Println("wrong campaign owner:", c.ID)
		}

		IsValidCampaignsTokens(c)
	}
}

func IsValidCampaignsTokens(c Campaign) {
	org := GetOrg(c.OrganizationID)

	if len(org.GroupToken)%20 != 0 {
		fmt.Println("invalid group token. OrgID:", org.ID)
	}
	var orgGroupTokens [][]byte
	for i := 0; i < len(org.GroupToken); i = i + 20 {
		orgGroupTokens = append(orgGroupTokens, org.GroupToken[i:i+20])
	}

	for _, sID := range c.Segments {
		s := GetSegment(sID)
		sORg := GetOrg(s.OrganizationID)

		// only check for private exchange and group token if segment's org is different than campaign's org
		if s.OrganizationID != c.OrganizationID && s.ParentID != GeoIPRoot {
			var isInBuyerList bool
			// check private exchange
			for _, buyerOrgID := range sORg.Buyers {
				if buyerOrgID == org.ID {
					isInBuyerList = true
					break
				}
			}
			if !isInBuyerList {
				fmt.Println("not in buyer list. OrgID:", org.ID, "segmentID:", s.ID)
			}

			// check group token
			buyerACLToken := make([]byte, len(sORg.ACLToken))
			copy(buyerACLToken, sORg.ACLToken)
			buyerACLToken[3] = 0x02

			var isGroupTokenValid bool
			for _, token := range orgGroupTokens {
				if bytes.Equal(buyerACLToken, token) {
					isGroupTokenValid = true
					break
				}
			}
			if !isGroupTokenValid {
				fmt.Println("group token no access. OrgID:", org.ID, "segmentID:", s.ID)
			}
		}

	}

	// validate campaign's token
	tokens := c.GetToken()
	if len(tokens)%22 != 0 {
		fmt.Println("invalid campaign token. campaignID:", c.ID)
	}
	var cTokens [][]byte
	for i := 0; i < len(tokens); i = i + 22 {
		cTokens = append(cTokens, tokens[i:i+22])
	}
	tokenCopy := cTokens

	// campaign's token should contains system-context and org token only
	orgToken := append(org.ACLToken, []byte{0x00, 0x07}...)
	var systemTokenFound, orgTokenFound bool
	for i, t := range cTokens {
		if bytes.Equal(t, systemContext) {
			systemTokenFound = true
			tokenCopy[i] = nil
		}
		if bytes.Equal(t, orgToken) {
			orgTokenFound = true
			tokenCopy[i] = nil
		}
	}
	if !systemTokenFound {
		fmt.Println("no system-context. campaignID:", c.ID)
	}
	if !orgTokenFound {
		fmt.Println("no org-token. campaignID:", c.ID)
	}

	for _, t := range tokenCopy {
		if t != nil {
			fmt.Println("campaign has additional token. campaignID:", c.ID)
		}
	}
}

func CompleteCheck() {
	// get all mappings

	// get all segments

	/* check all mappings must be pointing to correct segments
	- using private exchange properly
	- segments is not deleted */

	// check all linked segments must go through private exchange

	// check segments

	/* check campaigns
	- no campaigns pointing to only geoip
	- no campaigns pointing to 2 geoip
	- check ACL
	- check private exchange
	- check campaigns pointing to deleted segments
	*/

}
