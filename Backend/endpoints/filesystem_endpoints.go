package endpoints

import (
	"net/http"
	"reflect"

	"cms.csesoc.unsw.edu.au/database/repositories"
)

type EntityInfo struct {
	EntityID   int
	EntityName string
	IsDocument bool
	Parent     int
	Children   []EntityInfo
}

type ValidInfoRequest struct {
	EntityID int `schema:"EntityID"`
}

// Defines endpoints consumable via the API
func GetEntityInfo(w http.ResponseWriter, r *http.Request, df DependencyFactory) (int, interface{}, error) {
	var input ValidInfoRequest
	if status := ParseParamsToSchema(r, "GET", &input); status != http.StatusOK {
		return status, nil, nil
	}

	fs := reflect.TypeOf((*repositories.IFilesystemRepository)(nil))
	repository := df.GetDependency(fs).(repositories.IFilesystemRepository)

	if entity, err := repository.GetEntryWithID(input.EntityID); err == nil {
		childrenArr := []EntityInfo{}
		for _, id := range entity.ChildrenIDs {
			x, _ := repository.GetEntryWithID(id)
			childrenArr = append(childrenArr, EntityInfo{
				EntityID:   id,
				EntityName: x.LogicalName,
				IsDocument: x.IsDocument,
				Parent:     x.ParentFileID,
				Children:   nil,
			})
		}

		return http.StatusOK, EntityInfo{
			EntityID:   entity.EntityID,
			EntityName: entity.LogicalName,
			IsDocument: entity.IsDocument,
			Parent:     entity.ParentFileID,
			Children:   childrenArr,
		}, nil
	} else {
		return http.StatusNotFound, nil, nil
	}
}

// TODO: this needs to be wrapped around auth and permissions later
type ValidEntityCreationRequest struct {
	Parent      int
	LogicalName string `schema:"LogicalName,required"`
	OwnerGroup  int    `schema:"OwnerGroup,required"`
	IsDocument  bool   `schema:"IsDocument,required"`
}

func CreateNewEntity(w http.ResponseWriter, r *http.Request, df DependencyFactory) (int, interface{}, error) {
	var input ValidEntityCreationRequest
	if status := ParseParamsToSchema(r, "POST", &input); status != http.StatusOK {
		return status, nil, nil
	}

	fs := reflect.TypeOf((*repositories.IFilesystemRepository)(nil))
	repository := df.GetDependency(fs).(repositories.IFilesystemRepository)
	entityToCreate := repositories.FilesystemEntry{
		LogicalName:  input.LogicalName,
		ParentFileID: input.Parent,
		IsDocument:   input.IsDocument,
		OwnerUserId:  input.OwnerGroup,
	}

	if e, err := repository.CreateEntry(entityToCreate); err != nil {
		return http.StatusNotAcceptable, nil, err
	} else {
		return http.StatusOK, struct {
			NewID int
		}{
			NewID: e.EntityID,
		}, nil
	}

}

// Handler for deleting filesystem entities
func DeleteFilesystemEntity(w http.ResponseWriter, r *http.Request, df DependencyFactory) (int, interface{}, error) {
	var input ValidInfoRequest
	if status := ParseParamsToSchema(r, "POST", &input); status != http.StatusOK {
		return status, nil, nil
	}

	fs := reflect.TypeOf((*repositories.IFilesystemRepository)(nil))
	repository := df.GetDependency(fs).(repositories.IFilesystemRepository)
	if err := repository.DeleteEntryWithID(input.EntityID); err != nil {
		return http.StatusNotAcceptable, nil, err
	} else {
		return http.StatusOK, nil, nil
	}
}

// Handler for retrieving children
func GetChildren(w http.ResponseWriter, r *http.Request, df DependencyFactory) (int, interface{}, error) {
	var input ValidInfoRequest
	if status := ParseParamsToSchema(r, "GET", &input); status != http.StatusOK {
		return status, nil, nil
	}

	fs := reflect.TypeOf((*repositories.IFilesystemRepository)(nil))
	repository := df.GetDependency(fs).(repositories.IFilesystemRepository)
	if fileInfo, err := repository.GetEntryWithID(input.EntityID); err != nil {
		return http.StatusNotFound, nil, nil
	} else {
		return http.StatusOK, struct {
			Children []int
		}{
			Children: fileInfo.ChildrenIDs,
		}, nil
	}
}

type ValidPathRequest struct {
	Path string `schema:"Path,required"`
}

func GetIDWithPath(w http.ResponseWriter, r *http.Request, df DependencyFactory) (int, interface{}, error) {
	var input ValidPathRequest
	if status := ParseParamsToSchema(r, "GET", &input); status != http.StatusOK {
		return status, nil, nil
	}

	fs := reflect.TypeOf((*repositories.IFilesystemRepository)(nil))
	repository := df.GetDependency(fs).(repositories.IFilesystemRepository)
	if entityID, err := repository.GetIDWithPath(input.Path); err != nil {
		return http.StatusNotFound, nil, nil
	} else {
		return http.StatusOK, entityID, nil
	}
}

type ValidRenameRequest struct {
	EntityID int    `schema:"EntityID,required"`
	NewName  string `schema:"NewName,required"`
}

// Handler for renaming filesystem entities
func RenameFilesystemEntity(w http.ResponseWriter, r *http.Request, df DependencyFactory) (int, interface{}, error) {
	var input ValidRenameRequest
	if status := ParseParamsToSchema(r, "POST", &input); status != http.StatusOK {
		return status, nil, nil
	}

	fs := reflect.TypeOf((*repositories.IFilesystemRepository)(nil))
	repository := df.GetDependency(fs).(repositories.IFilesystemRepository)
	if err := repository.RenameEntity(input.EntityID, input.NewName); err != nil {
		return http.StatusNotAcceptable, nil, nil
	} else {
		return http.StatusOK, nil, nil
	}
}
