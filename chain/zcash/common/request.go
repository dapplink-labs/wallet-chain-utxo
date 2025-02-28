package common

const METHOD = "method"
const PARAMS = "params"

var RequestData = map[string]interface{}{
	"jsonrpc": "1.0",
	"id":      "curltest",
	"method":  "",
	"params":  []any{},
}
