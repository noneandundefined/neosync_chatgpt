/** @type {import('tailwindcss').Config} */
export default {
	content: ['./index.html', './apps/**/*.{js,ts,jsx,tsx}'],
	theme: {
		extend: {
			keyframes: {
				blink: {
					'0%, 100%': { opacity: '1' },
					'50%': { opacity: '0' },
				},
			},
		},
	},
};
