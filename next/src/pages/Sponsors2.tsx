import React, { useState } from "react";
import Image from "next/image";
import Footer from "../components/footer/Footer";
import Navbar from "../components/navbar/Navbar";
import { NavbarOpenHandler, NavbarType } from "../components/navbar/types";
import styled from 'styled-components'
import Link from 'next/link'
import { device } from '../styles/device'


import Otter from '../svgs/otter.png'

const Text = styled.p`
  color: white;
  @media ${device.tablet} {
    padding: 3vh 0;
  }
`

const Grid = styled.div`
  display: flex;
  justify-content: center;
  margin-left: auto;
  margin-right: auto;
  max-width: 70vw;
  max-width: 75vw;
  border-radius: 0.5rem;
  margin-bottom: 10vh;
  background: radial-gradient(50% 50% at 50% 50%, rgba(146, 67, 166, 0.2407) 0%, rgba(119, 158, 237, 0.83) 100%);
 
`

const SmallGrid = styled.div`
  display: flex;
  justify-content: center;
  flex-direction: column;
  flex: 2;
  @media ${device.tablet} {
    min-height: 100%;
    flex-wrap: wrap;
    flex-direction: row;
    flex: 3;
  }
`

const OurSponsorsCol = styled.div`
  flex: 2;
  padding: 1.5rem;
  
`

const SponsorCol = styled.div`
  flex: 3;
  display: flex; 
  justify-content: center;
  max-height: 100%;
  max-width: 100%;
  border-radius: 0.5rem;
  padding: 1vw;
  @media ${device.tablet} {
    flex-basis: calc(50% - 40px);  
    flex-direction: column;
    padding: 8vh 0;
  }


`

const TextContainer = styled.div`
  display:flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
`

const ButtonContainer = styled.div`
  display: flex;
`

const H3 = styled.div`
  color: white;
  font-family: 'Raleway';
  font-weight: 800;
  font-size: 20px;
  @media ${device.tablet} {
    font-size: 2vw;
  }
`;


const Button = styled.button`
  background-color: white;
  margin: 10px;
  padding: 1vh 2vw;
  font-size: 18px;
  color: #817fff;
  border: none;
  border-radius: 6px;
  cursor: pointer;

  &:hover {
    transform: scale(1.05);
  }
  @media ${device.tablet} {
    padding: 10px 45px;
  }
`

const ImgContainer = styled.div`

  display: block;
  margin-left: auto;
  margin-right: auto;
  width: 30vw;

  @media ${device.laptop} {
    width: 350px;
    height: 400px;
  }

`

const HeadingContainer = styled.div`
  display: flex;
  justify-content: center;
`

const Heading = styled.div`
  color: #A09FE3;
  font-family: 'Raleway';
  font-weight: 800;
  font-size: 30px;
  padding: 5vh 0;
  @media ${device.tablet} {
    font-size: 2.8vw;
  }
`;



export default function Sponsors2() {

  const [navbarOpen, setNavbarOpen] = useState(false);


  const handleToggle: NavbarOpenHandler = () => {
    setNavbarOpen(!navbarOpen);
  };

  return (
    <div>

      <Navbar open={navbarOpen} setNavbarOpen={handleToggle} variant={NavbarType.MINIPAGE} />

      <HeadingContainer>
        <Heading>Support CSESoc</Heading>
      </HeadingContainer>
      <ImgContainer>
        <Image src={Otter} />
      </ImgContainer>

      <Grid>
        
        <OurSponsorsCol>
          <TextContainer>
            <H3>Our Sponsors</H3>
            <Text>Check out our very cool sponsors</Text>
            <ButtonContainer>
              <Link href="/Sponsors">
                <Button>View our sponsors</Button>
              </Link>
            </ButtonContainer>
          </TextContainer>
        </OurSponsorsCol>

        <SmallGrid>
          <SponsorCol style={{ backgroundColor: 'rgba(0, 71, 255, 0.33)' }}>
            <Image
              src={`/assets/sponsors_white/atl.svg`}
              width={100}
              height={30}
              objectFit="contain"
            />
          </SponsorCol>
          <SponsorCol style={{ backgroundColor: 'rgba(82, 130, 255, 0.47)' }}>
            <Image
                src={`/assets/sponsors_white/imc.svg`}
                width={60}
                height={50}
                objectFit="contain"
              />
          </SponsorCol>
        
          <SponsorCol style={{ backgroundColor: 'rgba(48, 93, 255, 0.2)' }}>
          <Image
              src={`/assets/sponsors_white/deloitte.svg`}
              width={50}
              height={40}
              objectFit="contain"
            />
          </SponsorCol>
          <SponsorCol style={{ backgroundColor: 'rgba(122, 137, 236, 0.27)' }}>
          <Image
              src={`/assets/sponsors_white/js.svg`}
              width={50}
              height={55}
              objectFit="contain"
            />
          </SponsorCol>
        </SmallGrid>
      
      </Grid>

      <Footer />
    </div>
  );

}