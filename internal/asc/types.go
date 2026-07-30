package asc

import "encoding/json"

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
		Health     string            `json:"health"`
		NextAction string            `json:"nextAction"`
		Blockers   []json.RawMessage `json:"blockers"`
	} `json:"summary"`
	AppStore struct {
		CreatedDate string `json:"createdDate"`
		Platform    string `json:"platform"`
		State       string `json:"state"`
		Version     string `json:"version"`
		VersionID   string `json:"versionId"`
	} `json:"appstore"`
	Submission struct {
		InFlight       bool              `json:"inFlight"`
		BlockingIssues []json.RawMessage `json:"blockingIssues"`
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

type reviewStatusResponse struct {
	ReviewDetailConfigured bool `json:"reviewDetailConfigured"`
	Version                struct {
		ID          string `json:"id"`
		Version     string `json:"version"`
		Platform    string `json:"platform"`
		State       string `json:"state"`
		CreatedDate string `json:"createdDate"`
	} `json:"version"`
	LatestSubmission struct {
		ID            string `json:"id"`
		State         string `json:"state"`
		Platform      string `json:"platform"`
		SubmittedDate string `json:"submittedDate"`
	} `json:"latestSubmission"`
	ReviewState string `json:"reviewState"`
	NextAction  string `json:"nextAction"`
}

type reviewHistoryResponse []struct {
	SubmissionID string `json:"submissionId"`
	Version      string `json:"versionString"`
	Platform     string `json:"platform"`
	State        string `json:"state"`
	Outcome      string `json:"outcome"`
	SubmittedAt  string `json:"submittedDate"`
}

type versionViewResponse struct {
	BuildID string `json:"buildId"`
}

type buildInfoResponse struct {
	Data struct {
		Attributes struct {
			ProcessingState string `json:"processingState"`
		} `json:"attributes"`
	} `json:"data"`
}
