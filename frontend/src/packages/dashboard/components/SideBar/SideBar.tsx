import React from "react";
import styled from "styled-components";
import Button from "@mui/material/Button";
import { useNavigate } from "react-router-dom";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import MenuIcon from "@mui/icons-material/Menu";

const Container = styled.div`
  position: relative;
  width: 250px;
  background: #c9bff2;
  height: 100vh;
  transition: left 0.3s ease-in-out;
  margin-right: 25px;
`;

const IconWrapper = styled.div`
  z-index: 999;
  position: absolute;
  top: 50%;
  transform: translate(0, -50%);

  right: -45px;
  width: 30px;
  height: 30px;
  cursor: pointer;

  display: flex;
  justify-content: center;
  align-items: center;
`;

const Circle = styled.div`
  position: absolute;
  border: 1px solid black;
  border-radius: 999px;
  width: 40px;
  height: 40px;
  box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.25);
`;

const SidebarTitle = styled.div`
  font-size: xx-large;
  margin: 2rem;
  font-weight: bold;
`;

const ButtonFlex = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  grid-gap: 80px;
`;

const ButtonGroup = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  grid-gap: 30px;
`;

interface SideBarButtonProps {
  bgcolor: string;
}

const SidebarButton = styled(Button)<SideBarButtonProps>`
  && {
    width: 160px;
    variant: contained;
    background-color: ${(props) => props.bgcolor};
    color: white;
    border-radius: 20px;
    text-transform: none;
    color: black;
    &:hover {
      transform: scale(1.04);
      background-color: darkgrey;
    }
    &:active {
      transform: scale(0.96);
      background-color: darkgrey;
    }
  }
`;

type Props = {
  setModalState: (state: { open: boolean; type: string }) => void;
  selectedFile: string | null;
  isOpen: boolean;
  setOpen: (state: boolean) => void;
};

// Wrapper component ${props => props.color}
export default function SideBar({
  setModalState,
  selectedFile,
  isOpen,
  setOpen,
}: Props) {
  const handleNewFile = () => {
    setModalState({
      open: true,
      type: "file",
    }); // sets modal to be open
  };

  const handleNewFolder = () => {
    setModalState({
      open: true,
      type: "folder",
    });
  };

  const navigate = useNavigate();
  const handleEdit = () => {
    if (selectedFile !== null) {
      navigate("/editor/" + selectedFile, { replace: false }), [navigate];
    }
  };

  // TODO
  const handleRecycle = () => {
    return;
  };

  return (
    <Container style={{ left: isOpen ? "0px" : "-250px" }}>
      <IconWrapper onClick={() => setOpen(!isOpen)}>
        <Circle />
        {isOpen ? <ArrowBackIcon /> : <MenuIcon />}
      </IconWrapper>
      <SidebarTitle>Welcome \name\</SidebarTitle>
      <ButtonFlex>
        {/* <ButtonGroup>
          <SidebarButton bgcolor="#F88282">
            Blog
          </SidebarButton>
          <SidebarButton bgcolor="#F88282">
            Core pages
          </SidebarButton>
        </ButtonGroup> */}
        <ButtonGroup>
          <SidebarButton
            bgcolor="#b4c6ff"
            onClick={handleNewFile}
            data-anchor="NewPageButton"
          >
            New page
          </SidebarButton>
          <SidebarButton
            bgcolor="#b4c6ff"
            onClick={handleNewFolder}
            data-anchor="NewFolderButton"
          >
            New folder
          </SidebarButton>
        </ButtonGroup>
        <ButtonGroup>
          <SidebarButton bgcolor="white" onClick={handleEdit}>
            Edit
          </SidebarButton>
          {/* <SidebarButton bgcolor="#B8E8E8">
            Feature
          </SidebarButton>
          <SidebarButton bgcolor="#B8E8E8" onClick={handleRecycle}>
            Recycle
          </SidebarButton> */}
        </ButtonGroup>
      </ButtonFlex>
    </Container>
  );
}
