interface IconProps {
	fill: string;
	size?: number;
	style?: React.CSSProperties;
	className?: string;
	onClick?: () => void;
}

export const FilterVariant: React.FC<IconProps> = ({ fill, size = 24, className, style, onClick }) => {
	return (
		<svg xmlns="http://www.w3.org/2000/svg" className={className} viewBox={`0 0 24 24`} fill={fill} style={style} width={size} height={size} onClick={onClick}>
			<path d="M6,13H18V11H6M3,6V8H21V6M10,18H14V16H10V18Z"></path>
		</svg>
	);
};
export default FilterVariant;
