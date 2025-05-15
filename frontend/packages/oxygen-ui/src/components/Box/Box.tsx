import React from 'react';
import { Box as MUIBox, BoxProps } from '@mui/material';

const Box = (props: BoxProps): React.JSX.Element => {
  return <MUIBox {...props} />;
};

export default Box;
