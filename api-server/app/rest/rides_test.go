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
	"ride_sharing_api/app/rest"
	"ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestHandleCreateRide(t *testing.T) {
	db := utils.InitTestDB(path.Join(utils.ProjectRoot(), "db/testing/setup/0001-handle-create-rides.sql"))
	queries := sqlc.New(db)
	handler := rest.NewRESTApi(queries)

	api := httptest.NewServer(handler)
	defer api.Close()

	testAuth(api, "/rides", "POST")

	// Invalid HTTP method
	resp, err := api.Client().Get(api.URL + "/rides")
	assert.Nil(err)
	assert.Eq(resp.StatusCode, 405)

	// Missing request body
	req, err := http.NewRequest("POST", api.URL+"/rides", bytes.NewReader([]byte{}))
	assert.Nil(err)
	req.Header.Add("Authorization", "yr0osJ1-kQrQsOXzMGNVhbGUzdA0seGWAMK70WERFRWU5NzKcrQ1R2U_8ofXubwLbWxJYQK9hvj9xonabMMroA6oPjfnFuFR_zwOugdNGZVOwo6l8zczvFYRnGUdncOWv5Ckdy5eyB0leWpH7sDI_hbAxyiKljceGnKX-hcvB9MwjnsAiJMZ6EC_nAV-6ujEwM-YbPbYwndTEyY7CgDBp9gYrcOlvs9z_yf5sM_WQlziZFVVyGVoJyWDl-a1XbyLiagscmTeDs0pxQO0BH0oBF5qW8IRDIWAOuaSz3K9eygpqKQIxTFVq_psqaZT_qrhHI-3k-OPBbtWq9pF32-wVxNoFJMB3YvY17DgQfxxvzgckUH5YFlNks1cUgroHk2CIjtgs-9eskUzOrCzBKW3-EBcuyNrttnIePAkdVl2NC586fkBCVnKqfVIKYwm-ZrdCHxQVTZwGcswGnUP-YajlwZhmM-jgBjXIAJfWihcQTrDGmWz-0z8R8kycMdASguZXnQolGTvUOsOT21kFC4fwF-XQRi0tPh4mg0Bj1QN9y5sgibripVhCXQ7ma9QbbYL9ooAax6wqEU7b5-Gfai_r1ZLI5WcjOkI0ePAa2PikIC1b5nAMaz0c9y7Sv-hVAYtVzW5VB6PRJ4f5DoI_6KlGx6jE1AzmIEyPp_4_ImIUhHBlGUa7kikZkqUTtr9vSaz84EvQzT81wt3ULBLvA89Cr5rOWgAlNfmul3JZtJwfUuW39Mxc6QQN1mLUyKIUiofZImwkLqlACuriArAhMM_E8qo2V9sHSRVhZA_NOnKOYujsoFTTdr4vb2CWyeVIAEWT2YCueSMXinGL1Gmbxcczy9Hi2LoupnGYlQr9KgP5V_UrRvl_isC1MgUArQ25nIkdBNpUREW7a31bqWibAamOCgLP8bS20DERUD3-bKcDYDSDq9cEP2pKBRm_WyVQqCNYPIUpDOmDd9SEAZ3J_WveApSIDJlDt0j_nTibImctu6he92Kp63L5_A8nG6wBWW363CZ7tgktoY3KidPwbByX35BQRTUyE7wYxAqzdcF8Jd_n24SLHxC")
	resp, err = api.Client().Do(req)
	assert.Nil(err)
	assert.Eq(resp.StatusCode, 400)
	data, err := io.ReadAll(resp.Body)
	assert.Nil(err)
	assert.True(strings.Contains(string(data), "Invalid JSON"), "Invalid response body", string(data))

	// Missing fields in request body
	req, err = http.NewRequest("POST", api.URL+"/rides", bytes.NewReader([]byte(`{ "locationTo": "New York" }`)))
	assert.Nil(err)
	req.Header.Add("Authorization", "yr0osJ1-kQrQsOXzMGNVhbGUzdA0seGWAMK70WERFRWU5NzKcrQ1R2U_8ofXubwLbWxJYQK9hvj9xonabMMroA6oPjfnFuFR_zwOugdNGZVOwo6l8zczvFYRnGUdncOWv5Ckdy5eyB0leWpH7sDI_hbAxyiKljceGnKX-hcvB9MwjnsAiJMZ6EC_nAV-6ujEwM-YbPbYwndTEyY7CgDBp9gYrcOlvs9z_yf5sM_WQlziZFVVyGVoJyWDl-a1XbyLiagscmTeDs0pxQO0BH0oBF5qW8IRDIWAOuaSz3K9eygpqKQIxTFVq_psqaZT_qrhHI-3k-OPBbtWq9pF32-wVxNoFJMB3YvY17DgQfxxvzgckUH5YFlNks1cUgroHk2CIjtgs-9eskUzOrCzBKW3-EBcuyNrttnIePAkdVl2NC586fkBCVnKqfVIKYwm-ZrdCHxQVTZwGcswGnUP-YajlwZhmM-jgBjXIAJfWihcQTrDGmWz-0z8R8kycMdASguZXnQolGTvUOsOT21kFC4fwF-XQRi0tPh4mg0Bj1QN9y5sgibripVhCXQ7ma9QbbYL9ooAax6wqEU7b5-Gfai_r1ZLI5WcjOkI0ePAa2PikIC1b5nAMaz0c9y7Sv-hVAYtVzW5VB6PRJ4f5DoI_6KlGx6jE1AzmIEyPp_4_ImIUhHBlGUa7kikZkqUTtr9vSaz84EvQzT81wt3ULBLvA89Cr5rOWgAlNfmul3JZtJwfUuW39Mxc6QQN1mLUyKIUiofZImwkLqlACuriArAhMM_E8qo2V9sHSRVhZA_NOnKOYujsoFTTdr4vb2CWyeVIAEWT2YCueSMXinGL1Gmbxcczy9Hi2LoupnGYlQr9KgP5V_UrRvl_isC1MgUArQ25nIkdBNpUREW7a31bqWibAamOCgLP8bS20DERUD3-bKcDYDSDq9cEP2pKBRm_WyVQqCNYPIUpDOmDd9SEAZ3J_WveApSIDJlDt0j_nTibImctu6he92Kp63L5_A8nG6wBWW363CZ7tgktoY3KidPwbByX35BQRTUyE7wYxAqzdcF8Jd_n24SLHxC")
	resp, err = api.Client().Do(req)
	assert.Nil(err)
	assert.Eq(resp.StatusCode, 400)
	data, err = io.ReadAll(resp.Body)
	assert.Nil(err)
	assert.True(strings.Contains(string(data), "Missing/Invalid fields"), "Invalid response body", string(data))

	// Valid request
	req, err = http.NewRequest("POST", api.URL+"/rides", bytes.NewReader([]byte(fmt.Sprintf(`{
		"locationTo": "New York",
		"locationFrom": "San Francisco",
		"tackingPlaceAt": "2024-12-31T12:35:00+02:00",
		"driver": "nmBSHcxyvn",
		"transportLimit": 4
	}`))))
	assert.Nil(err)
	req.Header.Add("Authorization", "yr0osJ1-kQrQsOXzMGNVhbGUzdA0seGWAMK70WERFRWU5NzKcrQ1R2U_8ofXubwLbWxJYQK9hvj9xonabMMroA6oPjfnFuFR_zwOugdNGZVOwo6l8zczvFYRnGUdncOWv5Ckdy5eyB0leWpH7sDI_hbAxyiKljceGnKX-hcvB9MwjnsAiJMZ6EC_nAV-6ujEwM-YbPbYwndTEyY7CgDBp9gYrcOlvs9z_yf5sM_WQlziZFVVyGVoJyWDl-a1XbyLiagscmTeDs0pxQO0BH0oBF5qW8IRDIWAOuaSz3K9eygpqKQIxTFVq_psqaZT_qrhHI-3k-OPBbtWq9pF32-wVxNoFJMB3YvY17DgQfxxvzgckUH5YFlNks1cUgroHk2CIjtgs-9eskUzOrCzBKW3-EBcuyNrttnIePAkdVl2NC586fkBCVnKqfVIKYwm-ZrdCHxQVTZwGcswGnUP-YajlwZhmM-jgBjXIAJfWihcQTrDGmWz-0z8R8kycMdASguZXnQolGTvUOsOT21kFC4fwF-XQRi0tPh4mg0Bj1QN9y5sgibripVhCXQ7ma9QbbYL9ooAax6wqEU7b5-Gfai_r1ZLI5WcjOkI0ePAa2PikIC1b5nAMaz0c9y7Sv-hVAYtVzW5VB6PRJ4f5DoI_6KlGx6jE1AzmIEyPp_4_ImIUhHBlGUa7kikZkqUTtr9vSaz84EvQzT81wt3ULBLvA89Cr5rOWgAlNfmul3JZtJwfUuW39Mxc6QQN1mLUyKIUiofZImwkLqlACuriArAhMM_E8qo2V9sHSRVhZA_NOnKOYujsoFTTdr4vb2CWyeVIAEWT2YCueSMXinGL1Gmbxcczy9Hi2LoupnGYlQr9KgP5V_UrRvl_isC1MgUArQ25nIkdBNpUREW7a31bqWibAamOCgLP8bS20DERUD3-bKcDYDSDq9cEP2pKBRm_WyVQqCNYPIUpDOmDd9SEAZ3J_WveApSIDJlDt0j_nTibImctu6he92Kp63L5_A8nG6wBWW363CZ7tgktoY3KidPwbByX35BQRTUyE7wYxAqzdcF8Jd_n24SLHxC")
	resp, err = api.Client().Do(req)
	assert.Nil(err)
	assert.Eq(resp.StatusCode, 201)
	data, err = io.ReadAll(resp.Body)
	assert.Nil(err)
	var ride sqlc.Ride
	err = json.Unmarshal(data, &ride)
	assert.Nil(err)
	assert.Eq(ride.CreatedBy, "NnCaPHQLC9")
	assert.Eq(ride.Driver, "nmBSHcxyvn")
	assert.Eq(ride.LocationFrom, "San Francisco")
	assert.Eq(ride.TackingPlaceAt, "2024-12-31T10:35:00Z")
}

