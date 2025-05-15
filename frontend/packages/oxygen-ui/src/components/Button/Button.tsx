import React from 'react';
import { Button as MUIButton, ButtonProps } from '@mui/material';

const Button = (props: ButtonProps): React.JSX.Element => {
  return <MUIButton disableElevation {...props} />;
};

export default Button;
