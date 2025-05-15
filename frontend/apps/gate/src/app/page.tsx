"use client";

import Box from '@oxygen-ui/react/src/components/Box/Box';
import Button from '@oxygen-ui/react/src/components/Button/Button';
import Checkbox from '@oxygen-ui/react/src/components/Checkbox/Checkbox';
import ColorSchemeToggle from '@oxygen-ui/react/src/components/ColorSchemeToggle/ColorSchemeToggle';
import Divider from '@oxygen-ui/react/src/components/Divider/Divider';
import FormControl from '@oxygen-ui/react/src/components/FormControl/FormControl';
import FormControlLabel from '@oxygen-ui/react/src/components/FormControlLabel/FormControlLabel';
import FormLabel from '@oxygen-ui/react/src/components/FormLabel/FormLabel';
import Grid from '@oxygen-ui/react/src/components/Grid/Grid';
import IconButton, { IconButtonProps } from '@oxygen-ui/react/src/components/IconButton/IconButton';
import OutlinedInput from '@oxygen-ui/react/src/components/OutlinedInput/OutlinedInput';
import Paper from '@oxygen-ui/react/src/components/Paper/Paper';
import Link from '@oxygen-ui/react/src/components/Link/Link';
import Input from '@oxygen-ui/react/src/components/Input/Input';
import InputLabel from '@oxygen-ui/react/src/components/InputLabel/InputLabel';
// import { useColorScheme } from '@oxygen-ui/react/src/hooks/useColorScheme';
import Stack from '@oxygen-ui/react/src/components/Stack/Stack';
import TextField from '@oxygen-ui/react/src/components/TextField/TextField';
import { useThemeMode } from '@oxygen-ui/react/src/contexts/ThemeContext';
import Typography from '@oxygen-ui/react/src/components/Typography/Typography';
import React, { useState} from 'react';
import Image from 'next/image';

const GoogleIcon = (props: React.SVGProps<SVGSVGElement>) => (
  <svg
    focusable="false"
    aria-hidden="true"
    viewBox="0 0 24 24"
    height={24}
    width={24}
    {...props}
  >
    <g transform="matrix(1, 0, 0, 1, 27.009001, -39.238998)">
      <path fill="#4285F4" d="M -3.264 51.509 C -3.264 50.719 -3.334 49.969 -3.454 49.239 L -14.754 49.239 L -14.754 53.749 L -8.284 53.749 C -8.574 55.229 -9.424 56.479 -10.684 57.329 L -10.684 60.329 L -6.824 60.329 C -4.564 58.239 -3.264 55.159 -3.264 51.509 Z"></path>
      <path fill="#34A853" d="M -14.754 63.239 C -11.514 63.239 -8.804 62.159 -6.824 60.329 L -10.684 57.329 C -11.764 58.049 -13.134 58.489 -14.754 58.489 C -17.884 58.489 -20.534 56.379 -21.484 53.529 L -25.464 53.529 L -25.464 56.619 C -23.494 60.539 -19.444 63.239 -14.754 63.239 Z"></path>
      <path fill="#FBBC05" d="M -21.484 53.529 C -21.734 52.809 -21.864 52.039 -21.864 51.239 C -21.864 50.439 -21.724 49.669 -21.484 48.949 L -21.484 45.859 L -25.464 45.859 C -26.284 47.479 -26.754 49.299 -26.754 51.239 C -26.754 53.179 -26.284 54.999 -25.464 56.619 L -21.484 53.529 Z"></path>
      <path fill="#EA4335" d="M -14.754 43.989 C -12.984 43.989 -11.404 44.599 -10.154 45.789 L -6.734 42.369 C -8.804 40.429 -11.514 39.239 -14.754 39.239 C -19.444 39.239 -23.494 41.939 -25.464 45.859 L -21.484 48.949 C -20.534 46.099 -17.884 43.989 -14.754 43.989 Z"></path>
    </g>
  </svg>
);

const Logo = () => {
  const { resolvedMode } = useThemeMode();

  const src = resolvedMode === 'dark' ? '/images/logo-dark.svg' : '/images/logo-light.svg';

  return <Image src={src} alt="Logo" width={200} height={50} />;
};

export default function LoginPage() {
  return (
    <Box sx={{ height: '100vh', display: 'flex' }}>
      <Grid container sx={{ flex: 1 }}>
        <Grid
          size={{ xs: 12, md: 7 }}
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: 4,
          }}
        >
          <Box>
            <Box sx={{ mb: 3 }}>
              <Logo />
            </Box>
            <Typography variant="h4" gutterBottom>
              Login to Asgardeo - Thunder!
            </Typography>
            <Typography variant="body1">
              Asgardeo empowers developers to implement login experiances in minutes.
            </Typography>
            {/* <Box sx={{ mt: 3 }}>
              <ColorSchemeToggle />
            </Box> */}
          </Box>
        </Grid>

        <Grid
          size={{ xs: 12, md: 5 }}
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: 4,
          }}
        >
          <Paper elevation={3} sx={{ padding: 4, width: '100%', maxWidth: 400 }}>
            <Typography variant="h5" gutterBottom>
              Login to Account
            </Typography>
            <Typography>
              New to company?{' '}
              <Link href="">Sign up!</Link>
            </Typography>
            <Button
                fullWidth
                variant='contained'
                startIcon={<GoogleIcon />}
                color='secondary'
                sx={{ mt: 5 }}
              >
                Continue with Google
            </Button>

            <Divider sx={{ my: 3 }}>or</Divider>

            <Box component="form" display="flex" flexDirection="column" gap={2}>
              <Box display="flex" flexDirection="column" gap={0.5}>
                <InputLabel htmlFor="email-input">Email</InputLabel>
                <OutlinedInput id="email-input" placeholder="Enter your email address" size="small" />
              </Box>
              <Box display="flex" flexDirection="column" gap={0.5}>
                <InputLabel htmlFor="password-input">Password</InputLabel>
                <OutlinedInput id="password-input" placeholder="Enter your password" size="small" />
              </Box>
              <Box
                sx={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                }}
              >
                <FormControlLabel
                  control={<Checkbox name="remember-me-checkbox" />}
                  label="Remember me"
                />
                <Link href="#replace-with-a-link">
                  Forgot your password?
                </Link>
              </Box>
              <Button variant="contained" color="primary" fullWidth sx={{ mt: 2 }}>
                Sign In
              </Button>
            </Box>

            <Box component="footer" sx={{ mt: 5 }}>
              <Typography sx={{ textAlign: 'center' }}>
                © Asgardeo {new Date().getFullYear()}
              </Typography>
            </Box>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
}
