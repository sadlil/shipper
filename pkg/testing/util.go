package testing

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/pmezard/go-difflib/difflib"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	kubetesting "k8s.io/client-go/testing"
	"sigs.k8s.io/yaml"
)

const (
	NoResyncPeriod time.Duration = 0

	ContextLines = 4

	TestNamespace = "test-namespace"
	TestApp       = "shipper-test"
	TestRegion    = "eu-west"
	TestCluster   = "test-cluster"

	E2ETestNamespaceLabel = "shipper-e2e-test"
)

// CheckActions takes a slice of expected actions and a slice of observed
// actions (typically obtained from fakeClient.Actions()) and compares them.
// Calls Errorf on t for every difference it finds.
func CheckActions(expected, actual []kubetesting.Action, t *testing.T) {
	for i, action := range actual {
		if len(expected) < i+1 {
			t.Errorf("%d unexpected actions:", len(actual)-len(expected))
			for _, unexpectedAction := range actual[i:] {
				t.Logf("\n%s", prettyPrintAction(unexpectedAction))
			}
			break
		}

		CheckAction(expected[i], action, t)
	}

	if len(expected) > len(actual) {
		t.Errorf("missing %d expected actions:", len(expected)-len(actual))
		for _, missingExpectedAction := range expected[len(actual):] {
			t.Logf("\n%s", prettyPrintAction(missingExpectedAction))
		}
	}
}

// ShallowCheckActions takes a slice of expected actions and a slice of observed
// actions (typically obtained from fakeClient.Actions()) and compares them
// shallowly. Calls Errorf on t for every difference it finds.
func ShallowCheckActions(expected, actual []kubetesting.Action, t *testing.T) {
	for i, action := range actual {
		if len(expected) < i+1 {
			t.Errorf("%d unexpected actions: %+v", len(actual)-len(expected), actual[i:])
			for _, unexpectedAction := range actual[i:] {
				t.Logf("\n%s", prettyPrintAction(unexpectedAction))
			}
			break
		}

		ShallowCheckAction(expected[i], action, t)
	}

	if len(expected) > len(actual) {
		t.Errorf("missing %d expected actions: %+v", len(expected)-len(actual), expected[len(actual):])
	}
}

// ShallowCheckAction checks the verb, resource, and namespace without looking
// at the objects involved. This is a stand-in until we port the Installation
// controller to not use 'nil' as the object involved in the kubetesting.Actions
// it expects.
func ShallowCheckAction(expected, actual kubetesting.Action, t *testing.T) {
	if !(expected.Matches(actual.GetVerb(), actual.GetResource().Resource) &&
		actual.GetSubresource() == expected.GetSubresource() &&
		actual.GetResource() == expected.GetResource()) {

		t.Errorf("expected\n\t%#v\ngot\n\t%#v", expected, actual)
		return
	}

	if expected.GetNamespace() != actual.GetNamespace() {
		t.Errorf("expected action in ns %q, got ns %q", expected.GetNamespace(), actual.GetNamespace())
		return
	}

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("expected action %T but got %T", expected, actual)
		return
	}
}

// CheckAction compares two individual actions and calls Errorf on t if it finds
// a difference.
func CheckAction(expected, actual kubetesting.Action, t *testing.T) {
	if !equality.Semantic.DeepEqual(expected, actual) {
		prettyExpected := prettyPrintAction(expected)
		prettyActual := prettyPrintAction(actual)

		diff, err := YamlDiff(prettyActual, prettyExpected)
		if err != nil {
			panic(fmt.Sprintf("couldn't generate yaml diff: %s", err))
		}

		t.Errorf("expected action is different from actual:\n%s", diff)
	}
}

// PrettyPrintActions pretty-prints a slice of actions, useful for
// creating a human-readable list for debugging.
func PrettyPrintActions(actions []kubetesting.Action, t *testing.T) {
	for _, action := range actions {
		t.Logf("\n%s", prettyPrintAction(action))
	}
}

// FilterActions, given a slice of observed actions, returns only those that
// change state. Useful for reducing the number of actions needed to check in
// tests.
func FilterActions(actions []kubetesting.Action) []kubetesting.Action {
	ignore := func(action kubetesting.Action) bool {
		for _, v := range []string{"list", "watch"} {
			for _, r := range []string{
				"applications",
				"capacitytargets",
				"clusters",
				"configmaps",
				"deployments",
				"endpoints",
				"installationtargets",
				"pods",
				"releases",
				"rolloutblocks",
				"secrets",
				"services",
				"traffictargets",
			} {
				if action.Matches(v, r) {
					return true
				}
			}
		}

		return false
	}

	var ret []kubetesting.Action
	for _, action := range actions {
		if ignore(action) {
			continue
		}

		ret = append(ret, action)
	}

	return ret
}

