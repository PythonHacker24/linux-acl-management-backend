package types

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
