package param

var (
	VALID   = "1"
	INVALID = "0"
)

type Response struct {
	Status bool        `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}
