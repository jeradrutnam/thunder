'use client';

import React, {use} from 'react';
import { Button } from '@mui/material';
import { useThemeMode } from '../../contexts/ThemeContext';

const ColorSchemeToggle = () => {
  const { resolvedMode, source, setSource } = useThemeMode();

  const next = source === 'light' ? 'dark' : source === 'dark' ? 'system' : 'light';

  return (
    <Button color='secondary' onClick={() => setSource(next)}>
      {resolvedMode === 'dark' ? '🌙' : '☀️'}{' '}{source}
    </Button>
  );
};

export default ColorSchemeToggle;
