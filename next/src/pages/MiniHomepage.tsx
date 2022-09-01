import styled from "styled-components";
import Image from 'next/image';
import Navbar from "../components/navbar/Navbar";
import { NavbarOpenHandler } from "../components/navbar/types";
import HamburgerMenu from "../components/navbar/HamburgerMenu";
import { device } from "../styles/device"
import { useState } from "react";

type Props = {};

const HomepageContainer = styled.div`
	display: flex;
	flex-direction: column;
	height: 100vh;
	width: 100%;
`;
const Container = styled.div`
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  @media ${device.laptop} {
    flex-direction: row;
  }
`;

const ColumnContainer = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  @media ${device.desktop} {
    padding: 0 10vw;
  }

  @media ${device.laptop} {

  }
`;

const ImageContainer = styled.div`
	position: relative; /* IMPORTANT */
	width: 550px;
	height: 400px;
	margin-left: 100px;
	display: flex;
	align-items: center;
`;

const Text1 = styled.p`
	color: #010033;
	font-size: 36px;
	padding: 10px 0;
	margin-top: 100px;
`;

const Text2 = styled.p`
	color: #3977f8;
	font-size: 36px;
`;

const Text3 = styled.p`
	color: #010033;
	font-style: italic;
	font-size: 36px;
	display: inline;
`;

const Scroll = styled.p`
	transform: rotate(90deg);
	position: absolute;
	right: 0px;
	bottom: 10%;
`;

export default function Homepage({}: Props) {
  return (
    <>
      <HomepageContainer>
        <Container>
          <ColumnContainer>
            <Image src="/assets/logo.svg" width="600px" height="300px"/>
            <Text1>
              Empowering
              <Text3> future</Text3>
              <Text2>Technological Leaders</Text2>
            </Text1>
          </ColumnContainer>
          {/* <Button>Visit on Blog</Button> */}
          <ImageContainer>
            <Image src="/assets/WebsitesIcon.png" layout="fill" objectFit="contain" />
          </ImageContainer>
        </Container>
      </HomepageContainer>
    </>
  )
}
