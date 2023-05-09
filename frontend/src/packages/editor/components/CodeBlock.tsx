import React, { FC, useMemo, useCallback } from 'react';
import { 
  Slate,
  Editable,
  withReact,
  RenderLeafProps,
  RenderElementProps,
  useSlateStatic,
  ReactEditor,
  useSlate
} from 'slate-react';

import './styles/prism.css';

import { CustomElement } from '../types';

import Prism from 'prismjs';
import 'prismjs/components/prism-javascript';
import 'prismjs/components/prism-jsx';
import 'prismjs/components/prism-typescript';
import 'prismjs/components/prism-tsx';
import 'prismjs/components/prism-markdown';
import 'prismjs/components/prism-python';
import 'prismjs/components/prism-php';
import 'prismjs/components/prism-sql';
import 'prismjs/components/prism-java';

import styled from 'styled-components';
import { 
  Range,
  createEditor,
  Editor,
  Element, 
  NodeEntry,
  Transforms,
  Node
} from 'slate';

import { CMSBlockProps } from '../types';

import CodeContentBlock from "../../../cse-ui-kit/codeblock/codecontentblock-wrapper";
import { handleKey } from "./buttons/buttonHelpers";
import EditorCodeButton from "./buttons/EditorCodeButton";
import { normalizeTokens } from './util/normalize-tokens';


const ToolbarContainer = styled.div`
  display: flex;
  flex-direction: row;
  width: 100%;
  max-width: 600px;
  margin: 5px;
  justify-content: flex-end;
`;

const LanguageSelectWrapper = styled.select`
  content-editable: false;
  right: 5px;
  top: 5px;
  z-index: 1;
`;

const LanguageSelect = (props: JSX.IntrinsicElements['select']) => {
  return (
    <LanguageSelectWrapper value={props.value} onChange={props.onChange}>
      <option value="css">CSS</option>
      <option value="html">HTML</option>
      <option value="java">Java</option>
      <option value="javascript">JavaScript</option>
      <option value="jsx">JSX</option>
      <option value="markdown">Markdown</option>
      <option value="php">PHP</option>
      <option value="python">Python</option>
      <option value="sql">SQL</option>
      <option value="tsx">TSX</option>
      <option value="typescript">TypeScript</option>
    </LanguageSelectWrapper>
  )
}


const CodeBlock: FC<CMSBlockProps> = ({
  id,
  update,
  initialValue,
  showToolBar,
  onEditorClick,
}) => {
  const editor = useMemo(() => withReact(createEditor()), []);

  const renderLeaf: (props: RenderLeafProps) => JSX.Element = useCallback(
    ({ attributes, children, leaf }) => {
      const props = {
        ...attributes,
      };

      const {text, ...rest} = leaf;

      return (
        <span {...props} className={Object.keys(rest).join(' ')}>{
          children}
        </span>
      )
    },
    []
  );

    const language = editor.children.length > 0 
      ? (editor.children[0] as Element).language 
      : "css";

    console.log(editor.children);

  // const blocks = Editor.nodes<Element>(editor, {
  //   at: [],
  //   mode: 'highest',
  //   //  Find all code block nodes
  //   match: n => Element.isElement(n) && n.type === 'code-block',
  // })

  // const language = "python";

  const decorate = useDecorate(editor);

  const setLanguage = (newLanguage : string) => {
    // Select entire editor and change all lines to have language `newLanguage`.
    Transforms.select(editor, {
      anchor: Editor.start(editor, []),
      focus: Editor.end(editor, []),
    })
    Transforms.setNodes(editor, { language: newLanguage });
  };

  return (
    <Slate
      editor={editor}
      
      value={initialValue}
      onChange={(value) => update(id, editor.children, editor.operations)}
    >
      {showToolBar && 
        <ToolbarContainer>
          <LanguageSelect 
            value={language}
            onChange={event => setLanguage(event.target.value)}
          />
        </ToolbarContainer>
      }
      <SetNodeToDecorations/>
      <CodeContentBlock focused={showToolBar}>
        <Editable
          decorate={decorate}
          renderElement={ElementWrapper}
          renderLeaf={renderLeaf}
          onClick={() => onEditorClick()}
          autoFocus
          style={{ width: "100%", height: "100%" }}
        />
        <style>{prismThemeCss}</style>
      </CodeContentBlock>
    </Slate>
  );
};

const ElementWrapper = (props: RenderElementProps) => {
  const { attributes, children, element } = props;
  return (
    <div {...attributes} style={{ position: 'relative' }}>
        {children}
    </div>
  );
}

const useDecorate = (editor: Editor) => {
  return useCallback(
    ([node, path]) => {
      const ranges = editor.nodeToDecorations?.get(node) || [];
      return ranges;
    },
    [editor.nodeToDecorations]
  )
}