func TestHandleGetManyRides(t *testing.T) {
	db := utils.InitTestDB(path.Join(utils.ProjectRoot(), "db/testing/setup/0002-handle-get-many-rides-three-items.sql"))
	queries := sqlc.New(db)
	handler := rest.NewRESTApi(queries)

	api := httptest.NewServer(handler)
	defer api.Close()

	testAuth(api, "/rides/all", "GET")

	req, err := http.NewRequest("GET", api.URL+"/rides/all", bytes.NewReader([]byte{}))
	assert.Nil(err)
	req.Header.Add("Authorization", "yr0osJ1-kQrQsOXzMGNVhbGUzdA0seGWAMK70WERFRWU5NzKcrQ1R2U_8ofXubwLbWxJYQK9hvj9xonabMMroA6oPjfnFuFR_zwOugdNGZVOwo6l8zczvFYRnGUdncOWv5Ckdy5eyB0leWpH7sDI_hbAxyiKljceGnKX-hcvB9MwjnsAiJMZ6EC_nAV-6ujEwM-YbPbYwndTEyY7CgDBp9gYrcOlvs9z_yf5sM_WQlziZFVVyGVoJyWDl-a1XbyLiagscmTeDs0pxQO0BH0oBF5qW8IRDIWAOuaSz3K9eygpqKQIxTFVq_psqaZT_qrhHI-3k-OPBbtWq9pF32-wVxNoFJMB3YvY17DgQfxxvzgckUH5YFlNks1cUgroHk2CIjtgs-9eskUzOrCzBKW3-EBcuyNrttnIePAkdVl2NC586fkBCVnKqfVIKYwm-ZrdCHxQVTZwGcswGnUP-YajlwZhmM-jgBjXIAJfWihcQTrDGmWz-0z8R8kycMdASguZXnQolGTvUOsOT21kFC4fwF-XQRi0tPh4mg0Bj1QN9y5sgibripVhCXQ7ma9QbbYL9ooAax6wqEU7b5-Gfai_r1ZLI5WcjOkI0ePAa2PikIC1b5nAMaz0c9y7Sv-hVAYtVzW5VB6PRJ4f5DoI_6KlGx6jE1AzmIEyPp_4_ImIUhHBlGUa7kikZkqUTtr9vSaz84EvQzT81wt3ULBLvA89Cr5rOWgAlNfmul3JZtJwfUuW39Mxc6QQN1mLUyKIUiofZImwkLqlACuriArAhMM_E8qo2V9sHSRVhZA_NOnKOYujsoFTTdr4vb2CWyeVIAEWT2YCueSMXinGL1Gmbxcczy9Hi2LoupnGYlQr9KgP5V_UrRvl_isC1MgUArQ25nIkdBNpUREW7a31bqWibAamOCgLP8bS20DERUD3-bKcDYDSDq9cEP2pKBRm_WyVQqCNYPIUpDOmDd9SEAZ3J_WveApSIDJlDt0j_nTibImctu6he92Kp63L5_A8nG6wBWW363CZ7tgktoY3KidPwbByX35BQRTUyE7wYxAqzdcF8Jd_n24SLHxC")
	resp, err := api.Client().Do(req)
	assert.Nil(err)
	assert.Eq(resp.StatusCode, 200)
	data, err := io.ReadAll(resp.Body)
	assert.Nil(err)
	var rides []sqlc.Ride
	err = json.Unmarshal(data, &rides)
	assert.Nil(err)
	assert.Eq(len(rides), 3)
	assert.Eq(rides[0].LocationFrom, "Tokyo")
	assert.Eq(rides[0].CreatedBy, "nmBSHcxyvn")
	assert.Eq(rides[1].LocationFrom, "NYC")
}

