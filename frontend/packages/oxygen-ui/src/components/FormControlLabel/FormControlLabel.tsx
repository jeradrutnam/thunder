import React from 'react';
import { FormControlLabel as MUIFormControlLabel, FormControlLabelProps } from '@mui/material';

const FormControlLabel = (props: FormControlLabelProps): React.JSX.Element => {
  return <MUIFormControlLabel {...props} />;
};

export { FormControlLabelProps };

export default FormControlLabel;