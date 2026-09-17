import React, { useMemo } from 'react';

export interface GUIDonutSegment {
	key: string;
	value: number;
	color: string;
}

interface GUIDonutProps {
	segments: GUIDonutSegment[];
	size?: number;
	stroke?: number;
	trackColor?: string;
	center?: React.ReactNode;
	className?: string;
}

const GUIDonut: React.FC<GUIDonutProps> = ({ segments, size = 160, stroke = 20, trackColor = '#eef1f4', center, className }) => {
	const radius = (size - stroke) / 2;
	const circumference = 2 * Math.PI * radius;

	const arcs = useMemo(() => {
		const total = segments.reduce((sum, segment) => sum + segment.value, 0);
		let offset = 0;

		return segments.map((segment) => {
			const length = total > 0 ? (segment.value / total) * circumference : 0;
			const arc = {
				...segment,
				length,
				dashOffset: -offset,
			};

			offset += length;

			return arc;
		});
	}, [segments, circumference]);

	return (
		<div className={`relative shrink-0 ${className ?? ''}`} style={{ width: size, height: size }}>
			<svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} className="-rotate-90">
				<circle cx={size / 2} cy={size / 2} r={radius} fill="none" stroke={trackColor} strokeWidth={stroke} />

				{arcs.map((arc) =>
					arc.length > 0 ? (
						<circle
							key={arc.key}
							cx={size / 2}
							cy={size / 2}
							r={radius}
							fill="none"
							stroke={arc.color}
							strokeWidth={stroke}
							strokeDasharray={`${arc.length} ${circumference - arc.length}`}
							strokeDashoffset={arc.dashOffset}
							strokeLinecap="butt"
						/>
					) : null
				)}
			</svg>

			{center ? <div className="absolute inset-0 flex flex-col items-center justify-center text-center pointer-events-none">{center}</div> : null}
		</div>
	);
};

export default GUIDonut;
