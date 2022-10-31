import styled from "styled-components";
import React, { useState, FC, useRef, useEffect } from "react";

import Client from "./websocketClient";

import HeadingBlock from "./components/HeadingBlock";
import EditorBlock from "./components/EditorBlock";
import MediaBlock from "./components/MediaBlock";
import { BlockData, UpdateHandler } from "./types";
import CreateContentBlock from "src/cse-ui-kit/CreateContentBlock_button";
import CreateHeadingBlock from "src/cse-ui-kit/CreateHeadingBlock_button";
import CreateMediaBlock from "src/cse-ui-kit/CreateMediaBlock_button";
import SyncDocument from "src/cse-ui-kit/SyncDocument_button";
import PublishDocument from "src/cse-ui-kit/PublishDocument_button";
import EditorHeader from "src/deprecated/components/Editor/EditorHeader";
import { addContentBlock } from "./state/actions";
import { useParams } from "react-router-dom";
import { defaultContent, headingContent, mediaContent } from "./state/helpers";

// Redux
import { useDispatch } from "react-redux";

const Container = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const InsertContentWrapper = styled.div`
  display: flex;
`;

const EditorPage: FC = () => {
  const dispatch = useDispatch();
  const { id } = useParams();
  const wsClient = useRef<Client | null>(null);

  const [blocks, setBlocks] = useState<BlockData[]>([]);
  const [focusedId, setFocusedId] = useState<number>(0);

  const updateValues: UpdateHandler = (idx, updatedBlock) => {
    if (JSON.stringify(blocks[idx]) === JSON.stringify(updateValues)) return;
    setBlocks((prev) =>
      prev.map((block, i) => (i === idx ? updatedBlock : block))
    );
  };

  useEffect(() => {
    function cleanup() {
      wsClient.current?.close();
    }

    wsClient.current = new Client(
      id as string,
      (data) => {
        console.log(id, JSON.stringify(data));
        setBlocks(data as BlockData[]);
      },
      (reason) => {
        console.log(reason);
      }
    );
    window.addEventListener("beforeunload", cleanup);
    return () => {
      console.log("Editor component destroyed");
      wsClient.current?.close();
      window.removeEventListener("beforeunload", cleanup);
    };
  }, []);

  return (
    <div style={{ height: "100%" }}>
      <EditorHeader />
      <Container>
        {blocks.map((block, idx) => {
          console.log(block[0].type);
          switch (block[0].type) {
            case "heading":
              return (
                <HeadingBlock
                  id={idx}
                  key={idx}
                  update={updateValues}
                  showToolBar={focusedId === idx}
                  onEditorClick={() => setFocusedId(idx)}
                />
              )
            case "image":
              return (
                <MediaBlock
                  id={idx}
                  key={idx}
                  update={updateValues}
                  showToolBar={focusedId === idx}
                  onMediaClick={() => setFocusedId(idx)}
                />
              )
            default:
              return (
                <EditorBlock
                  id={idx}
                  key={idx}
                  initialValue={block}
                  update={updateValues}
                  showToolBar={focusedId === idx}
                  onEditorClick={() => setFocusedId(idx)}
                />

              )
          }
        })}

        <InsertContentWrapper>
          <CreateHeadingBlock
            onClick={() => {
              setBlocks((prev) => [
                ...prev,
                [{ type: "heading", children: [{ text: "" }] }],
              ]);

              // create the initial state of the content block to Redux
              dispatch(
                addContentBlock({
                  id: blocks.length,
                  data: headingContent,
                })
              );
              setFocusedId(blocks.length);
            }}
          />
          <CreateContentBlock
            onClick={() => {
              setBlocks((prev) => [
                ...prev,
                [{ type: "paragraph", children: [{ text: "" }] }],
              ]);

              // create the initial state of the content block to Redux
              dispatch(
                addContentBlock({
                  id: blocks.length,
                  data: defaultContent,
                })
              );
              setFocusedId(blocks.length);
            }}
          />
          <CreateMediaBlock
            onClick={() => {
              setBlocks((prev) => [
                ...prev,
                [{ type: "image", url: "" }],
              ]);

              // create the initial state of the content block to Redux
              dispatch(
                addContentBlock({
                  id: blocks.length,
                  data: mediaContent,
                })
              );
              setFocusedId(blocks.length);
            }}
          />
          <SyncDocument
            onClick={() => {
              if (wsClient.current?.socket.readyState === WebSocket.OPEN) {
                console.log(JSON.stringify(blocks));
                wsClient.current?.pushDocumentData(JSON.stringify(blocks));
              }
            }}
          />
          <PublishDocument
            onClick={() => {
              fetch("/api/filesystem/publish-document", {
                method: "POST",
                headers: {
                  "Content-Type": "application/x-www-form-urlencoded",
                },
                body: new URLSearchParams({
                  DocumentID: `${id}`,
                }),
              });
            }}
          />
        </InsertContentWrapper>
      </Container>
    </div>
  );
};

export default EditorPage;
