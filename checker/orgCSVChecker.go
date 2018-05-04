package checker

import (
	"fmt"
	"strconv"
	"strings"
)

func CheckORG() {
	data := openCSV("orgs-data.csv")
	header := data[0]
	orgs := []Organization{}
	for i := 1; i < len(data); i++ {
		isPublisher, _ := strconv.ParseBool(data[i][3])
		isDMP, _ := strconv.ParseBool(data[i][4])
		isBuyer, _ := strconv.ParseBool(data[i][5])
		org := Organization{
			ID:          data[i][0],
			Name:        data[i][1],
			ExternalID:  data[i][2],
			IsDMP:       isDMP,
			IsBuyer:     isBuyer,
			IsPublisher: isPublisher,
		}

		tags := []string{}
		for j := 6; j < len(data[i]); j++ {
			if data[i][j] == "" {
				continue
			}
			subtags := strings.Split(data[i][j], ",")
			if len(subtags) > 1 {
				for tagIndex, tag := range subtags {
					if strings.Contains(tag, "_") {
						continue
					}
					tag = header[j] + "_" + strings.Trim(tag, " ")
					subtags[tagIndex] = tag
				}
			}
			tags = append(tags, subtags...)
		}
		org.Tags = tags

		// checkUpdatedFields(org)

		orgs = append(orgs, org)
	}

	orgsToCSV(orgs)
}

func orgsToCSV(orgs []Organization) {
	data := [][]string{}

	header := []string{
		"ID",
		// "Name",
		// "External ID",
		// "Is DMP",
		// "Is Buyer",
		// "IS Publisher",
		"Tags",
	}

	data = append(data, header)

	for _, org := range orgs {
		d := []string{
			org.ID,
			// org.Name,
			// org.ExternalID,
			// strconv.FormatBool(org.IsDMP),
			// strconv.FormatBool(org.IsBuyer),
			// strconv.FormatBool(org.IsPublisher),
			strings.Join(org.Tags, ","),
		}

		data = append(data, d)
	}

	ToCSV("new-org-data.csv", data)
}

func checkUpdatedFields(org Organization) {
	oldOrg := GetOrg(org.ID)

	if oldOrg.Name != org.Name {
		fmt.Printf("%s,%s\n", oldOrg.Name, org.Name)
	}
	if oldOrg.IsBuyer != org.IsBuyer {
		fmt.Print("isbuyer is diff")
	}
	if oldOrg.IsDMP != org.IsDMP {
		fmt.Print("isdmp is diff")
	}
	if oldOrg.IsPublisher != org.IsPublisher {
		fmt.Print("is-pub is diff")
	}
}
