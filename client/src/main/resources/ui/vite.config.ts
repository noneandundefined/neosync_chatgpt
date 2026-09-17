import { defineConfig } from 'vite';
import { visualizer } from 'rollup-plugin-visualizer';
import react from '@vitejs/plugin-react-swc';
import path from 'path';

// https://vite.dev/config/
export default defineConfig({
	plugins: [
		react(),
		visualizer({
			open: false,
			filename: 'bundle-report.html',
		}),
	],
	resolve: {
		alias: {
			'@': path.resolve(__dirname, './apps'),
		},
	},
});
