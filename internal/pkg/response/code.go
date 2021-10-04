package response

const (
	SUCCESS = 200
	ERROR   = 500
)

var codeMsg = map[int]string{
	SUCCESS:                "OK",
	ERROR:                  "FAIL",
}

func GetErrMsg(code int) string {
	msg := "提示信息"
	if _, ok := codeMsg[code]; ok {
		msg = codeMsg[code]
	}
	return msg
}



