package asc

type appListResponse struct {
	Data []appRecord `json:"data"`
}

type appRecord struct {
	ID         string        `json:"id"`
	Attributes appAttributes `json:"attributes"`
}

type appAttributes struct {
	Name     string `json:"name"`
	BundleID string `json:"bundleId"`
}

type reviewListResponse struct {
	Data []reviewRecord `json:"data"`
}

type reviewRecord struct {
	ID         string           `json:"id"`
	Attributes reviewAttributes `json:"attributes"`
}

type reviewAttributes struct {
	Platform      string `json:"platform"`
	State         string `json:"state"`
	SubmittedDate string `json:"submittedDate"`
}

type statusResponse struct {
	App struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		BundleID string `json:"bundleId"`
	} `json:"app"`
	Summary struct {
		Health     string `json:"health"`
		NextAction string `json:"nextAction"`
	} `json:"summary"`
	AppStore struct {
		CreatedDate string `json:"createdDate"`
		Platform    string `json:"platform"`
		State       string `json:"state"`
		Version     string `json:"version"`
		VersionID   string `json:"versionId"`
	} `json:"appstore"`
	Submission struct {
		InFlight bool `json:"inFlight"`
	} `json:"submission"`
	Review struct {
		LatestSubmissionID string `json:"latestSubmissionId"`
		Platform           string `json:"platform"`
		State              string `json:"state"`
		SubmittedDate      string `json:"submittedDate"`
	} `json:"review"`
	Links struct {
		AppStoreConnect string `json:"appStoreConnect"`
		Review          string `json:"review"`
		TestFlight      string `json:"testFlight"`
	} `json:"links"`
}
