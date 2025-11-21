package code

// 通用: 基本错误
// Code must start with 1xxxxx
const (
	// ErrSuccess - 200: OK.
	ErrSuccess int = iota + 100001

	// ErrUnknown - 500: Internal server error.
	ErrUnknown

	// ErrBind - 400: Error occurred while binding the request body to the struct.
	ErrBind

	// ErrValidation - 400: Validation failed.
	ErrValidation

	// ErrTokenInvalid - 401: Token invalid.
	ErrTokenInvalid

	// ErrQuery - 400: Query parameter is invalid.
	ErrQuery

	// ErrDuplicatedRequest - 400: Duplicated request.
	ErrDuplicatedRequest
)

// 通用：数据库类错误
const (
	// ErrDatabase - 500: Database error.
	ErrDatabase int = iota + 100101
)

// 通用：认证授权类错误
const (
	// ErrEncrypt - 401: Error occurred while encrypting the user password.
	ErrEncrypt int = iota + 100201

	// ErrSignatureInvalid - 401: Signature is invalid.
	ErrSignatureInvalid

	// ErrExpired - 401: Token expired.
	ErrExpired

	// ErrInvalidAuthHeader - 401: Invalid authorization header.
	ErrInvalidAuthHeader

	// ErrMissingHeader - 401: The `Authorization` header was empty.
	ErrMissingHeader

	// ErrorExpired - 401: Token expired.
	ErrorExpired

	// ErrPasswordIncorrect - 401: Password was incorrect.
	ErrPasswordIncorrect

	// PermissionDenied - 403: Permission denied.
	ErrPermissionDenied

	// ErrUnauthorized - 401: Unauthorized access.
	ErrUnauthorized

	// ErrInvalidParams - 400: Invalid parameters.
	ErrInvalidParams

	// ErrNotFound - 404: Resource not found.
	ErrNotFound

	// ErrInternalServer - 500: Internal server error.
	ErrInternalServer
)

// 通用：编解码类错误
const (
	// ErrEncodingFailed - 500: Encoding failed due to an error with the data.
	ErrEncodingFailed int = iota + 100301

	// ErrDecodingFailed - 500: Decoding failed due to an error with the data.
	ErrDecodingFailed

	// ErrInvalidJSON - 500: Data is not valid JSON.
	ErrInvalidJSON

	// ErrEncodingJSON - 500: JSON data could not be encoded.
	ErrEncodingJSON

	// ErrDecodingJSON - 500: JSON data could not be decoded.
	ErrDecodingJSON

	// ErrInvalidYaml - 500: Data is not valid Yaml.
	ErrInvalidYaml

	// ErrEncodingYaml - 500: Yaml data could not be encoded.
	ErrEncodingYaml

	// ErrDecodingYaml - 500: Yaml data could not be decoded.
	ErrDecodingYaml
)

// 用户数据类错误
const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 100401
)

// 用户账号类错误
const (
	// ErrAccountNotFound - 404: Account not found.
	ErrAccountNotFound int = iota + 100501

	// ErrAccountInvalidPassword - 400: Invalid password.
	ErrAccountInvalidPassword

	// ErrAccountAlreadyBound - 400: Account already bound.
	ErrAccountAlreadyBound int = iota + 100601
)

func init() {
	Register(ErrSuccess, 200, "OK")
	Register(ErrUnknown, 500, "Internal server error")
	Register(ErrBind, 400, "Error occurred while binding the request body to the struct")
	Register(ErrValidation, 400, "Validation failed")
	Register(ErrTokenInvalid, 401, "Token invalid")
	Register(ErrDatabase, 500, "Database error")
	Register(ErrEncrypt, 401, "Error occurred while encrypting the user password")
	Register(ErrSignatureInvalid, 401, "Signature is invalid")
	Register(ErrExpired, 401, "Token expired")
	Register(ErrInvalidAuthHeader, 401, "Invalid authorization header")
	Register(ErrMissingHeader, 401, "The `Authorization` header was empty")
	Register(ErrPasswordIncorrect, 401, "Password was incorrect")
	Register(ErrPermissionDenied, 403, "Permission denied")
	Register(ErrEncodingFailed, 500, "Encoding failed due to an error with the data")
	Register(ErrDecodingFailed, 500, "Decoding failed due to an error with the data")
	Register(ErrInvalidJSON, 500, "Data is not valid JSON")
	Register(ErrEncodingJSON, 500, "JSON data could not be encoded")
	Register(ErrDecodingJSON, 500, "JSON data could not be decoded")
	Register(ErrInvalidYaml, 500, "Data is not valid Yaml")
	Register(ErrEncodingYaml, 500, "Yaml data could not be encoded")
	Register(ErrDecodingYaml, 500, "Yaml data could not be decoded")
	Register(ErrQuery, 400, "Query parameter is invalid")
	Register(ErrUnauthorized, 401, "Unauthorized")
	Register(ErrInvalidParams, 400, "Invalid parameters")
	Register(ErrNotFound, 404, "Resource not found")
	Register(ErrInternalServer, 500, "Internal server error")
	Register(ErrUserNotFound, 404, "User not found")
	Register(ErrAccountNotFound, 404, "Account not found")
	Register(ErrAccountInvalidPassword, 400, "Invalid password")
	Register(ErrAccountAlreadyBound, 400, "Account already bound")
	Register(ErrDuplicatedRequest, 400, "Duplicated request")
}
