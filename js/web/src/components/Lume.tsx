import {Handle, Position, useConnection} from '@xyflow/react';
import Image from 'next/image';
import {memo} from 'react';
import {HoverCard, HoverCardContent, HoverCardTrigger} from '@/components/ui/hover-card';
import {LumeType} from '@/genproto/lume/v1/lume_pb';
import '@/styles/hover-card.css';

interface LumeNodeData {
  id: string;
  type: LumeType;
  name: string;
  description?: string;
}

interface LumeNodeProps {
  data: LumeNodeData;
  selected?: boolean;
  id: string;
}

const lumeTypeToIcon: Record<LumeType, string> = {
  [LumeType.UNSPECIFIED]: '/file.svg',
  [LumeType.CITY]: '/building-2.svg',
  [LumeType.ATTRACTION]: '/map-pin.svg',
  [LumeType.ACCOMMODATION]: '/bed-double.svg',
  [LumeType.RESTAURANT]: '/map-pin.svg',
  [LumeType.TRANSPORT_HUB]: '/bus.svg',
  [LumeType.ACTIVITY]: '/sun.svg',
  [LumeType.SHOPPING]: '/map-pin.svg',
  [LumeType.ENTERTAINMENT]: '/map-pin.svg',
  [LumeType.CUSTOM]: '/file.svg',
};

const lumeTypeToColor: Record<LumeType, string> = {
  [LumeType.UNSPECIFIED]: 'from-stone-50 to-stone-100 border-stone-300',
  [LumeType.CITY]: 'from-blue-50 to-blue-100 border-blue-300',
  [LumeType.ATTRACTION]: 'from-purple-50 to-purple-100 border-purple-300',
  [LumeType.ACCOMMODATION]: 'from-emerald-50 to-emerald-100 border-emerald-300',
  [LumeType.RESTAURANT]: 'from-orange-50 to-orange-100 border-orange-300',
  [LumeType.TRANSPORT_HUB]: 'from-slate-50 to-slate-100 border-slate-300',
  [LumeType.ACTIVITY]: 'from-yellow-50 to-yellow-100 border-yellow-300',
  [LumeType.SHOPPING]: 'from-pink-50 to-pink-100 border-pink-300',
  [LumeType.ENTERTAINMENT]: 'from-indigo-50 to-indigo-100 border-indigo-300',
  [LumeType.CUSTOM]: 'from-gray-50 to-gray-100 border-gray-300',
};

function Lume({data, selected, id}: LumeNodeProps) {
  const iconPath = lumeTypeToIcon[data.type] || '/file.svg';
  const colorClass = lumeTypeToColor[data.type] || lumeTypeToColor[LumeType.UNSPECIFIED];

  const connection = useConnection();

  const isTarget = connection.inProgress && connection.fromNode.id !== id;
  const isConnecting = connection.inProgress && connection.fromNode.id === id;

  return (
	<div className="customNode">
	  {!connection.inProgress && (
		<Handle
		  position={Position.Right}
		  type="source"
		  style={{
			width: '100%',
			height: '100%',
			opacity: 0,
			top: 0,
			left: 0,
			borderRadius: '50%',
			transform: 'none',
		  }}
		/>
	  )}
	  {(!connection.inProgress || isTarget) && (
		<Handle
		  position={Position.Left}
		  type="target"
		  style={{
			width: '100%',
			height: '100%',
			opacity: 0,
			top: 0,
			left: 0,
			borderRadius: '50%',
			transform: 'none',
		  }}
		  isConnectableStart={false}
		/>
	  )}

	  <HoverCard>
		<HoverCardTrigger asChild>
		  <div
			className={`
	          w-14 h-14 rounded-full bg-gradient-to-br ${colorClass}
	          ${selected ? 'ring-2 ring-offset-2 ring-blue-500 shadow-lg' : 'shadow-sm'}
	          ${isConnecting ? 'ring-4 ring-green-400 ring-opacity-75 animate-pulse' : ''}
	          ${isTarget ? 'ring-2 ring-blue-400 ring-opacity-50 hover:ring-opacity-75' : ''}
	          flex items-center justify-center
	          hover:shadow-md hover:scale-105 transition-all duration-200 cursor-pointer
	          border
	        `}
		  >
			<Image
			  src={iconPath}
			  alt={data.name}
			  width={24}
			  height={24}
			  className="opacity-70"
			/>
		  </div>
		</HoverCardTrigger>
		<HoverCardContent className="w-80">
		  <div className="space-y-2">
			<div className="flex items-center gap-2">
			  <Image
				src={iconPath}
				alt={data.name}
				width={16}
				height={16}
				className="opacity-70"
			  />
			  <h4 className="text-sm font-semibold">{data.name}</h4>
			</div>
			{data.description && (
			  <p className="text-sm text-muted-foreground">
				{data.description}
			  </p>
			)}
			<div className="text-xs text-muted-foreground">
			  Type: {LumeType[data.type]}
			</div>
		  </div>
		</HoverCardContent>
	  </HoverCard>
	</div>
  );
}

export default memo(Lume);