const getChildNodeToDecorations = ([block, blockPath]: NodeEntry<
  CustomElement
>) => {
  const nodeToDecorations = new Map<Element, Range[]>()

  const text = block.children.map((line) => Node.string(line)).join('\n');
  const language = block.language ?? "javascript";
  const tokens = Prism.tokenize(text, Prism.languages[language])
  const normalizedTokens = normalizeTokens(tokens) // make tokens flat and grouped by line
  const blockChildren = block.children as unknown as Element[];

  for (let index = 0; index < normalizedTokens.length; index++) {
    const tokens = normalizedTokens[index];
    const element = blockChildren[index];

    if (!nodeToDecorations.has(element)) {
      nodeToDecorations.set(element, []);
    }

    let start = 0;
    for (const token of tokens) {
      const length = token.content.length;
      if (!length) {
        continue;
      }

      const end = start + length;

      const path = [...blockPath, index, 0];
      const range = {
        anchor: { path, offset: start },
        focus: { path, offset: end },
        token: true,
        ...Object.fromEntries(token.types.map(type => [type, true])),
      };

      nodeToDecorations.get(element)?.push(range);

      start = end;
    }
  }

  return nodeToDecorations
}


// precalculate editor.nodeToDecorations map to use it inside decorate function then
const SetNodeToDecorations = () => {
  const editor = useSlate();

  const blockEntries = Array.from(
    Editor.nodes<Element>(editor, {
      at: [],
      mode: 'highest',
      //  Find all code block nodes
      match: n => Element.isElement(n) && n.type === 'code-block',
    })
  );
  
  const nodeToDecorations = mergeMaps(
    ...blockEntries.map(getChildNodeToDecorations)
  );
  
  editor.nodeToDecorations = nodeToDecorations;
  
  return null;
}
    
const mergeMaps = <K, V>(...maps: Map<K, V>[]) => {
  const map = new Map<K, V>();

  for (const m of maps) {
    for (const item of m) {
      map.set(...item);
    }
  }

  return map;

}

const prismThemeCss = `
/**
 * prism.js default theme for JavaScript, CSS and HTML
 * Based on dabblet (http://dabblet.com)
 * @author Lea Verou
 */
code[class*="language-"],
pre[class*="language-"] {
    color: black;
    background: none;
    text-shadow: 0 1px white;
    font-family: Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
    font-size: 1em;
    text-align: left;
    white-space: pre;
    word-spacing: normal;
    word-break: normal;
    word-wrap: normal;
    line-height: 1.5;
    -moz-tab-size: 4;
    -o-tab-size: 4;
    tab-size: 4;
    -webkit-hyphens: none;
    -moz-hyphens: none;
    -ms-hyphens: none;
    hyphens: none;
}
pre[class*="language-"]::-moz-selection, pre[class*="language-"] ::-moz-selection,
code[class*="language-"]::-moz-selection, code[class*="language-"] ::-moz-selection {
    text-shadow: none;
    background: #b3d4fc;
}
pre[class*="language-"]::selection, pre[class*="language-"] ::selection,
code[class*="language-"]::selection, code[class*="language-"] ::selection {
    text-shadow: none;
    background: #b3d4fc;
}
@media print {
    code[class*="language-"],
    pre[class*="language-"] {
        text-shadow: none;
    }
}
/* Code blocks */
pre[class*="language-"] {
    padding: 1em;
    margin: .5em 0;
    overflow: auto;
}
:not(pre) > code[class*="language-"],
pre[class*="language-"] {
    background: #f5f2f0;
}
/* Inline code */
:not(pre) > code[class*="language-"] {
    padding: .1em;
    border-radius: .3em;
    white-space: normal;
}
.token.comment,
.token.prolog,
.token.doctype,
.token.cdata {
    color: slategray;
}
.token.punctuation {
    color: #999;
}
.token.namespace {
    opacity: .7;
}
.token.property,
.token.tag,
.token.boolean,
.token.number,
.token.constant,
.token.symbol,
.token.deleted {
    color: #905;
}
.token.selector,
.token.attr-name,
.token.string,
.token.char,
.token.builtin,
.token.inserted {
    color: #690;
}
.token.operator,
.token.entity,
.token.url,
.language-css .token.string,
.style .token.string {
    color: #9a6e3a;
    /* This background color was intended by the author of this theme. */
    background: hsla(0, 0%, 100%, .5);
}
.token.atrule,
.token.attr-value,
.token.keyword {
    color: #07a;
}
.token.function,
.token.class-name {
    color: #DD4A68;
}
.token.regex,
.token.important,
.token.variable {
    color: #e90;
}
.token.important,
.token.bold {
    font-weight: bold;
}
.token.italic {
    font-style: italic;
}
.token.entity {
    cursor: help;
}
`;


export default CodeBlock;
