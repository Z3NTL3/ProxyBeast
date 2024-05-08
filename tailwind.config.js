/** @type {import('tailwindcss').Config} */
module.exports = {
  content: {
    files: ["./frontend/**/*.{html,js,css}", "./frontend/*.{html,js,css}", "*.{html,js,css}"]
  },
  theme: {
    extend: {},
  },
  plugins: [
    require("daisyui"),
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
  ],
  daisyui: {
    themes: [
      {
        proxybeast: {
          "base-100": "#212024",
          "base-200": "#2C2D30", // 2C2D30
          "primary": "#9379FF",
          "secondary": "#5EC9FF"
        }
      }
    ]
  }
}

