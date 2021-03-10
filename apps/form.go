package apps

// Form defines what inputs a Call accepts, and how they can be gathered from
// the user, in Modal and Autocomplete modes.
//
// For a Modal, the form defines the modal entirely, and displays it when
// returned in response to a submit Call. Modals are dynamic in the sense that
// they may be refreshed entirely when values of certain fields change, and may
// contain dynamic select fields.
//
// For Autocomplete, a form can be bound to a sub-command. The behavior of
// autocomplete once the subcommand is selected is designed to mirror the
// functionality of the Modal. Some gaps and differences still remain.
//
// Requests for forms are calls, have expanded contexts, making it easy to
// generate forms specific to the user, channel, etc.
type Form struct {
	// Title, Header, and Footer are used for Modals only.
	Title  string `json:"title,omitempty"`
	Header string `json:"header,omitempty"`
	Footer string `json:"footer,omitempty"`

	// A fully-qualified URL, or a path to the form icon.
	// TODO do we default to the App icon?
	Icon string `json:"icon,omitempty"`

	// Call is the same definition used to submit, refresh the form, and to
	// lookup dynamic select options.
	Call *Call `json:"call,omitempty"`

	// SubmitButtons refers to a field name that must be a FieldTypeStaticSelect
	// or FieldTypeDynamicSelect.
	//
	// In Modal view, the field will be rendered as a list of buttons at the
	// bottom. Clicking one of them submits the Call, providing the button
	// reference as the corresponding Field's value. Leaving this property
	// blank, displays the default "OK".
	//
	// In Autocomplete, it is ignored.
	SubmitButtons string `json:"submit_buttons,omitempty"`

	// Adds a default "Cancel" button in the modal view
	CancelButton  bool `json:"cancel_button,omitempty"`
	SubmitOnCanel bool `json:"submit_on_cancel,omitempty"`

	// DependsOn is the list of field names that when changed force reloading
	// the form. Values of the fields with values that are not included in the
	// refreshed form are lost.
	DependsOn []string `json:"depends_on,omitempty"`

	// Fields is the list of fields in the form.
	Fields []*Field `json:"fields,omitempty"`
}
