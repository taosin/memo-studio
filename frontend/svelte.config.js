export default {
  compilerOptions: {
    css: "injected",
    warningFilter: (warning) => {
      // Ignore a11y warnings during build
      if (warning.code && warning.code.startsWith('a11y_')) return false;
      return true;
    }
  }
};
