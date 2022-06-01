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
	FailedToLogin      = "Login Failed!"
	FailedToSignup     = "Signup Failed!"
)

// Success constant
const (
	SuccessStatus   = 1
	LoginSuccess    = "Login successfully."
	SignupSuccess   = "Signup successfully."
	GetTokenSuccess = "Access token fetch successfully."
	SendLinkSuccess = "We've just sent you an email to reset your password."
)

// Common constant
const (
	NilString  = ""
	EmptyData  = "Data Not Provided!"
	ProfileURL = "https://res.cloudinary.com/dobanpo5b/image/upload/v1652076493/user_zrlnnh.jpg"
)
