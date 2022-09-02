package endpoints

import (
	"fmt"
	"net/http"

	"cms.csesoc.unsw.edu.au/database/repositories"
	. "cms.csesoc.unsw.edu.au/endpoints/models"
	"cms.csesoc.unsw.edu.au/internal/logger"
	"github.com/google/uuid"
)

// Defines endpoints consumable via the API
func GetEntityInfo(form ValidInfoRequest, df DependencyFactory) handlerResponse[EntityInfoResponse] {
	log := getDependency[*logger.Log](df)
	fsRepo := getDependency[repositories.IFilesystemRepository](df)

	// Query the repository for an existing entity with the given ID
	entity, err := fsRepo.GetEntryWithID(form.EntityID)
	if err != nil {
		return handlerResponse[EntityInfoResponse]{
			Status:   http.StatusNotFound,
			Response: EntityInfoResponse{},
		}
	}

	log.Write(fmt.Sprintf("retrieved entity: %v.", entity))

	return handlerResponse[EntityInfoResponse]{
		Status:   http.StatusOK,
		Response: FsEntryToEntityInfo(entity, fsRepo, true),
	}
}

// CreateNewEntity is the public handler for constructing and creating new entities
func CreateNewEntity(form ValidEntityCreationRequest, df DependencyFactory) handlerResponse[NewEntityResponse] {
	fsRepo := getDependency[repositories.IFilesystemRepository](df)
	pubRepo := getDependency[repositories.IUnpublishedVolumeRepository](df)
	log := getDependency[*logger.Log](df)

	entityToCreate := CreationReqToFsEntry(form)
	newEntity, err := fsRepo.CreateEntry(entityToCreate)
	if err != nil {
		return handlerResponse[NewEntityResponse]{
			Status: http.StatusNotAcceptable,
		}
	}

	log.Write(fmt.Sprintf("created new entity %v.", entityToCreate))
	pubRepo.AddToVolume(newEntity.EntityID.String())
	return handlerResponse[NewEntityResponse]{
		Status:   http.StatusOK,
		Response: NewEntityResponse{NewID: newEntity.EntityID},
	}
}

// Handler for deleting filesystem entities
func DeleteFilesystemEntity(form ValidInfoRequest, df DependencyFactory) handlerResponse[empty] {
	fsRepo := getDependency[repositories.IFilesystemRepository](df)
	log := getDependency[*logger.Log](df)

	err := fsRepo.DeleteEntryWithID(form.EntityID)
	if err != nil {
		return handlerResponse[empty]{
			Status: http.StatusNotAcceptable,
		}
	}

	log.Write(fmt.Sprintf("deleted entity with ID: %s", form.EntityID))
	return handlerResponse[empty]{
		Status: http.StatusOK,
	}
}

// Handler for retrieving children
func GetChildren(form ValidInfoRequest, df DependencyFactory) handlerResponse[ChildrenRequestResponse] {
	fsRepo := getDependency[repositories.IFilesystemRepository](df)
	log := getDependency[*logger.Log](df)

	fileInfo, err := fsRepo.GetEntryWithID(form.EntityID)
	if err != nil {
		return handlerResponse[ChildrenRequestResponse]{
			Status: http.StatusNotFound,
		}
	}

	log.Write(fmt.Sprintf("fetched children for %s, got %v.", form.EntityID, fileInfo.ChildrenIDs))
	return handlerResponse[ChildrenRequestResponse]{
		Status: http.StatusOK,
		Response: ChildrenRequestResponse{
			Children: fileInfo.ChildrenIDs,
		},
	}
}

func GetIDWithPath(form ValidPathRequest, df DependencyFactory) handlerResponse[uuid.UUID] {
	repository := getDependency[repositories.IFilesystemRepository](df)
	log := getDependency[*logger.Log](df)

	entityID, err := repository.GetIDWithPath(form.Path)
	if err != nil {
		return handlerResponse[uuid.UUID]{
			Status: http.StatusNotFound,
		}
	}

	log.Write(fmt.Sprintf("got ID %s for %s", entityID, form.Path))
	return handlerResponse[uuid.UUID]{
		Status: http.StatusOK, Response: entityID,
	}
}

// Handler for renaming filesystem entities
func RenameFilesystemEntity(form ValidRenameRequest, df DependencyFactory) handlerResponse[empty] {
	repository := getDependency[repositories.IFilesystemRepository](df)
	err := repository.RenameEntity(form.EntityID, form.NewName)
	if err != nil {
		return handlerResponse[empty]{Status: http.StatusNotAcceptable}
	}

	return handlerResponse[empty]{Status: http.StatusOK}
}
