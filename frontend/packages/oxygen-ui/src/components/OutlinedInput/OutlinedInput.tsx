import React from 'react';
import { OutlinedInput as MUIOutlinedInput, OutlinedInputProps } from '@mui/material';

const OutlinedInput = (props: OutlinedInputProps): React.JSX.Element => {
  return <MUIOutlinedInput {...props} />;
};

export { OutlinedInputProps };

export default OutlinedInput;
