import React from 'react';
import { Input as MUIInput, InputProps } from '@mui/material';

const Input = (props: InputProps): React.JSX.Element => {
  return <MUIInput {...props} />;
};

export { InputProps };

export default Input;
