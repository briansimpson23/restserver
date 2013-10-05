package api

type Member struct {
	MemberID           int64
	FirstName          string
	LastName           string
	Email              string
	ZipCode            string
	Gender             string
	FacebookID         string
	FacebookCity       string
	FacebookLocationID string
	TwitterUsername    string
	MemberStatusID     int
	MarketingSourceID  int
}
