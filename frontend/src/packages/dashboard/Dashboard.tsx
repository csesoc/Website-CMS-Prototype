import React, { useEffect, useState } from "react";
import { useDispatch } from "react-redux";
import styled from "styled-components";
import { Breadcrumbs, Link } from "@mui/material";

// local imports
import SideBar from "src/packages/dashboard/components/SideBar/SideBar";
import Renderer from "./components/FileRenderer/Renderer";
import { initAction, traverseBackFolder } from "./state/folders/actions";
import ConfirmationWindow from "./components/ConfirmationModal/ConfirmationWindow";
import Directory from "./components/Directory";

const Container = styled.div`
  display: flex;
  flex-direction: row;
`;

const FlexWrapper = styled.div`
  position: relative;
  display: flex;
  flex-direction: column;
  transition: left 0.3s ease-in-out;
`;

export default function Dashboard() {
  const [isOpen, setOpen] = useState(true);

  const [modalState, setModalState] = useState<{ open: boolean; type: string }>(
    {
      open: false,
      type: "",
    }
  );

  const [selectedFile, setSelectedFile] = useState<string | null>(null);
  // const parentFolder = getFolderState().parentFolder;
  const dispatch = useDispatch();

  useEffect(() => {
    // fetches all folders and files from backend and displays it
    dispatch(initAction());
  }, []);

  return (
    <Container>
      <SideBar
        setModalState={setModalState}
        selectedFile={selectedFile}
        isOpen={isOpen}
        setOpen={setOpen}
      />
      <FlexWrapper style={{ left: isOpen ? "0px" : "-250px" }}>
        <Directory />
        <Renderer
          selectedFile={selectedFile}
          setSelectedFile={setSelectedFile}
        />
      </FlexWrapper>
      <ConfirmationWindow
        open={modalState.open}
        modalState={modalState}
        setModalState={setModalState}
      />
    </Container>
  );
}
