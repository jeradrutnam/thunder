import React from 'react';
import { Typography as MUITypography, TypographyProps } from '@mui/material';

const Typography = (props: TypographyProps): React.JSX.Element => {
  return <MUITypography {...props} />;
};

export default Typography;
