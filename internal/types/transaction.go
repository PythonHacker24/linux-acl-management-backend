package types

/*
	contains shared definations where compete modulation was not possible
	Eg. session and transprocesser need same transaction structure and updating seperate definations
	needs rewriting same code multiple times.
*/

/* transaction represents a permission management transaction */
type Transaction struct {
	ID          string
	UserID      string
	Action      string
	Resource    string
	Permissions string
	Status      string
	Timestamp   string
}

/* represents the result of a processed transaction */
type TransactionResult struct {
	Transaction Transaction
	Success     bool
	Error       string
	Timestamp   string
}
