package utils

type successResponse struct {
	RpcVersion string      `json:"jsonrpc"`
	Data       interface{} `json:"result"`
}

type errorResponse struct {
	RpcVersion string       `json:"jsonrpc"`
	Error      *errorObject `json:"error"`
}

type errorObject struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	JRPC_Version = "2.0"
)

func JRPCSuccess(data interface{}) *successResponse {
	response := &successResponse{
		JRPC_Version,
		data,
	}

	return response
}

func JRPCError(err error, data interface{}) *errorResponse {

	message := "Unknown"
	if err != nil {
		message = err.Error()
	}

	return JRPCErrorF(0, message, data)
}

func JRPCErrorF(code int, message string, data interface{}) *errorResponse {

	response := &errorResponse{
		"2.0",
		&errorObject{
			code,
			message,
			data,
		},
	}

	return response
}
