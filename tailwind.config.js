/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/ui/**/*.html',
    './src/ui/**/*.tmpl',
  ],
  theme: {
    extend: {
      colors: {
        primary: '#2980b9',
        secondary: '#ffa500',
        background: '#ffffff',
        surface: '#eeeeee',
      },
    },
    fontFamily: {
      sans: ['Josefin Sans', 'sans-serif'],
      pacifico: ['Pacifico', 'cursive'],
    },
  },
  plugins: [],
}

