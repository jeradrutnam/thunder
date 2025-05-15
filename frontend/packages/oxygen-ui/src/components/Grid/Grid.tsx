import React from 'react';
import { Grid as MUIGrid, GridProps } from '@mui/material';

const Grid = (props: GridProps): React.JSX.Element => {
  return <MUIGrid {...props} />;
};

export default Grid;
