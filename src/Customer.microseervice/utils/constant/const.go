package constant

// Error constant
const (
	FailedStatus       = 0
	InternalError      = "Somethings Went Wrong!"
	DatabaseError      = "Database Error!"
	UserNotFound       = "User doesn't exist!"
	BadRequestError    = "Bad Request Error!"
	InvalidAccess      = "invalid_access"
	FailedToFetchToken = "Failed to fetch token!"
	NoDataFound        = "sql: no rows in result set"
	Unauthorized       = "Unauthorized"
)

// Success constant
const (
	SuccessStatus     = 1
	GetProfileSuccess = "Get user profile successfully."
)

// Common constant
const (
	NilString        = ""
	EmptyData        = "Data Not Provided!"
	SetAsNotSelected = "Set0"
	SetAsSelected    = "Set1"
)
