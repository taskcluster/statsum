//go:generate msgp.v0
//go:generate ffjson $GOFILE

package payload

// The Response hold the result code and message from any request.
type Response struct {
	Code    string `json:"code" msg:"code"`
	Message string `json:"message" msg:"message"`
}
