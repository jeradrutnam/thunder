const fs = require('fs');
const path = require('path');

const rootDir = path.resolve(__dirname, '..');
const componentsDir = path.join(rootDir, 'src', 'components');
const packageJsonPath = path.join(rootDir, 'package.json');

const packageJson = require(packageJsonPath);

const exportsField = {
  '.': {
    import: './dist/index.js',
    require: './dist/index.js',
  },
};

fs.readdirSync(componentsDir).forEach((file) => {
  const ext = path.extname(file);
  const name = path.basename(file, ext);
  if (ext === '.ts' || ext === '.tsx') {
    exportsField[`./${name}`] = {
      import: `./dist/components/${name}.js`,
      require: `./dist/components/${name}.js`,
    };
  }
});

packageJson.exports = exportsField;

fs.writeFileSync(packageJsonPath, JSON.stringify(packageJson, null, 2));
