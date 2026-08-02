package contract

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const SyntheticResourceView = "acceptance.synthetic.resource.view"

var resourceIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)
var ErrUnknownRoute = errors.New("route is not allowlisted")
var ErrMethodMismatch = errors.New("method does not match route contract")
var ErrMalformedPath = errors.New("malformed route path")

type Route struct{ ID, ContractVersion, Method, Template, Action, ResourceType, ResourceID string }

func Resolve(routeID, method, rawPath string) (Route, error) {
	if routeID != SyntheticResourceView {
		return Route{}, ErrUnknownRoute
	}
	if method != "GET" {
		return Route{}, ErrMethodMismatch
	}
	u, err := url.ParseRequestURI(rawPath)
	if err != nil || u.RawQuery != "" || u.Fragment != "" {
		return Route{}, ErrMalformedPath
	}
	const prefix = "/_acceptance/authorization/resources/"
	if !strings.HasPrefix(u.EscapedPath(), prefix) {
		return Route{}, ErrMalformedPath
	}
	escapedID := strings.TrimPrefix(u.EscapedPath(), prefix)
	if escapedID == "" || strings.Contains(escapedID, "/") {
		return Route{}, ErrMalformedPath
	}
	resourceID, err := url.PathUnescape(escapedID)
	if err != nil || !resourceIDPattern.MatchString(resourceID) {
		return Route{}, fmt.Errorf("%w: invalid resource id", ErrMalformedPath)
	}
	return Route{
		ID: SyntheticResourceView, ContractVersion: "1", Method: "GET",
		Template: prefix + "{resource_id}", Action: "resource:view",
		ResourceType: "resource", ResourceID: resourceID,
	}, nil
}
