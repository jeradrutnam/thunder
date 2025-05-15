import React from 'react';
import { FormControl as MUIFormControl, FormControlProps } from '@mui/material';

const FormControl = (props: FormControlProps): React.JSX.Element => {
  return <MUIFormControl {...props} />;
};

export { FormControlProps };

export default FormControl;