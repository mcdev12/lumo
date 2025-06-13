import {
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
  Connection,
  Edge,
  EdgeChange,
  Node,
  NodeChange
} from '@xyflow/react';
import {create} from 'zustand';
import {LumeType} from '@/genproto/lume/v1/lume_pb';

interface CanvasState {
  nodes: Node[];
  edges: Edge[];
  onNodesChange: (changes: NodeChange[]) => void;
  onEdgesChange: (changes: EdgeChange[]) => void;
  onConnect: (connection: Connection) => void;
  addNode: (node: Node) => void;
  updateNode: (id: string, data: Record<string, unknown>) => void;
  deleteNode: (id: string) => void;
  setNodes: (nodes: Node[]) => void;
  setEdges: (edges: Edge[]) => void;
}

const useCanvasStore = create<CanvasState>((set, get) => ({
  nodes: [
	{
	  id: '1',
	  type: 'lume',
	  position: {x: 100, y: 100},
	  data: {
		id: '1',
		type: LumeType.CITY,
		name: 'Paris',
		description: 'The City of Light'
	  },
	},
	{
	  id: '2',
	  type: 'lume',
	  position: {x: 300, y: 200},
	  data: {
		id: '2',
		type: LumeType.ATTRACTION,
		name: 'Eiffel Tower',
		description: 'Iconic iron lattice tower'
	  },
	},
  ],
  edges: [],
  onNodesChange: (changes) => {
	set({
	  nodes: applyNodeChanges(changes, get().nodes),
	});
  },
  onEdgesChange: (changes) => {
	set({
	  edges: applyEdgeChanges(changes, get().edges),
	});
  },
  onConnect: (connection) => {
	set({
	  edges: addEdge(connection, get().edges),
	});
  },
  addNode: (node) => {
	set({
	  nodes: [...get().nodes, node],
	});
  },
  updateNode: (id, data) => {
	set({
	  nodes: get().nodes.map((node) =>
		node.id === id ? {...node, data: {...node.data, ...data}} : node
	  ),
	});
  },
  deleteNode: (id) => {
	set({
	  nodes: get().nodes.filter((node) => node.id !== id),
	  edges: get().edges.filter((edge) => edge.source !== id && edge.target !== id),
	});
  },
  setNodes: (nodes) => {
	set({nodes});
  },
  setEdges: (edges) => {
	set({edges});
  },
}));

export default useCanvasStore;
