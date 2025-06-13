'use client';

import {Background, BackgroundVariant, Controls, MiniMap, ReactFlow,} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import useCanvasStore from '@/stores/canvasStore';

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
                fitView
            >
                <Background
                    color="rgba(166, 174, 191, 0.3)"
                    variant={BackgroundVariant.Dots}
                    gap={20}
                    size={2}
                />
                <Controls/>
                <MiniMap
                    style={{
                        height: 80,
                        width: 120,
                    }}
                />
            </ReactFlow>
        </div>
    );
}
