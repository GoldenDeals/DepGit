module.exports = {
  root: true,
  extends: [
    'eslint:recommended'
  ],
  plugins: [
    'svelte',
    '@typescript-eslint'
  ],
  parser: '@typescript-eslint/parser',
  overrides: [
    {
      files: ['*.svelte'],
      parser: 'svelte-eslint-parser',
      parserOptions: {
        parser: '@typescript-eslint/parser'
      },
      rules: {
        'no-unused-vars': 'off'
      }
    },
    {
      files: ['*.ts'],
      parser: '@typescript-eslint/parser'
    }
  ],
  parserOptions: {
    sourceType: 'module',
    ecmaVersion: 2020,
    extraFileExtensions: ['.svelte']
  },
  env: {
    browser: true,
    es2017: true,
    node: true
  },
  rules: {
    // Basic rules to accept TypeScript and Svelte syntax
    'no-inner-declarations': 'off'
  }
}; 