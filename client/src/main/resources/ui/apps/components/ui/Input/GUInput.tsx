import { forwardRef } from 'react';

interface GUInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
	error?: string;
}

export const GUInput = forwardRef<HTMLInputElement, GUInputProps>(({ error, className, ...props }, ref) => {
	return (
		<div className="w-full">
			<input ref={ref} {...props} aria-invalid={!!error} className={`px-3 py-2 ${error && 'border-red-500 focus:border-red-500'} ${className}`} />

			{error && <span className="text-sm text-red-500">{error}</span>}
		</div>
	);
});
