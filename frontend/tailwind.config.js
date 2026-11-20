/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
        display: ['"Plus Jakarta Sans"', 'Inter', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'monospace'],
      },
      colors: {
        canvas: {
          light: '#FAFAFA',
          dark: '#09090B',
        },
        surface: {
          light: '#FFFFFF',
          dark: '#18181B',
          subtleLight: '#F4F4F5',
          subtleDark: '#27272A',
          borderLight: '#E4E4E7',
          borderDark: '#27272A',
        },
        zinc: {
          850: '#1f1f23',
          950: '#0c0c0e',
        },
        brand: {
          50: '#EEF2FF',
          100: '#E0E7FF',
          200: '#C7D2FE',
          300: '#A5B4FC',
          400: '#818CF8',
          500: '#6366F1',
          600: '#4F46E5', // Primary Linear Accent
          700: '#4338CA',
          800: '#3730A3',
          900: '#312E81',
        },
      },
      boxShadow: {
        'subtle': '0 1px 2px 0 rgba(0, 0, 0, 0.04)',
        'card': '0 1px 2px 0 rgba(0, 0, 0, 0.03)',
        'card-hover': '0 4px 14px 0 rgba(0, 0, 0, 0.06)',
        'drawer': '-4px 0 28px 0 rgba(0, 0, 0, 0.12)',
        'palette': '0 16px 40px 0 rgba(0, 0, 0, 0.24)',
      }
    },
  },
  plugins: [],
}
