import { CSSProperties } from 'react';

interface GUIRangeProps {
	max: number;
	min?: number;
	value: number;
	onChange: (value: number) => void;
	onChangeComplete?: (value: number) => void;
	disabled?: boolean;
	className?: string;
}

const GUIRange: React.FC<GUIRangeProps> = ({ max, min = 0, value, onChange, onChangeComplete, disabled, className }) => {
	const span = Math.max(max - min, 1);
	const percent = ((value - min) / span) * 100;
	const ticks = Array.from({ length: span + 1 }, (_, i) => min + i);

	const commit = (next: number) => {
		onChangeComplete?.(next);
	};

	return (
		<div className={`gui-range ${className || ''}`}>
			<input
				type="range"
				min={min}
				max={max}
				step={1}
				value={value}
				disabled={disabled}
				style={{ '--gui-range-fill': `${percent}%` } as CSSProperties}
				onChange={(e) => onChange(Number(e.target.value))}
				onMouseUp={(e) => commit(Number((e.target as HTMLInputElement).value))}
				onTouchEnd={(e) => commit(Number((e.target as HTMLInputElement).value))}
			/>

			<div className="gui-range__ticks">
				{ticks.map((tick) => (
					<span key={tick} className="gui-range__tick">
						{(tick === min || tick === max) && <span className="gui-range__label">{tick}</span>}
					</span>
				))}
			</div>
		</div>
	);
};

export default GUIRange;
