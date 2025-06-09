import React from 'react';
import { Lume, LumeType } from '../genproto/lume/v1/lume_pb.js';

interface LumeCardProps {
    lume: Lume;
}

const getLumeTypeName = (type: LumeType): string => {
    switch (type) {
        case LumeType.CITY:
            return "City";
        case LumeType.ATTRACTION:
            return "Attraction";
        case LumeType.RESTAURANT:
            return "Restaurant";
        default:
            return "Unknown";
    }
};

export const LumeCard: React.FC<LumeCardProps> = ({ lume }) => {
    return (
        <div className="lume-item">
            <h3>{lume.name}</h3>
            <p><strong>Type:</strong> {getLumeTypeName(lume.type)}</p>
            <p><strong>ID:</strong> {lume.lumeId}</p>
            <p><strong>Description:</strong> {lume.description}</p>
            <p><strong>Address:</strong> {lume.address}</p>
            <p><strong>Coordinates:</strong> {lume.latitude}, {lume.longitude}</p>
            {lume.categoryTags.length > 0 && (
                <p><strong>Tags:</strong> {lume.categoryTags.join(", ")}</p>
            )}
            {lume.bookingLink && (
                <p><strong>Booking:</strong> <a href={lume.bookingLink} target="_blank" rel="noopener noreferrer">Link</a></p>
            )}
        </div>
    );
};