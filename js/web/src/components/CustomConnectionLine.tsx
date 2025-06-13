import { getStraightPath } from '@xyflow/react';
import React from 'react';

interface CustomConnectionLineProps {
  fromX: number;
  fromY: number;
  toX: number;
  toY: number;
  connectionLineStyle?: React.CSSProperties;
}

function CustomConnectionLine({ fromX, fromY, toX, toY, connectionLineStyle }: CustomConnectionLineProps) {
  const [edgePath] = getStraightPath({
    sourceX: fromX,
    sourceY: fromY,
    targetX: toX,
    targetY: toY,
  });

  return (
    <g>
      <path 
        style={{
          stroke: '#b1b1b7',
          strokeWidth: 2,
          fill: 'none',
          ...connectionLineStyle
        }} 
        d={edgePath} 
      />
    </g>
  );
}

export default CustomConnectionLine;
