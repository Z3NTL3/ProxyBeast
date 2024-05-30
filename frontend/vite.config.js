// vite.config.js
export default {
  build: {
    // Make sure to set the base if you are serving from a subdirectory
    base: '/',
    rollupOptions: {
      external: ['/js/sweetalert2.all.min.js']
    }
  },
  // Other configurations
};