/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/templates/pages/*.{templ,js}", 
    "./internal/templates/layouts/*.{templ,js}", 
    "./internal/templates/components/*.{templ,js}", 
    "./tmpl/**/*.{html,tmpl}"],
  theme: {
    container: {
      center: true,
      padding: '2rem',
    },
    colors: {
      gray: '#c3c3c3',
      white: '#ffffff',
      'dark-charkol': '#222831',
      'light-gray': '#393E46',
      primaryCyan: '#00ADB5',
      'light-white': '#EEEEEE'
    },
    extend: {
      fontFamily: {
        nunito: ['"Nunito Regular"', 'sans-serif'],
        'nunito-bold': ['"Nunito Bold"', 'sans-serif'],
        'nunito-light': ['"Nunito Light"', 'sans-serif'],
      },
    },
  },
  plugins: [],
}

