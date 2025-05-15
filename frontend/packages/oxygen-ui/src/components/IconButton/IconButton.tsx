import React from 'react';
import { IconButton as MUIIconButton, IconButtonProps } from '@mui/material';

const IconButton = (props: IconButtonProps): React.JSX.Element => {
  return <MUIIconButton {...props} />;
};

export { IconButtonProps };

export default IconButton;
