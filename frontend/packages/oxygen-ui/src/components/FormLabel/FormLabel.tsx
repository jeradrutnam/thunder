import React from 'react';
import { FormLabel as MUIFormLabel, FormLabelProps } from '@mui/material';

const FormLabel = (props: FormLabelProps): React.JSX.Element => {
  return <MUIFormLabel {...props} />;
};

export { FormLabelProps };

export default FormLabel;
