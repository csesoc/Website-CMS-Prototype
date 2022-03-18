import React from "react";
import { ComponentStory, ComponentMeta } from '@storybook/react';

import { ReactComponent as Italic } from '../../assets/italics-button.svg';
import ItalicButton from './ItalicButton';

// const stories = generateStories("Buttons");

// stories.add("Buttons", () => {
//   return (
//     <div>
//       <h1>this is a button</h1>
//       <Button/>
//     </div>
//   )
// })

export default {
  title: 'CSE-UIKIT/Italic-Button',
  component: ItalicButton,
  argTypes: {
    backgroundColor: { control: 'color' },
  },
} as ComponentMeta<typeof ItalicButton>;

const Template: ComponentStory<typeof ItalicButton> = (args) =>
(
  <div
    style={{
      margin: "30px"
    }}
  >
    Italic Button
    <ItalicButton {...args}><Italic height={parseInt(args.size)*0.55} width={parseInt(args.size)*0.55}/></ItalicButton>
  </div>
)

export const Primary = Template.bind({});
Primary.args = {
  background: "#E2E1E7",
  size: "45px"
}