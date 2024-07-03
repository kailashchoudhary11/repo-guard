/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./templates/**/*.templ'],
  theme: {
    extend: {
      fontFamily: {
        maven: ['"Maven Pro"', 'sans-serif'],
        rubik: ['"Rubik', 'sans-serif'],
      }
    },
  },
  plugins: [],
}

