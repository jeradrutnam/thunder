import React from 'react';
import { Divider as MUIDivider, DividerProps } from '@mui/material';

const Divider = (props: DividerProps): React.JSX.Element => {
  return <MUIDivider {...props} />;
};

export { DividerProps };

export default Divider;