func CheckEvents(expectedOrderedEvents []string, receivedEvents []string, t *testing.T) {
	eq, diff := DeepEqualDiff(expectedOrderedEvents, receivedEvents)
	if !eq {
		t.Errorf("Events don't match expectation:\n%s", diff)
	}
}

func YamlDiff(a interface{}, b interface{}) (string, error) {
	yamlActual, err := yaml.Marshal(a)
	if err != nil {
		return "", err
	}

	yamlExpected, err := yaml.Marshal(b)
	if err != nil {
		return "", err
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(yamlExpected)),
		B:        difflib.SplitLines(string(yamlActual)),
		FromFile: "Expected",
		ToFile:   "Actual",
		Context:  ContextLines,
	}

	return difflib.GetUnifiedDiffString(diff)
}

func DeepEqualDiff(expected, actual interface{}) (bool, string) {
	if !reflect.DeepEqual(actual, expected) {
		diff, err := YamlDiff(actual, expected)
		if err != nil {
			panic(fmt.Sprintf("couldn't generate yaml diff: %s", err))
		}

		return false, diff
	}

	return true, ""
}

func prettyPrintAction(a kubetesting.Action) string {
	verb := a.GetVerb()
	gvk := a.GetResource()
	ns := a.GetNamespace()
	extra := ""

	template := fmt.Sprintf("Verb: %s\nGVK: %s\nNamespace: %s\n%%s--------\n%%s", verb, gvk.String(), ns)

	switch action := a.(type) {

	case kubetesting.CreateActionImpl:
		extra := fmt.Sprintf("Name: %s\n", action.Name)
		obj, err := yaml.Marshal(action.GetObject())
		if err != nil {
			panic(fmt.Sprintf("could not marshal %+v: %q", action.GetObject(), err))
		}

		return fmt.Sprintf(template, extra, string(obj))

	case kubetesting.UpdateActionImpl:
		obj, err := yaml.Marshal(action.GetObject())
		if err != nil {
			panic(fmt.Sprintf("could not marshal %+v: %q", action.GetObject(), err))
		}

		return fmt.Sprintf(template, extra, string(obj))

	case kubetesting.PatchActionImpl:
		extra = fmt.Sprintf("Name: %s\n", action.Name)
		patch := prettyPrintActionPatch(action)
		return fmt.Sprintf(template, extra, patch)

	case kubetesting.GetActionImpl:
		message := fmt.Sprintf("(no object body: GET %s)", action.GetName())
		return fmt.Sprintf(template, extra, message)

	case kubetesting.ListActionImpl:
		message := fmt.Sprintf("(no object body: GET %s)", action.GetKind())
		return fmt.Sprintf(template, extra, message)

	case kubetesting.WatchActionImpl:
		return fmt.Sprintf(template, "(no object body: WATCH)")

	case kubetesting.DeleteActionImpl:
		message := fmt.Sprintf("(no object body: DELETE %s)", action.GetName())
		return fmt.Sprintf(template, extra, message)

	case kubetesting.ActionImpl:
		message := fmt.Sprintf("(no object body: %s %s)", action.GetVerb(), action.GetResource())
		return fmt.Sprintf(template, extra, message)
	}

	panic(fmt.Sprintf("unknown action! patch printAction to support %T %+v", a, a))
}

func prettyPrintActionPatch(action kubetesting.PatchActionImpl) string {
	switch action.GetPatchType() {
	case types.MergePatchType:
		var obj map[string]interface{}

		err := json.Unmarshal(action.GetPatch(), &obj)
		if err != nil {
			panic(fmt.Sprintf("could not unmarshal %v: %s", action.GetPatch(), err))
		}

		str, err := yaml.Marshal(obj)
		if err != nil {
			panic(fmt.Sprintf("could not marshal %v: %s", obj, err))
		}

		return string(str)
	case types.JSONPatchType:
		panic("not implemented")
	case types.StrategicMergePatchType:
		panic("not implemented")
	}

	return ""
}

func NewDiscoveryAction(_ string) kubetesting.ActionImpl {
	// FakeDiscovery has a very odd way of generating fake actions for
	// discovery. We try to paper over that as best we can. The ignored
	// parameter is trying to be future proof in case FakeDiscovery ever
	// decides to fix its nasty ways and actually report the resource we're
	// trying to discover.
	return kubetesting.ActionImpl{
		Verb:     "get",
		Resource: schema.GroupVersionResource{Resource: "resource"},
	}
}
