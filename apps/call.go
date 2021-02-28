package apps

import (
	"encoding/json"
	"io"

	"github.com/mattermost/mattermost-plugin-apps/server/utils/md"
)

type CallType string

const (
	// CallTypeSubmit (default) indicates the intent to take action.
	CallTypeSubmit = CallType("")
	// CallTypeForm retrieves the form definition for the current set of values,
	// and the context.
	CallTypeForm = CallType("form")
	// CallTypeCancel is used for for the (rare?) case of when the form with
	// SubmitOnCancel set is dismissed by the user.
	CallTypeCancel = CallType("cancel")
	// CallTypeLookup is used to fetch items for dynamic select elements
	CallTypeLookup = CallType("lookup")
)

// A Call invocation is supplied a BotAccessToken as part of the context. If a
// call needs acting user's or admin tokens, it should be specified in the
// Expand section.
//
// If a user or admin token are required and are not available from previous
// consent, the appropriate OAuth flow is launched, and the Call is executed
// upon its success.
//
// Call should be abbreviated as `call`.
type Call struct {
	URL        string                 `json:"url,omitempty"`
	Type       CallType               `json:"type,omitempty"`
	Values     map[string]interface{} `json:"values,omitempty"`
	Context    *Context               `json:"context,omitempty"`
	RawCommand string                 `json:"raw_command,omitempty"`
	Expand     *Expand                `json:"expand,omitempty"`
}

type CallResponseType string

const (
	// CallResponseTypeOK indicates that the call succeeded, and returns
	// Markdown and Data.
	// TODO update OK to be ["ok", "" default], update redux, webapp?
	CallResponseTypeOK = CallResponseType("")

	// CallResponseTypeOK indicates an error, returns Error.
	CallResponseTypeError = CallResponseType("error")

	// CallResponseTypeForm returns the definition of the form to display for
	// the inputs.
	CallResponseTypeForm = CallResponseType("form")

	// CallResponseTypeCall indicates that another Call that should be executed
	// (from the user-agent?). Call field is returned.
	CallResponseTypeCall = CallResponseType("call")

	// CallResponseTypeNavigate indicates that the user should be forcefully
	// navigated to a URL, which may be a channel in Mattermost. NavigateToURL
	// and UseExternalBrowser are expected to be returned.
	// TODO should CallResponseTypeNavigate be a variation of CallResponseTypeOK?
	CallResponseTypeNavigate = CallResponseType("navigate")
)

// CallResponse should be abbreviated as `cr`.
type CallResponse struct {
	Type CallResponseType `json:"type,omitempty"`

	// Used in CallResponseTypeOK to return the displayble, and JSON results
	Markdown md.MD       `json:"markdown,omitempty"`
	Data     interface{} `json:"data,omitempty"`

	// Used in CallResponseTypeError
	ErrorText string `json:"error,omitempty"`

	// Used in CallResponseTypeNavigate
	NavigateToURL      string `json:"navigate_to_url,omitempty"`
	UseExternalBrowser bool   `json:"use_external_browser,omitempty"`

	// Used in CallResponseTypeCall
	Call *Call `json:"call,omitempty"`

	// Used in CallResponseTypeForm
	Form *Form `json:"form,omitempty"`
}

func NewCallResponse(txt md.MD, data interface{}, err error) *CallResponse {
	if err != nil {
		return NewErrorCallResponse(err)
	}
	return &CallResponse{
		Type:     CallResponseTypeOK,
		Markdown: txt,
		Data:     data,
	}
}

func NewErrorCallResponse(err error) *CallResponse {
	return &CallResponse{
		Type: CallResponseTypeError,
		// TODO <><> ticket use MD and Data, remove Error
		ErrorText: err.Error(),
	}
}

// Error() makes CallResponse a valid error, for convenience
func (cr *CallResponse) Error() string {
	if cr.Type == CallResponseTypeError {
		return cr.ErrorText
	}
	return ""
}

func UnmarshalCallFromData(data []byte) (*Call, error) {
	call := Call{}
	err := json.Unmarshal(data, &call)
	if err != nil {
		return nil, err
	}
	return &call, nil
}

func UnmarshalCallFromReader(in io.Reader) (*Call, error) {
	call := Call{}
	err := json.NewDecoder(in).Decode(&call)
	if err != nil {
		return nil, err
	}
	return &call, nil
}

func MakeCall(url string, namevalues ...string) *Call {
	call := &Call{
		URL: url,
	}

	values := map[string]interface{}{}
	for len(namevalues) > 0 {
		switch len(namevalues) {
		case 1:
			values[namevalues[0]] = ""
			namevalues = namevalues[1:]

		default:
			values[namevalues[0]] = namevalues[1]
			namevalues = namevalues[2:]
		}
	}
	if len(values) > 0 {
		call.Values = values
	}
	return call
}

func (call *Call) GetStringValue(name, defaultValue string) string {
	if len(call.Values) == 0 {
		return defaultValue
	}
	v := call.Values[name]
	if v == nil {
		return defaultValue
	}
	switch v := v.(type) {
	case string:
		return v

	case map[string]interface{}:
		if len(v) == 0 {
			return defaultValue
		}
		if s, ok := v["value"].(string); ok {
			return s
		}
		return defaultValue

	default:
		return defaultValue
	}
}

func (call *Call) GetBoolValue(name string) bool {
	if len(call.Values) == 0 {
		return false
	}
	v := call.Values[name]
	if v == nil {
		return false
	}
	b, ok := v.(bool)
	if !ok {
		return false
	}
	return b
}
