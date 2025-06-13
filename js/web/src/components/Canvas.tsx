'use client';

import {Background, BackgroundVariant, Controls, MarkerType, MiniMap, ReactFlow,} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import '../styles/react-flow.css';
import useCanvasStore from '@/stores/canvasStore';
import CustomConnectionLine from './CustomConnectionLine';
import FloatingEdge from './FloatingEdge';
import Lume from './Lume';

const nodeTypes = {
  lume: Lume,
};

const edgeTypes = {
  floating: FloatingEdge,
};

const defaultEdgeOptions = {
  type: 'floating',
  markerEnd: {
	type: MarkerType.ArrowClosed,
	color: '#b1b1b7',
  },
};

export default function Canvas() {
  const {nodes, edges, onNodesChange, onEdgesChange, onConnect} = useCanvasStore();
  return (
	<div className="w-full h-screen bg-[#FDFAF6]">
	  <ReactFlow
		nodes={nodes}
		edges={edges}
		onNodesChange={onNodesChange}
		onEdgesChange={onEdgesChange}
		onConnect={onConnect}
		nodeTypes={nodeTypes}
		edgeTypes={edgeTypes}
		connectionLineComponent={CustomConnectionLine}
		connectionLineStyle={{
		  stroke: '#b1b1b7',
		  strokeWidth: 2,
		}}
		fitView
		connectionRadius={25}
		snapToGrid
		snapGrid={[15, 15]}
		defaultEdgeOptions={defaultEdgeOptions}

	  >
		<Background
		  color="rgba(166, 174, 191, 0.3)"
		  variant={BackgroundVariant.Dots}
		  gap={24}
		  size={1.5}
		/>
		<Controls/>
		<MiniMap
		  style={{
			height: 80,
			width: 120,
			backgroundColor: 'rgba(253, 250, 246, 0.75)',
		  }}
		  nodeColor={(node) => {
			const typeColors: Record<string, string> = {
			  lume: '#fed7aa',
			};
			return typeColors[node.type || ''] || '#e5e7eb';
		  }}
		/>
	  </ReactFlow>
	</div>
  );
}
