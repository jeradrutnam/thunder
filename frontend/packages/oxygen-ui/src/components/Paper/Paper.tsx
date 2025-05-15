import React from 'react';
import { Paper as MUIPaper, PaperProps } from '@mui/material';

const Paper = (props: PaperProps): React.JSX.Element => {
  return <MUIPaper {...props} />;
};

export default Paper;
