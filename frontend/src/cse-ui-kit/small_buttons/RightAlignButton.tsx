import React, { MouseEventHandler } from 'react';
import { StyledButton, buttonProps, scaleRate } from './small_buttons-Styled';
import { ReactComponent as RightAlign } from 'src/cse-ui-kit/assets/leftrightalign-button.svg';

type Props = {
  onClick?: MouseEventHandler<HTMLDivElement>;
  onMouseDown?: MouseEventHandler<HTMLDivElement>;
} & buttonProps;

export default function RightAlignButton({
  onClick,
  onMouseDown,
  ...styleProps
}: Props) {
  return (
    <StyledButton onClick={onClick} onMouseDown={onMouseDown} {...styleProps}>
      <RightAlign
        height={styleProps.size * scaleRate}
        width={styleProps.size * scaleRate}
        transform={"scale(-1, 1)"}
      />
    </StyledButton>
  );
}
