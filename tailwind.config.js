/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/templates/pages/*.{templ,js}", 
    "./internal/templates/layouts/*.{templ,js}", 
    "./internal/templates/components/*.{templ,js}", 
    "./tmpl/**/*.{html,tmpl}"],
  theme: {
    container: {
      center: true
    },
    extend: {},
  },
  plugins: [],
}

