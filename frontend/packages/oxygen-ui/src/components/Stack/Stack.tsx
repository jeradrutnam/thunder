import React from 'react';
import { Stack as MUIStack, StackProps } from '@mui/material';

const Stack = (props: StackProps): React.JSX.Element => {
  return <MUIStack {...props} />;
};

export { StackProps };

export default Stack;
