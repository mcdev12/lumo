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
	  connectable: true,
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
	  position: {x: 350, y: 150},
	  connectable: true,
	  data: {
		id: '2',
		type: LumeType.ATTRACTION,
		name: 'Eiffel Tower',
		description: 'Iconic iron lattice tower'
	  },
	},
	{
	  id: '3',
	  type: 'lume',
	  position: {x: 200, y: 250},
	  connectable: true,
	  data: {
		id: '3',
		type: LumeType.RESTAURANT,
		name: 'Le Comptoir',
		description: 'Traditional French bistro'
	  },
	},
	{
	  id: '4',
	  type: 'lume',
	  position: {x: 450, y: 100},
	  connectable: true,
	  data: {
		id: '4',
		type: LumeType.ACCOMMODATION,
		name: 'Hotel Plaza',
		description: 'Luxury hotel in city center'
	  },
	},
	{
	  id: '5',
	  type: 'lume',
	  position: {x: 150, y: 350},
	  connectable: true,
	  data: {
		id: '5',
		type: LumeType.TRANSPORT_HUB,
		name: 'Metro Station',
		description: 'Central metro hub'
	  },
	},
	{
	  id: '6',
	  type: 'lume',
	  position: {x: 400, y: 300},
	  connectable: true,
	  data: {
		id: '6',
		type: LumeType.ACTIVITY,
		name: 'Seine Cruise',
		description: 'Scenic river tour'
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
