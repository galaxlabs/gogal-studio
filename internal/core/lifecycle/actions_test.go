package lifecycle

import "testing"

func TestSubmittableDraftActions(t *testing.T) {
	actions := AvailableActions(0, true)

	if !containsAction(actions, ActionSave) {
		t.Fatal("draft should allow Save")
	}

	if !containsAction(actions, ActionSubmit) {
		t.Fatal("draft should allow Submit")
	}

	if containsAction(actions, ActionCancel) {
		t.Fatal("draft should not allow Cancel")
	}
}

func TestSubmittableSubmittedActions(t *testing.T) {
	actions := AvailableActions(1, true)

	if containsAction(actions, ActionSave) {
		t.Fatal("submitted document should not allow Save")
	}

	if !containsAction(actions, ActionCancel) {
		t.Fatal("submitted document should allow Cancel")
	}
}

func TestSubmittableCancelledActions(t *testing.T) {
	actions := AvailableActions(2, true)

	if !containsAction(actions, ActionAmend) {
		t.Fatal("cancelled document should allow Amend")
	}

	if containsAction(actions, ActionSave) {
		t.Fatal("cancelled document should not allow Save")
	}
}

func TestNonSubmittableActions(t *testing.T) {
	actions := AvailableActions(0, false)

	if !containsAction(actions, ActionSave) {
		t.Fatal("non-submittable draft should allow Save")
	}

	if containsAction(actions, ActionSubmit) {
		t.Fatal("non-submittable document should not allow Submit")
	}
}

func containsAction(actions []Action, target Action) bool {
	for _, action := range actions {
		if action == target {
			return true
		}
	}

	return false
}
