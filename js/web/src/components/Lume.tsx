import {Handle, Position} from '@xyflow/react';
import Image from 'next/image';
import {memo} from 'react';
import {LumeType} from '@/genproto/lume/v1/lume_pb';

interface LumeNodeData {
  id: string;
  type: LumeType;
  name: string;
  description?: string;
}

interface LumeNodeProps {
  data: LumeNodeData;
  selected?: boolean;
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

function Lume({data, selected}: LumeNodeProps) {
  const iconPath = lumeTypeToIcon[data.type] || '/file.svg';
  const colorClass = lumeTypeToColor[data.type] || lumeTypeToColor[LumeType.UNSPECIFIED];

  return (
	<>
	  <Handle type="target" position={Position.Top} className="opacity-10"/>
	  <Handle type="target" position={Position.Left} className="opacity-10"/>
	  <div
		className={`
          w-14 h-14 rounded-full bg-gradient-to-br ${colorClass}
          ${selected ? 'ring-2 ring-offset-2 ring-blue-500 shadow-lg' : 'shadow-sm'}
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
	  <Handle type="target" position={Position.Right} className="opacity-10"/>
	  <Handle type="source" position={Position.Bottom} className="opacity-10"/>
	</>
  );
}

export default memo(Lume);
