/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      keyframes: {
        fadeInCpu: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        'cpu-glow': {
          '0%, 100%': { boxShadow: '0 0 20px 0px rgba(6,182,212,0.3)' },
          '50%': { boxShadow: '0 0 60px 10px rgba(6,182,212,0.6)' },
        },
        'cpu-scan': {
          '0%': { transform: 'translateY(-100%)' },
          '100%': { transform: 'translateY(100%)' },
        },
        'path-draw': {
          '0%': { 'stroke-dashoffset': '1' },
          '100%': { 'stroke-dashoffset': '0' },
        },
      },
      animation: {
        'fade-in-cpu': 'fadeInCpu 2s ease forwards',
        'fade-in-short': 'fadeInCpu 0.4s ease',
        'cpu-glow': 'cpu-glow 3s ease-in-out infinite',
        'cpu-scan': 'cpu-scan 2s linear infinite',
        'path-draw': 'path-draw 2s ease-in-out infinite alternate',
        'path-draw-slow': 'path-draw 2.5s ease-in-out infinite alternate',
        'path-draw-fast': 'path-draw 1.8s ease-in-out infinite alternate',
      },
    },
  },
  plugins: [],
}
