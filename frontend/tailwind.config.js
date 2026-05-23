/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        brand: {
          red: '#E11D48',
          ink: '#111827',
          soft: '#FAFAFA'
        }
      },
      boxShadow: {
        soft: '0 16px 40px rgba(17, 24, 39, 0.08)'
      }
    }
  },
  plugins: []
};
