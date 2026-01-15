package errs

type ErrorCode int

func (e ErrorCode) Error() string {
	return errorMessages[e]
}

const (
	DatabaseError                ErrorCode = 1001
	DatabaseNoRecordError        ErrorCode = 1002
	UnmarshalError               ErrorCode = 1003
	MarshalError                 ErrorCode = 1004
	PasswordValidationError      ErrorCode = 1005
	EncryptionError              ErrorCode = 1006
	DecryptionError              ErrorCode = 1007
	DuplicateUserError           ErrorCode = 1008
	UserNotFoundError            ErrorCode = 1009
	IdentityNotFoundError        ErrorCode = 1010
	UserLockedError              ErrorCode = 1011
	CredentialsValidationError   ErrorCode = 1012
	TokenGenerationError         ErrorCode = 1013
	TokenValidationError         ErrorCode = 1014
	UserCategoryError            ErrorCode = 1015
	EmailSendingError            ErrorCode = 1016
	BusinessCategoryError        ErrorCode = 1017
	StatusesError                ErrorCode = 1018
	ErrorUnauthorized            ErrorCode = 1019
	EmailFormatError             ErrorCode = 1020
	ValidEmailHostError          ErrorCode = 1021
	ValidLeetaDomainError        ErrorCode = 1022
	FormParseError               ErrorCode = 1023
	OrderStatusesError           ErrorCode = 1024
	ProductCategoryError         ErrorCode = 1025
	ProductSubCategoryError      ErrorCode = 1026
	ProductStatusError           ErrorCode = 1027
	ForgotPasswordError          ErrorCode = 1028
	MissingUserNames             ErrorCode = 1029
	InvalidUserRoleError         ErrorCode = 1030
	InvalidIdentityError         ErrorCode = 1031
	InvalidOTPError              ErrorCode = 1032
	CartStatusesError            ErrorCode = 1033
	AmountPaidError              ErrorCode = 1034
	FeesStatusesError            ErrorCode = 1035
	InvalidPageRequestError      ErrorCode = 1036
	CartItemQuantityError        ErrorCode = 1037
	CartItemRequestQuantityError ErrorCode = 1038
	InvalidRequestError          ErrorCode = 1039
	InternalError                ErrorCode = 1040
	InvalidProductIdError        ErrorCode = 1041
	InvalidDeliveryFeeError      ErrorCode = 1042
	InvalidServiceFeeError       ErrorCode = 1043
	RestrictedAccessError        ErrorCode = 1044
	FeesError                    ErrorCode = 1045
	TemplateCreationError        ErrorCode = 1046
	AwsSessionError              ErrorCode = 1047
	SesSendEmailError            ErrorCode = 1048
	SnsSendSMSError              ErrorCode = 1049
	LGANotFoundError             ErrorCode = 1050
	PushNotificationError        ErrorCode = 1051
	DuplicateVendorBusinessError ErrorCode = 1052
	InvalidVendorIdError         ErrorCode = 1053
	TooManyVendorsError          ErrorCode = 1054
	S3ObjectNotFoundError        ErrorCode = 1055
	DuplicateRecordError         ErrorCode = 1056
	ErrorForbidden               ErrorCode = 1057
	ErrorImcompleteOrder         ErrorCode = 1058
)

