import React from 'react';
import { Checkbox as MUICheckbox, CheckboxProps } from '@mui/material';

const Checkbox = (props: CheckboxProps): React.JSX.Element => {
  return <MUICheckbox {...props} />;
};

export { CheckboxProps };

export default Checkbox;
