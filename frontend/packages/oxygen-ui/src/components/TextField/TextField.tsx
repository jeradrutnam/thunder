import React from 'react';
import { TextField as MUITextField, TextFieldProps } from '@mui/material';

const TextField = (props: TextFieldProps): React.JSX.Element => {
  return (
    <MUITextField
      variant="outlined"
      fullWidth
      slotProps={{
        inputLabel: {
          shrink: true,
        },
      }}
      {...props} />
  );
};

export default TextField;
