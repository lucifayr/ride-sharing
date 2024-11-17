package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/common"
	"ride_sharing_api/app/rest"
	"ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestHandleCreateRide(t *testing.T) {
	utils.SetupTestDBs()

	db, err := utils.InitDb(path.Join(utils.ProjectRoot(), "db/testing/0001-with-three-users.db"))
	assert.Nil(err)
	queries := sqlc.New(db)
	handler := rest.NewRESTApi(queries)

	api := httptest.NewServer(handler)
	defer api.Close()

	// Missing authentication token
	resp, err := api.Client().Post(api.URL+"/rides", "application/json", bytes.NewReader([]byte{}))
	assert.Nil(err)
	assert.True(resp.StatusCode == 400, "Invalid status code", resp.StatusCode)

	// Invalid HTTP method
	resp, err = api.Client().Get(api.URL + "/rides")
	assert.Nil(err)
	assert.True(resp.StatusCode == 405, "Invalid status code", resp.StatusCode)

	// Invalid token
	req, err := http.NewRequest("POST", api.URL+"/rides", bytes.NewReader([]byte{}))
	assert.Nil(err)
	req.Header.Add("Authorization", "invalid-token")
	resp, err = api.Client().Do(req)
	assert.Nil(err)
	assert.True(resp.StatusCode == 401, "Invalid status code", resp.StatusCode)

	// Missing request body
	req, err = http.NewRequest("POST", api.URL+"/rides", bytes.NewReader([]byte{}))
	assert.Nil(err)
	req.Header.Add("Authorization", common.TEST_USERS["0001-with-three-users"][0].AccessToken.String)
	resp, err = api.Client().Do(req)
	assert.Nil(err)
	assert.True(resp.StatusCode == 400, "Invalid status code", resp.StatusCode)
	data, err := io.ReadAll(resp.Body)
	assert.Nil(err)
	assert.True(strings.Contains(string(data), "Invalid JSON"), "Invalid response body", string(data))

	// Missing fields in request body
	req, err = http.NewRequest("POST", api.URL+"/rides", bytes.NewReader([]byte(`{ "locationTo": "New York" }`)))
	assert.Nil(err)
	req.Header.Add("Authorization", common.TEST_USERS["0001-with-three-users"][0].AccessToken.String)
	resp, err = api.Client().Do(req)
	assert.Nil(err)
	assert.True(resp.StatusCode == 400, "Invalid status code", resp.StatusCode)
	data, err = io.ReadAll(resp.Body)
	assert.Nil(err)
	assert.True(strings.Contains(string(data), "Missing/Invalid fields"), "Invalid response body", string(data))

	// Valid request
	req, err = http.NewRequest("POST", api.URL+"/rides", bytes.NewReader([]byte(fmt.Sprintf(`{
		"locationTo": "New York",
		"locationFrom": "San Francisco",
		"tackingPlaceAt": "2024-12-31T12:35:00+02:00",
		"driver": "%s"
	}`, common.TEST_USERS["0001-with-three-users"][1].ID))))
	assert.Nil(err)
	req.Header.Add("Authorization", common.TEST_USERS["0001-with-three-users"][0].AccessToken.String)
	resp, err = api.Client().Do(req)
	assert.Nil(err)
	assert.True(resp.StatusCode == 201, "Invalid status code", resp.StatusCode)
	data, err = io.ReadAll(resp.Body)
	assert.Nil(err)
	var ride sqlc.Ride
	err = json.Unmarshal(data, &ride)
	assert.Nil(err)
	assert.True(ride.CreatedBy == common.TEST_USERS["0001-with-three-users"][0].ID)
	assert.True(ride.Driver == common.TEST_USERS["0001-with-three-users"][1].ID)
	assert.True(ride.LocationFrom == "San Francisco")
	assert.True(ride.TackingPlaceAt == "2024-12-31T10:35:00Z")
}
