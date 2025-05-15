import { createTheme } from '@mui/material/styles';
import { lightPalette, darkPalette } from './palettes';
import typography from './typography';
import components from './components';

const getTheme = (mode: 'light' | 'dark') => {
  const baseTheme = createTheme({
    palette: mode === 'light' ? lightPalette : darkPalette,
    typography,
  });

  // Now create the final theme with components, passing baseTheme so it can use palette colors
  return createTheme(baseTheme, {
    components: components(baseTheme),
  });
};

export default getTheme;
