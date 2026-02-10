/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        gray: {
          750: '#2d2d44',
          850: '#1f1f35',
        },
      },
    },
  },
  plugins: [],
}
