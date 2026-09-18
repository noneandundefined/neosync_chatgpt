interface IconProps {
	fill: string;
	size?: number;
	style?: React.CSSProperties;
	className?: string;
	onClick?: (events?: React.MouseEvent<SVGSVGElement, MouseEvent>) => void;
}

export const Information: React.FC<IconProps> = ({ fill, size = 24, className, style, onClick }) => {
	return (
		<svg xmlns="http://www.w3.org/2000/svg" className={className} viewBox={`0 0 24 24`} fill={fill} style={style} width={size} height={size} onClick={onClick}>
			<path d="M13,9H11V7H13M13,17H11V11H13M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z"></path>
		</svg>
	);
};
export default Information;
