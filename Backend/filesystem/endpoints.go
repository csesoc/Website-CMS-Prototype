package filesystem

import (
	"DiffSync/database"
	"DiffSync/httpUtil"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var httpDbPool database.Pool

// Todo: abstract out configuration logic elsewhere
func init() {
	var err error
	httpDbPool, err = database.NewPool(database.Config{
		HostAndPort: "db:5432",
		User:        "postgres",
		Password:    "postgres",
		Database:    "test_db",
	})

	if err != nil {
		log.Print(err.Error())
	}
}

type ValidInfoRequest struct {
	EntityID int `schema:"EntityID,required"`
}

// Defines endpoints consumable via the API
func GetEntityInfo(w http.ResponseWriter, r *http.Request) {
	routes := strings.Split(r.URL.RequestURI(), "/")

	switch routes[len(routes)-1] {
	case "root":
		fileInfo, err := getRootInfo(httpDbPool)
		if err != nil {
			httpUtil.ThrowRequestError(w, 500, "something went wrong")
			return
		}

		out, _ := json.Marshal(fileInfo)
		httpUtil.SendResponse(w, string(out))

		break
	default:
		var input ValidInfoRequest
		if validRequest := httpUtil.ParseParamsToSchema(w, r, []string{"GET"}, map[int]string{
			400: "missing EntityID paramater",
			405: "invalid method",
		}, &input); validRequest {

			fileInfo, err := getFilesystemInfo(httpDbPool, input.EntityID)
			if err != nil {
				httpUtil.ThrowRequestError(w, 404, "unable to find entity with requested ID")
				return
			}

			out, _ := json.Marshal(fileInfo)
			httpUtil.SendResponse(w, string(out))
		}
		break
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
			newID, err = createFilesystemEntityAtRoot(httpDbPool, input.LogicalName, input.OwnerGroup, input.IsDocument)
		} else {
			log.Print("hello there\n")
			newID, err = createFilesystemEntity(httpDbPool, input.Parent, input.LogicalName, input.OwnerGroup, input.IsDocument)
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
		err := deleteEntity(httpDbPool, input.EntityID)
		if err != nil {
			httpUtil.ThrowRequestError(w, 500, "unable to delete, the requested entity is either the root directory or has children")
		} else {
			httpUtil.SendResponse(w, fmt.Sprintf(`{"success": true, "deleted": %d}`, input.EntityID))
		}
	}
}