import type { NextPage } from "next";
import Head from "next/head";
import Image from 'next/image';
import styled from "styled-components";

const PageContainer = styled.div`
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  padding-left: 2rem;
  padding-right: 2rem;
`;

const ImageContainer = styled.div`
  @import url('https://fonts.googleapis.com/css?family=Raleway:300,400,700&display=swap');
  position: absolute;
  left: 10%;
  top: 25%;
  font-family: "Raleway", sans-serif;
  justify-content: flex-end;
`;

const Text1 = styled.body`
  color: #010033;
  font-size: 36px;
  padding: 10px 0;
  text-shadow: 0px 4px 4px;
  margin-top: 100px;
`;

const Text2 = styled.body`
  color: #3977F8;
  font-size: 36px;
`;

const Text3 = styled.body`
  color: #010033;
  font-style: italic;
  font-size: 36px;
  display:inline;
`;

const Button = styled.button`
  background-color:#FFFFFF;
  color: #3977F8;
  font-size: 22px;
  margin-top: 150px;
  padding: 0.25em 1em;
  border: 1px solid #3977F8;
  border-radius: 3px;
  position: absolute;
  width: 184px;
  height: 44px;
`;

const Home: NextPage = () => {
  return (
    <PageContainer>
      <Head>
        <title>CSESoc</title>
        <meta name="description" content="CSESoc Website Homepage" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <ImageContainer>
          <Image src="/logo.svg" width="362" height="84"/>
          <Text1>
            Empowering
            <Text3> future</Text3>
            <Text2>Technological Leaders</Text2>
          </Text1>
          <Button>Visit on Blog</Button>
        </ImageContainer>
      </main>

      <footer></footer>
    </PageContainer>
  );
};

export default Home;

