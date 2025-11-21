package code

// ErrCode implements `github.com/marmotedu/errors`.Coder interface.
type ErrCode struct {
	// C refers to the code of the ErrCode.
	C int

	// HTTP status that should be used for the associated error code.
	HTTP int

	// External (user) facing error text.
	Ext string

	// Ref specify the reference document.
	Ref string
}

// Code returns the integer code of ErrCode.
func (coder ErrCode) Code() int {
	return coder.C
}

// String implements stringer. String returns the external error message,
// if any.
func (coder ErrCode) String() string {
	return coder.Ext
}

// Reference returns the reference document.
func (coder ErrCode) Reference() string {
	return coder.Ref
}

// HTTPStatus returns the associated HTTP status code, if any. Otherwise,
// returns 200.
func (coder ErrCode) HTTPStatus() int {
	if coder.HTTP == 0 {
		return 500
	}
	return coder.HTTP
}

var ErrCodes = map[int]ErrCode{}

func Use(code int) ErrCode {
	if coder, ok := ErrCodes[code]; ok {
		return coder
	}
	return ErrCodes[ErrUnknown]
}

func Register(code int, httpStatus int, message string, refs ...string) {
	var reference string
	if len(refs) > 0 {
		reference = refs[0]
	}

	ErrCodes[code] = ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
		Ref:  reference,
	}

}
