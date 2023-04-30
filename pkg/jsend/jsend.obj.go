package jsend

import (
	"strings"

	"github.com/freekup/mini-wallet/pkg/cerror"
)

type (
	JSendStatusType string
	JSendResponse   struct {
		Status  JSendStatusType
		Code    int64
		Data    map[string]interface{}
		Message string
	}
)

const (
	JSendStatusTypeSuccess JSendStatusType = "success"
	JSendStatusTypeFail    JSendStatusType = "fail"
	JSendStatusTypeErr     JSendStatusType = "error"
)

func GenerateResponseSuccess(data map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"status": JSendStatusTypeSuccess,
		"data":   data,
	}

	return result
}

func GenerateResponseError(cerr cerror.CError) map[string]interface{} {
	var result = map[string]interface{}{}

	if cerr.IsSystemError() {
		result = map[string]interface{}{
			"status":  JSendStatusTypeErr,
			"message": cerr.Error(),
		}
	} else if cerr.IsValidationError() {
		data := make(map[string]interface{})

		for _, errs := range strings.Split(cerr.Error(), ";") {
			err := strings.Split(errs, "=")
			if len(err) >= 2 {
				data[err[0]] = err[1]
			}
		}

		result = map[string]interface{}{
			"status": JSendStatusTypeFail,
			"data":   data,
		}
	}

	return result
}
