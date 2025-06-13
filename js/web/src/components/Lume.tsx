import { memo } from 'react';
import Image from 'next/image';
import { Handle, Position } from '@xyflow/react';
import { LumeType } from '@/genproto/lume/v1/lume_pb';

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

function Lume({ data, selected }: LumeNodeProps) {
  const iconPath = lumeTypeToIcon[data.type] || '/file.svg';

  return (
    <>
      <Handle type="target" position={Position.Top} />
      <div
        className={`
          w-16 h-16 rounded-full bg-white border-2 
          ${selected ? 'border-blue-500 shadow-lg' : 'border-gray-300'}
          flex items-center justify-center
          hover:shadow-md transition-shadow cursor-pointer
        `}
      >
        <Image
          src={iconPath}
          alt={data.name}
          width={32}
          height={32}
          className="text-gray-700"
        />
      </div>
      <div className="text-center mt-2 max-w-[100px]">
        <div className="text-sm font-medium truncate">{data.name}</div>
      </div>
      <Handle type="source" position={Position.Bottom} />
    </>
  );
}

export default memo(Lume);