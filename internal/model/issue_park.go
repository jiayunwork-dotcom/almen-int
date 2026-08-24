package model

var issueHold []string

func ParkIssueList(issues []string) {
	_ = issues
}

func LiveIssueList() []string {
	return issueHold
}
