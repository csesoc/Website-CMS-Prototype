import { describe } from "@jest/globals";
import { configureApiUrl, FilesystemAPI, resetApiUrl } from "../filesystem";
import { CreateFilesystemEntryResponse, FilesystemEntry, IsCreateFilesystemEntryResponse, IsFilesystemEntry } from "../types/filesystem";
import { APIError, IsEmptyApiResponse } from "../types/general";

// filesystemConsistencyTests ensure that the contract maintained between the frontend and backend regarding endpoint input/response types are consistent
//  note: requires the BE container to be up and running
beforeAll(() => {
    if (process.env.E2E_MODE === "github") {
        configureApiUrl("http://backend:8080");
    } else {
        configureApiUrl("http://localhost:8080")
    }
});

afterAll(() => {
    resetApiUrl();
})
    
    
// Ensure that the returned response was actually indeed assignable to FilesystemEntry
//  - We need this additional check as the return value is merely a "promise" to the typescript compiler that the response
//    conforms to a specific format, we need to additionally confirm this at runtime using a test
// All of these tests basically just assert that AT RUNTIME the API response types are indeed assignable to our type definitions
describe("the filesystem api should", () => {
    test("structure root info responses properly", async () => {
        const rootInformation = await FilesystemAPI.GetRootInfo();
        expect(IsFilesystemEntry(rootInformation), 'Expected root information to be assignable to the FilesystemEntry type').toBe(true);
    });

    // Note: that this is basically a full integration test of our BE
    test("structure base API responses properly", async () => {
        const root = (await FilesystemAPI.GetRootInfo()) as FilesystemEntry;

        // Create a document
        const newDocument = await FilesystemAPI.CreateDocument("ebic document of truth", root.EntityID);
        expect(IsCreateFilesystemEntryResponse(newDocument), "Expected CreateDocument response to be assignable to CreateFilesystemEntryResponse").toBe(true);

        // fetch the information
        const newEntityId = (newDocument as CreateFilesystemEntryResponse).NewID;
        const documentInformation = await FilesystemAPI.GetEntityInfo(newEntityId);
        expect(IsFilesystemEntry(documentInformation), 'Expected document information to be assignable to the FilesystemEntry type').toBe(true);

        // rename it
        const renameResp = await FilesystemAPI.RenameEntity(newEntityId, "docMcStuffins");
        expect(IsEmptyApiResponse(renameResp), 'Expected deletion response to be assignable to empty');

        // publish it
        const publishResp = await FilesystemAPI.RenameEntity(newEntityId, "docMcStuffins");
        expect(IsEmptyApiResponse(publishResp), 'Expected deletion response to be assignable to empty');

        // delete it
        const deletionResp = await FilesystemAPI.DeleteEntity(newEntityId);
        expect(IsEmptyApiResponse(deletionResp), 'Expected deletion response to be assignable to empty');
    });
})