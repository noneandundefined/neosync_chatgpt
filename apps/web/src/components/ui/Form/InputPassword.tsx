import Eye from '@/components/@icons/eye';
import { forwardRef, useState } from 'react';
import EyeOff from '@/components/@icons/eye-off';

interface InputPasswordProps extends React.InputHTMLAttributes<HTMLInputElement> {
	label?: string;
	error?: string;
}

const InputPassword = forwardRef<HTMLInputElement, InputPasswordProps>(({ label, error, ...props }, ref) => {
	const [showPassword, setShowPassword] = useState(false);

	// const handlePreventClipboard = (e: any) => {
	//     e.preventDefault();
	// };

	return (
		<div className="w-full">
			{label && <label className="block text-[15px] font-medium text-gray-800 mb-1">{label}</label>}

			<div className="relative">
				<input
					ref={ref}
					{...props}
					spellCheck="false"
					autoCorrect="off"
					autoCapitalize="off"
					aria-invalid={!!error}
					type={showPassword ? 'text' : 'password'}
					className={`px-3 pr-10 ${error && 'border-red-500 focus:border-red-500'} ${props.className || ''}`}
				/>

				<span onClick={() => setShowPassword((prev) => !prev)} className="cursor-pointer absolute right-3 top-[50%] -translate-y-1/2">
					{showPassword ? <EyeOff size={17} fill="#999" /> : <Eye size={17} fill="#999" />}
				</span>
			</div>

			{error && <span className="text-sm text-red-500">{error}</span>}
		</div>
	);
});

export default InputPassword;
