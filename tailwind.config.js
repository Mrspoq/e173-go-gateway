/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.{html,tmpl}",
    "./pkg/web/templates/**/*.{html,tmpl}",
    "./templates/partials/**/*.{html,tmpl}"
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
