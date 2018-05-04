package checker

type (
	Campaign struct {
		ID             string   `json:"id"`
		Segments       []string `json:"segments"`
		OrganizationID string   `json:"organization_id"`
		Active         bool     `json:"active"`
		Piggyback      string   `json:"piggyback"`
		BuyerID        string   `json:"buyerid"`
		Name           string   `json:"name"`
		Deleted        bool     `json:"deleted"`
		PlatformID     string   `json:"platform_id"`
	}

	Segment struct {
		ID             string `json:"id"`
		OrganizationID string `json:"organization_id"`
		ParentID       string `json:"parent_id"`
		FullName       string `json:"fullname"`
	}

	Mapping struct {
		ID             string `json:"id"`
		OrganizationID string `json:"organization_id"`
		SegmentID      string `json:"segment_id"`
		SegmentName    string `json:"segment_name"`
		Parameter      string `json:"parameter"`
		Value          string `json:"value"`
		Ignore         string `json:"ignore"`
	}

	Organization struct {
		ID          string
		Name        string
		ExternalID  string
		ACLToken    []byte
		GroupToken  []byte
		Buyers      []string
		Publisher   []string
		IsBuyer     bool
		IsDMP       bool
		IsPublisher bool
		Tags        []string
	}
)
