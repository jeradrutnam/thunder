import { Components, Theme } from '@mui/material/styles';

const components = (theme: Theme): Components => ({
  MuiLink: {
    styleOverrides: {
      root: {
        color: theme.palette.primary.main,
        textDecoration: 'none',
        '&:hover': {
          color: theme.palette.primary.dark,
          textDecoration: 'underline',
        },
      },
    },
  },
  MuiButton: {
    styleOverrides: {
      root: {
        borderRadius: 20,
        textTransform: 'none',
      },
      containedPrimary: {
        background: 'linear-gradient(77.74deg, #eb4f63 11.16%, #fa7b3f 99.55%)',
        color: theme.palette.text?.primary,
        '&:hover': {
          background: 'linear-gradient(77.74deg, #e94459 11.16%, #f56f30 99.55%)',
        },
      },
      containedSecondary: {
        backgroundColor: theme.palette.secondary?.main,
        color: theme.palette.text?.secondary,
        '&:hover': {
          backgroundColor: theme.palette.secondary?.dark || theme.palette.secondary?.main,
        },
      },
    },
  },
});

export default components;
