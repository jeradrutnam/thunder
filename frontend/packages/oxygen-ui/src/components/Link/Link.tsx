import React from 'react';
import { Link as MUILink, LinkProps } from '@mui/material';

const Link = (props: LinkProps): React.JSX.Element => {
  return <MUILink {...props} />;
};

export { LinkProps };

export default Link;
