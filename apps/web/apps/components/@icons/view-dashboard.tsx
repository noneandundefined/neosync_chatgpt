interface IconProps {
	fill: string;
	size?: number;
	style?: React.CSSProperties;
	className?: string;
	onClick?: () => void;
}

export const ViewDashboard: React.FC<IconProps> = ({ fill, size = 24, className, style, onClick }) => {
	return (
		<svg xmlns="http://www.w3.org/2000/svg" className={className} viewBox={`0 0 24 24`} fill={fill} style={style} width={size} height={size} onClick={onClick}>
			<path d="M13,3V9H21V3M13,21H21V11H13M3,21H11V15H3M3,13H11V3H3V13Z"></path>
		</svg>
	);
};
export default ViewDashboard;
