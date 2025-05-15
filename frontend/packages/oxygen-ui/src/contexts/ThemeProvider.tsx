'use client';

import React, { useEffect, useState, useMemo } from 'react';
import { ThemeProvider as MUIThemeProvider, CssBaseline } from '@mui/material';
import getTheme from '../theme/index';
import { ThemeModeContext, ColorScheme } from './ThemeContext';

const STORAGE_KEY = 'color-scheme';

interface Props {
  children: React.ReactNode;
}

export const ThemeProvider: React.FC<Props> = ({ children }) => {
  const [source, setSource] = useState<ColorScheme>(() => {
    if (typeof window === 'undefined') return 'system';
    return (localStorage.getItem(STORAGE_KEY) as ColorScheme) || 'system';
  });

  const [resolvedMode, setResolvedMode] = useState<'light' | 'dark'>(() => {
    if (typeof window === 'undefined') return 'light';
    const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    return systemPrefersDark ? 'dark' : 'light';
  });

  // Keep resolvedMode synced with system if needed
  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');

    const updateMode = () => {
      if (source === 'system') {
        setResolvedMode(mediaQuery.matches ? 'dark' : 'light');
      } else {
        setResolvedMode(source);
      }
    };

    updateMode(); // Initial
    mediaQuery.addEventListener('change', updateMode);
    return () => mediaQuery.removeEventListener('change', updateMode);
  }, [source]);

  // Persist user preference
  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, source);
  }, [source]);

  const theme = useMemo(() => getTheme(resolvedMode), [resolvedMode]);

  return (
    <ThemeModeContext.Provider value={{ resolvedMode, source, setSource }}>
      <MUIThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </MUIThemeProvider>
    </ThemeModeContext.Provider>
  );
};
