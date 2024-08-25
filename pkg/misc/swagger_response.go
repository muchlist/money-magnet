package misc

import "github.com/muchlist/moneymagnet/pkg/paging"

// ResponseErr just for wrapping swaggo type generator
type ResponseErr struct {
	Error      string            `json:"error" example:"example error message"`
	ErrorField map[string]string `json:"error_field" example:"example_field:example_field is a required field"`
}

// Response500Err just for wrapping swaggo type generator
type Response500Err struct {
	Error string `json:"error" example:"name func: sub func: cause of error"`
}

// ResponseSuccess just for wrapping swaggo type generator
type ResponseSuccess struct {
	Data any `json:"data"`
}

// ResponseMessage just for wrapping swaggo type generator
type ResponseMessage struct {
	Data string `json:"data" example:"do thing success"`
}

// ResponseSuccessList just for wrapping swaggo type generator
type ResponseSuccessList struct {
	Data     []any           `json:"data"`
	Metadata paging.Metadata `json:"meta_data"`
}

// ResponseSuccessListCursor just for wrapping swaggo type generator
type ResponseSuccessListCursor struct {
	Data     []any                 `json:"data"`
	Metadata paging.CursorMetadata `json:"meta_data"`
}
