package main

type worklog struct {
	ID                  int                     `json:"id"`
	Comment             string                  `json:"comment"`
	Self                string                  `json:"self"`
	Issue               worklogIssue            `json:"issue"`
	TimeSpentSeconds    int                     `json:"timeSpentSeconds"`
	BilledSeconds       int                     `json:"billedSeconds"`
	DateStarted         string                  `json:"dateStarted"`
	Author              worklogAuthor           `json:"author"`
	WorkAttributeValues []worklogAttributeValue `json:"workAttributeValues"`
}

type worklogIssue struct {
	Key                      string `json:"key"`
	ID                       int    `json:"id"`
	Self                     string `json:"self"`
	RemainingEstimateSeconds int    `json:"remainingEstimateSeconds"`
	Summary                  string `json:"summary"`
	// IssueType                worklogIssueType `json:"issueType"`
	ProjectID int `json:"projectId"`
}

type worklogIssueType struct {
	Name    string `json:"name"`
	IconURL string `json:"iconUrl"`
}

type worklogAuthor struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar"`
	Self        string `json:"self"`
}

type worklogAttributeValue struct {
	Value         string                         `json:"value,omitempty"`
	ID            int                            `json:"id,omitempty"`
	WorkAttribute worklogAttributeValueAttribute `json:"workAttribute,omitempty"`
	WorklogID     int                            `json:"worklogId,omitempty"`
}

type worklogAttributeValueAttribute struct {
	Name        string                             `json:"name,omitempty"`
	Key         string                             `json:"key,omitempty"`
	ID          int                                `json:"id,omitempty"`
	Type        worklogAttributeValueAttributeType `json:"type,omitempty"`
	Required    bool                               `json:"required,omitempty"`
	Sequence    int                                `json:"sequence,omitempty"`
	ExternalURL string                             `json:"externalUrl,omitempty"`
}

type worklogAttributeValueAttributeType struct {
	Name       string `json:"name,omitempty"`
	Value      string `json:"value,omitempty"` // "ACCOUNT","BILLABLE_SECONDS","CHECKBOX","DYNAMIC_DROPDOWN","INPUT_FIELD","INPUT_NUMERIC","SCRIPT"
	SystemType bool   `json:"systemType,omitempty"`
}
