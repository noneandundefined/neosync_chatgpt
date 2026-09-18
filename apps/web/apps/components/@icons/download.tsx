interface IconProps {
	fill: string;
	size?: number;
	style?: React.CSSProperties;
	className?: string;
	onClick?: () => void;
}

export const Download: React.FC<IconProps> = ({ fill, size = 24, className, style, onClick }) => {
	return (
		<svg xmlns="http://www.w3.org/2000/svg" className={className} viewBox={`0 0 24 24`} fill={fill} style={style} width={size} height={size} onClick={onClick}>
			<path d="M5,20H19V18H5M19,9H15V3H9V9H5L12,16L19,9Z" />
		</svg>
	);
};
export default Download;
