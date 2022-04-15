import { render } from "src/cse-testing-lib"
import { queryByDataAnchor } from "src/cse-testing-lib/custom-queries";
import SideBar from "./SideBar"

describe("Side bar tests", () => {
  it("Side bar is rendered with proper buttons", () => {
    const mockSetModalState = jest.fn();
    const mockSelectedFileID = 5;
    const { queryByDataAnchor } = render(
      <SideBar 
        setModalState={mockSetModalState}
        selectedFile={mockSelectedFileID}
      />
    )
    
    expect(queryByDataAnchor("NewPageButton")).toBeTruthy();
    expect(queryByDataAnchor("NewFolderButton")).toBeTruthy();
  })
})