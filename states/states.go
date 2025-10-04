package states

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/leetatech/leeta_golang_libraries/restclient"
)

const getStatePath = "/states/"

// GetState fetches the details of a specific state by name from the given service URL.
// Returns the state information or an error if the request or decoding fails.
func GetState(ctx context.Context, name, url string) (state State, err error) {
	getStateURL := url + getStatePath + name
	resp, err := restclient.DoHTTPRequest(ctx, http.MethodGet, nil, getStateURL)
	if err != nil {
		return state, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&state)
	if err != nil {
		return state, err
	}

	return
}

// GetAllStates retrieves a list of all states from the given service URL.
// Returns a slice of State or an error if the request or decoding fails.
func GetAllStates(ctx context.Context, url string) (stateList []State, err error) {
	listStateURL := url + getStatePath
	resp, err := restclient.DoHTTPRequest(ctx, http.MethodGet, nil, listStateURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&stateList)
	if err != nil {
		return nil, err
	}

	return stateList, nil
}
