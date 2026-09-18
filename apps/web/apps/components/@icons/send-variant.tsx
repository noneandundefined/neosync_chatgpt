interface IconProps {
	fill: string;
	size?: number;
	style?: React.CSSProperties;
	className?: string;
	onClick?: () => void;
}

export const SentVariant: React.FC<IconProps> = ({ fill, size = 24, className, style, onClick }) => {
	return (
		<svg xmlns="http://www.w3.org/2000/svg" className={className} viewBox={`0 0 24 24`} fill={fill} style={style} width={size} height={size} onClick={onClick}>
			<path d="M3 20V14L11 12L3 10V4L22 12Z" />
		</svg>
	);
};
export default SentVariant;
