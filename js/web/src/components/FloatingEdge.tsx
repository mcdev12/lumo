import { BaseEdge, getStraightPath, useInternalNode, EdgeProps } from '@xyflow/react';
import { getEdgeParams } from '../utils/edgeUtils.js';

function FloatingEdge({ id, source, target, markerEnd, style }: EdgeProps) {
  const sourceNode = useInternalNode(source);
  const targetNode = useInternalNode(target);

  if (!sourceNode || !targetNode) {
    return null;
  }

  // Transform the nodes to match the expected interface
  const sourceNodeWithInternals = {
    ...sourceNode,
    internals: {
      positionAbsolute: sourceNode.internals.positionAbsolute,
    },
    measured: {
      width: sourceNode.measured.width || 0,
      height: sourceNode.measured.height || 0,
    },
  };

  const targetNodeWithInternals = {
    ...targetNode,
    internals: {
      positionAbsolute: targetNode.internals.positionAbsolute,
    },
    measured: {
      width: targetNode.measured.width || 0,
      height: targetNode.measured.height || 0,
    },
  };

  const { sx, sy, tx, ty } = getEdgeParams(sourceNodeWithInternals, targetNodeWithInternals);

  const [path] = getStraightPath({
    sourceX: sx,
    sourceY: sy,
    targetX: tx,
    targetY: ty,
  });

  return (
    <BaseEdge
      id={id}
      className="react-flow__edge-path"
      path={path}
      markerEnd={markerEnd}
      style={style}
    />
  );
}

export default FloatingEdge;
