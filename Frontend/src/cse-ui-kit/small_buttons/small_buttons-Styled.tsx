import styled from "styled-components";

export const scaleRate = {
  smallButtonRate: 0.8,
  underlineRate: 0.6,
}

export type buttonProps = {
  background?: string;
  size: number;
}
export const StyledButton = styled.div<buttonProps>`
  height: ${props => props.size}px;
  width: ${props => props.size}px;
  background: ${props => props.background};
  color: #000000;
  box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.25);
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: ${props => props.size / 10}px;

  &:hover {
    background: #EFEEF3;
    transform: scale(1.04);
  }
  &:active {
    background: #C8D1FA;
    color: #7482CB;
    transform: scale(0.96);
  }

  cursor: pointer;
`;