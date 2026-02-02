const { defineConfig } = require('@playwright/test');

// JUnit reporter enabled when output path is set (e.g. by Testkube workflow); includes execution time per test
const reporters = process.env.PLAYWRIGHT_JUNIT_OUTPUT_NAME
  ? ['list', ['junit', { outputFile: process.env.PLAYWRIGHT_JUNIT_OUTPUT_NAME }]]
  : 'list';

module.exports = defineConfig({
  testDir: './tests',
  timeout: 60 * 1000,
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: 0,
  reporter: reporters,
  use: {
    baseURL: 'https://testkube-test-page-lipsum.pages.dev/',
    video: 'off',
    screenshot: 'off',
    trace: 'off',
  },
});