func TestHandleGetManyRidesEmpty(t *testing.T) {
	db := utils.InitTestDB(path.Join(utils.ProjectRoot(), "db/testing/setup/0003-handle-get-many-rides-empty.sql"))
	queries := sqlc.New(db)
	handler := rest.NewRESTApi(queries)

	api := httptest.NewServer(handler)
	defer api.Close()

	testAuth(api, "/rides/all", "GET")

	req, err := http.NewRequest("GET", api.URL+"/rides/all", bytes.NewReader([]byte{}))
	assert.Nil(err)
	req.Header.Add("Authorization", "yr0osJ1-kQrQsOXzMGNVhbGUzdA0seGWAMK70WERFRWU5NzKcrQ1R2U_8ofXubwLbWxJYQK9hvj9xonabMMroA6oPjfnFuFR_zwOugdNGZVOwo6l8zczvFYRnGUdncOWv5Ckdy5eyB0leWpH7sDI_hbAxyiKljceGnKX-hcvB9MwjnsAiJMZ6EC_nAV-6ujEwM-YbPbYwndTEyY7CgDBp9gYrcOlvs9z_yf5sM_WQlziZFVVyGVoJyWDl-a1XbyLiagscmTeDs0pxQO0BH0oBF5qW8IRDIWAOuaSz3K9eygpqKQIxTFVq_psqaZT_qrhHI-3k-OPBbtWq9pF32-wVxNoFJMB3YvY17DgQfxxvzgckUH5YFlNks1cUgroHk2CIjtgs-9eskUzOrCzBKW3-EBcuyNrttnIePAkdVl2NC586fkBCVnKqfVIKYwm-ZrdCHxQVTZwGcswGnUP-YajlwZhmM-jgBjXIAJfWihcQTrDGmWz-0z8R8kycMdASguZXnQolGTvUOsOT21kFC4fwF-XQRi0tPh4mg0Bj1QN9y5sgibripVhCXQ7ma9QbbYL9ooAax6wqEU7b5-Gfai_r1ZLI5WcjOkI0ePAa2PikIC1b5nAMaz0c9y7Sv-hVAYtVzW5VB6PRJ4f5DoI_6KlGx6jE1AzmIEyPp_4_ImIUhHBlGUa7kikZkqUTtr9vSaz84EvQzT81wt3ULBLvA89Cr5rOWgAlNfmul3JZtJwfUuW39Mxc6QQN1mLUyKIUiofZImwkLqlACuriArAhMM_E8qo2V9sHSRVhZA_NOnKOYujsoFTTdr4vb2CWyeVIAEWT2YCueSMXinGL1Gmbxcczy9Hi2LoupnGYlQr9KgP5V_UrRvl_isC1MgUArQ25nIkdBNpUREW7a31bqWibAamOCgLP8bS20DERUD3-bKcDYDSDq9cEP2pKBRm_WyVQqCNYPIUpDOmDd9SEAZ3J_WveApSIDJlDt0j_nTibImctu6he92Kp63L5_A8nG6wBWW363CZ7tgktoY3KidPwbByX35BQRTUyE7wYxAqzdcF8Jd_n24SLHxC")
	resp, err := api.Client().Do(req)
	assert.Nil(err)
	assert.Eq(resp.StatusCode, 200)
	data, err := io.ReadAll(resp.Body)
	assert.Nil(err)
	var rides []sqlc.Ride
	err = json.Unmarshal(data, &rides)
	assert.Eq(len(rides), 0)
}

func testAuth(api *httptest.Server, endpoint string, method string) {
	// Missing authentication token
	req, err := http.NewRequest(method, api.URL+endpoint, bytes.NewReader([]byte{}))
	assert.Nil(err)
	resp, err := api.Client().Do(req)
	assert.Nil(err)
	assert.Eq(resp.StatusCode, 400)

	// Invalid token
	req, err = http.NewRequest(method, api.URL+endpoint, bytes.NewReader([]byte{}))
	assert.Nil(err)
	req.Header.Add("Authorization", "invalid-token")
	resp, err = api.Client().Do(req)
	assert.Nil(err)
	assert.Eq(resp.StatusCode, 401)
}
