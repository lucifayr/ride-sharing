package rest

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"
	"strconv"
)

func groupHandlers(h *http.ServeMux) {
	h.HandleFunc("POST /groups", handle(createGroup).with(bearerAuth(false)).build())
	h.HandleFunc("POST /groups/update", handle(updateGroup).with(bearerAuth(false)).build())
	h.HandleFunc("GET /groups/many", handle(getManyGroups).with(bearerAuth(false)).build())
	h.HandleFunc("GET /groups/by-id/{id}", handle(getGroupById).with(bearerAuth(false)).build())
}

type GroupData struct {
	GroupId     string  `json:"groupId"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatedBy   string  `json:"createdBy"`
}

type createGroupParams struct {
	Name        *string `json:"name" validate:"required"`
	Description *string `json:"description"`
}

type updateGroupParams struct {
	GroupId     *string `json:"groupId" validate:"required"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func createGroup(w http.ResponseWriter, r *http.Request) {
	user := getMiddlewareData[sqlc.User](r, "user")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error: Invalid request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid request body.")
		return
	}

	var createParams createGroupParams
	err = json.Unmarshal(data, &createParams)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid JSON in request body.", err.Error())
		return
	}

	err = utils.Validate.Struct(createParams)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Missing/Invalid fields in request body.", err.Error())
		return
	}

	desc := ""
	if createParams.Description != nil {
		desc = *createParams.Description
	}

	argsCreateGroup := sqlc.GroupsCreateParams{
		Name: *createParams.Name,
		Description: sql.NullString{
			String: desc,
			Valid:  createParams.Description != nil,
		},
		CreatedBy: user.ID,
	}
	group, err := state.queries.GroupsCreate(r.Context(), argsCreateGroup)

	var dataDesc *string = nil
	if group.Description.Valid {
		dataDesc = &group.Description.String
	}

	groupData := GroupData{
		GroupId:     group.ID,
		Name:        group.Name,
		Description: dataDesc,
		CreatedBy:   group.CreatedBy,
	}

	resp, err := json.Marshal(groupData)
	assert.Nil(err, "Failed to serialize ride.")
	w.WriteHeader(201)
	w.Write(resp)
}

func updateGroup(w http.ResponseWriter, r *http.Request) {
	user := getMiddlewareData[sqlc.User](r, "user")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error: Invalid request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid request body.")
		return
	}

	var updateParams updateGroupParams
	err = json.Unmarshal(data, &updateParams)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid JSON in request body.", err.Error())
		return
	}

	err = utils.Validate.Struct(updateParams)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Missing/Invalid fields in request body.", err.Error())
		return
	}

	group, err := state.queries.GroupsGetById(r.Context(), *updateParams.GroupId)
	if errors.Is(err, sql.ErrNoRows) {
		httpWriteErr(w, http.StatusNotFound, "No group exists with 'id'.")
		return
	}

	if group.CreatedBy != user.ID {
		httpWriteErr(w, http.StatusBadRequest, "You are not the owner of this group.")
		return
	}

	if updateParams.Name != nil {
		argsUpdateName := sqlc.GroupsUpdateNameParams{
			Name: *updateParams.Name,
			ID:   *updateParams.GroupId,
		}

		err := state.queries.GroupsUpdateName(r.Context(), argsUpdateName)
		assert.Nil(err)
	}

	if updateParams.Description != nil {
		argsUpdateDescription := sqlc.GroupsUpdateDescriptionParams{
			Description: utils.SqlNullStr(*updateParams.Description),
			ID:          *updateParams.GroupId,
		}

		err := state.queries.GroupsUpdateDescription(r.Context(), argsUpdateDescription)
		assert.Nil(err)
	}
}

func getGroupById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httpWriteErr(w, http.StatusBadRequest, "Must provide 'id' path parameter.")
		return
	}

	group, err := state.queries.GroupsGetById(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		httpWriteErr(w, http.StatusNotFound, "No group exists with 'id'.")
		return
	}

	var dataDesc *string = nil
	if group.Description.Valid {
		dataDesc = &group.Description.String
	}

	groupData := GroupData{
		GroupId:     group.ID,
		Name:        group.Name,
		Description: dataDesc,
		CreatedBy:   group.CreatedBy,
	}

	var resp []byte
	resp, err = json.Marshal(groupData)
	assert.Nil(err, "Failed to serialize group.")
	w.WriteHeader(200)
	w.Write(resp)
}

func getManyGroups(w http.ResponseWriter, r *http.Request) {
	var offset int64 = 0
	offsetStr := r.FormValue("offset")
	if parsed, err := strconv.ParseInt(offsetStr, 10, 64); err == nil && parsed > 0 {
		offset = parsed
	}

	rows, err := state.queries.GroupsGetMany(r.Context(), offset)
	if err != nil {
		log.Println("Error: Failed to get groups.", "error:", err)
		httpWriteErr(w, http.StatusInternalServerError, "Failed to get groups.")
		return
	}

	groups := make([]GroupData, len(rows))
	for idx, row := range rows {

		var desc *string = nil
		if row.Description.Valid {
			desc = &row.Description.String
		}

		groups[idx] = GroupData{
			GroupId:     row.ID,
			Name:        row.Name,
			Description: desc,
			CreatedBy:   row.CreatedBy,
		}
	}

	var resp []byte
	if len(groups) == 0 {
		resp = []byte("[]")
	} else {
		resp, err = json.Marshal(groups)
	}

	assert.Nil(err, "Failed to serialize groups.")
	w.WriteHeader(200)
	w.Write(resp)
}
