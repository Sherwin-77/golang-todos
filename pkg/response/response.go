package response

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	LastPage int `json:"last_page"`
	Total    int `json:"total"`
}

func NewResponse(code int, message string, data interface{}, meta *Meta) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}
