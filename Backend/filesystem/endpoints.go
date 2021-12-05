package filesystem

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cms.csesoc.unsw.edu.au/database"
	"cms.csesoc.unsw.edu.au/environment"
	"cms.csesoc.unsw.edu.au/internal/httpUtil"
)

var httpDBContext database.LiveContext

func init() {
	if !environment.IsTestingEnvironment() {
		var err error
		httpDBContext, err = database.NewLiveContext()
		if err != nil {
			log.Print(err.Error())
		}
	}
}

type ValidInfoRequest struct {
	EntityID int `schema:"EntityID"`
}

// Defines endpoints consumable via the API
func GetEntityInfo(w http.ResponseWriter, r *http.Request) {

	var input ValidInfoRequest
	if validRequest := httpUtil.ParseParamsToSchema(w, r, []string{"GET"}, map[int]string{
		400: "missing EntityID paramater",
		405: "invalid method",
	}, &input); validRequest {
		var fileInfo EntityInfo
		var err error

		if input.EntityID == 0 {
			fileInfo, err = GetRootInfo(httpDBContext)
		} else {
			fileInfo, err = GetFilesystemInfo(httpDBContext, input.EntityID)
		}

		if err != nil {
			httpUtil.ThrowRequestError(w, 404, "unable to find entity with requested ID")
			return
		}

		out, _ := json.Marshal(fileInfo)
		httpUtil.SendResponse(w, string(out))
	}

}

// TODO: this needs to be wrapped around auth and permissions later
type ValidEntityCreationRequest struct {
	Parent      int
	LogicalName string `schema:"LogicalName,required"`
	OwnerGroup  int    `schema:"OwnerGroup,required"`
	IsDocument  bool   `schema:"IsDocument,required"`
}

func CreateNewEntity(w http.ResponseWriter, r *http.Request) {

	var input ValidEntityCreationRequest
	if validRequest := httpUtil.ParseParamsToSchema(w, r, []string{"POST"}, map[int]string{
		400: "missing paramaters, must have: LogicalName, OwnerGroup, IsDocument",
		405: "invalid method",
	}, &input); validRequest {
		var newID int
		var err error

		if input.Parent == 0 {
			newID, err = CreateFilesystemEntityAtRoot(httpDBContext, input.LogicalName, input.OwnerGroup, input.IsDocument)
		} else {
			newID, err = CreateFilesystemEntity(httpDBContext, input.Parent, input.LogicalName, input.OwnerGroup, input.IsDocument)
		}

		if err != nil {
			httpUtil.ThrowRequestError(w, 500, "unable to create entity (may be a duplicate)")
		} else {
			httpUtil.SendResponse(w, fmt.Sprintf(`{"success": true, "newID": %d}`, newID))
		}
	}

}

// Handler for deleting filesystem entities
func DeleteFilesystemEntity(w http.ResponseWriter, r *http.Request) {
	var input ValidInfoRequest
	if validRequest := httpUtil.ParseParamsToSchema(w, r, []string{"POST"}, map[int]string{
		400: "missing paramaters, must have: LogicalName, OwnerGroup, IsDocument",
		405: "invalid method",
	}, &input); validRequest {
		err := DeleteEntity(httpDBContext, input.EntityID)
		if err != nil {
			httpUtil.ThrowRequestError(w, 500, "unable to delete, the requested entity is either the root directory or has children")
		} else {
			httpUtil.SendResponse(w, fmt.Sprintf(`{"success": true, "deleted": %d}`, input.EntityID))
		}
	}
}

// Handler for retrieving children
func GetChildren(w http.ResponseWriter, r *http.Request) {
	var input ValidInfoRequest
	if validRequest := httpUtil.ParseParamsToSchema(w, r, []string{"GET"}, map[int]string{
		400: "missing EntityID paramater",
		405: "invalid method",
	}, &input); validRequest {

		fileInfo, err := GetEntityChildren(httpDBContext, input.EntityID)
		if err != nil {
			httpUtil.ThrowRequestError(w, 404, "unable to find entity with requested ID")
			return
		}

		out, _ := json.Marshal(fileInfo)
		httpUtil.SendResponse(w, string(out))
	}
}

type ValidRenameRequest struct {
	EntityID int    `schema:"EntityID,required"`
	NewName  string `schema:"NewName,required"`
}

// Handler for renaming filesystem entities
func RenameFilesystemEntity(w http.ResponseWriter, r *http.Request) {
	var input ValidRenameRequest
	if validRequest := httpUtil.ParseParamsToSchema(w, r, []string{"POST"}, map[int]string{
		400: "missing paramaters, must have: NewName, EntityID",
		405: "invalid method",
	}, &input); validRequest {
		err := RenameEntity(httpDBContext, input.EntityID, input.NewName)
		if err != nil {
			httpUtil.ThrowRequestError(w, 500, "unable rename, the requested name is most likely taken")
		} else {
			httpUtil.SendResponse(w, fmt.Sprintf(`{"success": true, "renamed": %d}`, input.EntityID))
		}
	}
}
