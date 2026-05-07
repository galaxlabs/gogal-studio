package lifecycle

import "fmt"

type Action string

const (
	ActionSave   Action = "Save"
	ActionSubmit Action = "Submit"
	ActionCancel Action = "Cancel"
	ActionAmend  Action = "Amend"
)

type ActionResult struct {
	Allowed bool   `json:"allowed"`
	Message string `json:"message"`
}

func CanPerform(action Action, docstatus int, isSubmittable bool) ActionResult {
	status := DocStatus(docstatus)

	if !IsValid(docstatus) {
		return ActionResult{
			Allowed: false,
			Message: fmt.Sprintf("Invalid docstatus: %d", docstatus),
		}
	}

	switch action {
	case ActionSave:
		if CanEdit(docstatus, isSubmittable) {
			return ActionResult{Allowed: true, Message: "Document can be saved"}
		}
		return ActionResult{Allowed: false, Message: "Document cannot be edited in current status"}

	case ActionSubmit:
		if !isSubmittable {
			return ActionResult{Allowed: false, Message: "DocType is not submittable"}
		}
		if status == Draft {
			return ActionResult{Allowed: true, Message: "Document can be submitted"}
		}
		return ActionResult{Allowed: false, Message: "Only draft documents can be submitted"}

	case ActionCancel:
		if !isSubmittable {
			return ActionResult{Allowed: false, Message: "DocType is not submittable"}
		}
		if status == Submitted {
			return ActionResult{Allowed: true, Message: "Document can be cancelled"}
		}
		return ActionResult{Allowed: false, Message: "Only submitted documents can be cancelled"}

	case ActionAmend:
		if !isSubmittable {
			return ActionResult{Allowed: false, Message: "DocType is not submittable"}
		}
		if status == Cancelled {
			return ActionResult{Allowed: true, Message: "Document can be amended"}
		}
		return ActionResult{Allowed: false, Message: "Only cancelled documents can be amended"}

	default:
		return ActionResult{
			Allowed: false,
			Message: fmt.Sprintf("Unknown lifecycle action: %s", action),
		}
	}
}

func AvailableActions(docstatus int, isSubmittable bool) []Action {
	actions := []Action{}

	for _, action := range []Action{
		ActionSave,
		ActionSubmit,
		ActionCancel,
		ActionAmend,
	} {
		result := CanPerform(action, docstatus, isSubmittable)
		if result.Allowed {
			actions = append(actions, action)
		}
	}

	return actions
}
