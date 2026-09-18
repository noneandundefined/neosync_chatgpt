import { forwardRef } from 'react';

interface GUITextareaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
	error?: string;
}

export const GUITextarea = forwardRef<HTMLTextAreaElement, GUITextareaProps>(({ error, className, ...props }, ref) => {
	return (
		<div className="flex flex-col gap-1 w-full">
			<textarea ref={ref} {...props} aria-invalid={!!error} className={`px-3 py-2 resize-none ${error && 'border-red-500 focus:border-red-500'} ${className}`} />

			{error && <span className="text-sm text-red-500">{error}</span>}
		</div>
	);
});
