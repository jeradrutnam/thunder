import React from 'react';
import { InputLabel as MUIInputLabel, InputLabelProps } from '@mui/material';

const InputLabel = (props: InputLabelProps): React.JSX.Element => {
  return <MUIInputLabel {...props} />;
};

export { InputLabelProps };

export default InputLabel;