var (
	errorTypes = map[ErrorCode]string{
		DatabaseError:                "DatabaseError",
		DatabaseNoRecordError:        "DatabaseNoRecordError",
		UnmarshalError:               "UnmarshalError",
		MarshalError:                 "MarshalError",
		PasswordValidationError:      "PasswordValidationError",
		EncryptionError:              "EncryptionError",
		DecryptionError:              "DecryptionError",
		DuplicateUserError:           "DuplicateUserError",
		UserNotFoundError:            "UserNotFoundError",
		IdentityNotFoundError:        "IdentityNotFoundError",
		UserLockedError:              "UserLockedError",
		CredentialsValidationError:   "CredentialsValidationError",
		TokenGenerationError:         "TokenGenerationError",
		TokenValidationError:         "TokenValidationError",
		UserCategoryError:            "UserCategoryError",
		EmailSendingError:            "EmailSendingError",
		BusinessCategoryError:        "BusinessCategoryError",
		StatusesError:                "StatusesError",
		ErrorUnauthorized:            "ErrorUnauthorized",
		EmailFormatError:             "EmailFormatError",
		ValidEmailHostError:          "ValidEmailHostError",
		ValidLeetaDomainError:        "ValidLeetaDomainError",
		FormParseError:               "FormParseError",
		OrderStatusesError:           "OrderStatusesError",
		ProductCategoryError:         "ProductCategoryError",
		ProductSubCategoryError:      "ProductSubCategoryError",
		ProductStatusError:           "ProductStatusError",
		ForgotPasswordError:          "ForgotPasswordError",
		MissingUserNames:             "MissingUserNamesError",
		InvalidUserRoleError:         "InvalidUserRoleError",
		InvalidIdentityError:         "InvalidIdentityError",
		InvalidOTPError:              "InvalidOTPError",
		CartStatusesError:            "CartStatusesError",
		AmountPaidError:              "AmountPaidError",
		FeesStatusesError:            "FeesStatusesError",
		InvalidPageRequestError:      "InvalidPageRequestError",
		CartItemQuantityError:        "CartItemQuantityError",
		CartItemRequestQuantityError: "CartItemRequestQuantityError",
		InvalidRequestError:          "InvalidRequestError",
		InternalError:                "InternalError",
		InvalidProductIdError:        "InvalidProductIdError",
		InvalidDeliveryFeeError:      "InvalidDeliveryFeeError",
		InvalidServiceFeeError:       "InvalidServiceFeeError",
		RestrictedAccessError:        "RestrictedAccessError",
		FeesError:                    "FeesError",
		TemplateCreationError:        "TemplateCreationError",
		AwsSessionError:              "AwsSessionError",
		SesSendEmailError:            "SesSendEmailError",
		SnsSendSMSError:              "SnsSendSMSError",
		LGANotFoundError:             "LGANotFoundError",
		PushNotificationError:        "PushNotificationError",
		DuplicateVendorBusinessError: "DuplicateVendorBusinessError",
		InvalidVendorIdError:         "InvalidVendorIdError",
		TooManyVendorsError:          "TooManyVendorsError",
		S3ObjectNotFoundError:        "S3ObjectNotFoundError",
		DuplicateRecordError:         "DuplicateRecordError",
		ErrorForbidden:               "ErrorForbidden",
		ErrorImcompleteOrder:         "ErrorImcompleteOrder",
	}

	errorMessages = map[ErrorCode]string{
		DatabaseError:                "An error occurred while reading from the database",
		DatabaseNoRecordError:        "An error occurred because no record was found",
		UnmarshalError:               "An error occurred while unmarshalling data",
		MarshalError:                 "An error occurred while marshaling data",
		PasswordValidationError:      "An error occurred while validating password. | Password must contain at least six character long, one uppercase letter, one lowercase letter, one digit, and one special character | password and confirm password don't match",
		EncryptionError:              "An error occurred while encrypting",
		DecryptionError:              "An error occurred while decrypting",
		DuplicateUserError:           "An error occurred because user already exists",
		UserNotFoundError:            "An error occurred because this is not a registered user",
		IdentityNotFoundError:        "An error occurred because this user identity is not known",
		UserLockedError:              "An error occurred because this user is locked",
		CredentialsValidationError:   "An error occurred because the credentials are invalid",
		TokenGenerationError:         "An error occurred while generating token",
		TokenValidationError:         "An error occurred because the token is invalid | validated | expired",
		UserCategoryError:            "An error occurred because the user category is invalid",
		EmailSendingError:            "An error occurred while sending email",
		BusinessCategoryError:        "An error occurred because the business category is invalid",
		StatusesError:                "An error occurred because the statuses are invalid",
		ErrorUnauthorized:            "An error occurred because the user is unauthorized",
		EmailFormatError:             "An error occurred because the email format is invalid",
		ValidEmailHostError:          "An error occurred because the domain does not exist or cannot receive emails",
		ValidLeetaDomainError:        "An error occurred because the domain does not belong to leeta or cannot receive emails",
		FormParseError:               "An error occurred because the form parse failed or file retrieval failed",
		OrderStatusesError:           "An error occurred because the order status is invalid",
		ProductCategoryError:         "An error occurred because the product category is invalid",
		ProductSubCategoryError:      "An error occurred because the product subcategory is invalid",
		ProductStatusError:           "An error occurred because the product status is invalid",
		ForgotPasswordError:          "An error occurred while trying to reset a user password",
		MissingUserNames:             "An error occurred because user first name/last name was not found",
		InvalidUserRoleError:         "An error occurred because the user is trying to login with the wrong app",
		InvalidIdentityError:         "An error occurred because the user identity data is invalid",
		InvalidOTPError:              "An error occurred because the OTP is invalid",
		CartStatusesError:            "An error occurred because the cart status is invalid",
		AmountPaidError:              "An error occurred because the amount paid is invalid",
		FeesStatusesError:            "An error occurred because the fees status is invalid",
		InvalidPageRequestError:      "An error occurred because the page request field is required",
		CartItemQuantityError:        "An error occurred because the stored cart item quantity/weight is already 0. Please delete the item or increase the quantity to continue",
		CartItemRequestQuantityError: "An error occurred because the request quantity/weight field is 0. Please increase the quantity/weight to continue",
		InvalidRequestError:          "An error occurred because the request is invalid",
		InternalError:                "An error has occurred in the server",
		InvalidProductIdError:        "An error occurred because the product id is invalid",
		InvalidDeliveryFeeError:      "An error occurred because the delivery fee is invalid",
		InvalidServiceFeeError:       "An error occurred because the service fee is invalid",
		RestrictedAccessError:        "User do not have authorization to access this endpoint",
		FeesError:                    "There is an error with the application fees",
		TemplateCreationError:        "An error occurred while creating template",
		SnsSendSMSError:              "An error occurred while sending SMS",
		AwsSessionError:              "An error occurred while creating aws session",
		SesSendEmailError:            "An error occurred while sending email",
		LGANotFoundError:             "Leeta is not available in your region",
		PushNotificationError:        "An error occurred while sending push notification",
		DuplicateVendorBusinessError: "An error occurred because this vendor's business has already been registered",
		InvalidVendorIdError:         "An error occurred because the vendor id is invalid",
		TooManyVendorsError:          "An error occurred because the vendor already has another vendor item in cart",
		S3ObjectNotFoundError:        "Object not found in s3 bucket",
		DuplicateRecordError:         "The record with this unique id already exists in the db",
		ErrorForbidden:               "User does not have sufficient access",
		ErrorImcompleteOrder:         "The Order is not complete",
	}
)
