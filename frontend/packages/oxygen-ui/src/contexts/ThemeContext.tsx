'use client';

import React, { createContext, useContext } from 'react';

export type ColorScheme = 'light' | 'dark' | 'system';

export const ThemeModeContext = createContext<{
  resolvedMode: 'light' | 'dark';
  source: ColorScheme;
  setSource: (source: ColorScheme) => void;
}>({
  resolvedMode: 'light',
  source: 'system',
  setSource: () => {},
});

export const useThemeMode = () => {
  const context = useContext(ThemeModeContext);
  if (!context) {
    throw new Error('useThemeMode must be used within ThemeProvider');
  }
  return context;
};
