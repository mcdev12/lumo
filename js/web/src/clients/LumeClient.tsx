import { createConnectTransport } from "@connectrpc/connect-web";
import { createClient } from "@connectrpc/connect";
import {LumeService} from "../genproto/lume/v1/service_pb.js";
import { LumeType } from '../genproto/lume/v1/lume_pb.js';
import { create } from '@bufbuild/protobuf';
import { CreateLumeRequestSchema } from '../genproto/lume/v1/service_pb.js';

const transport = createConnectTransport({
    baseUrl: "http://localhost:8080",
});

export const lumeClient = createClient(LumeService, transport);

// Use Web Crypto API for UUID generation
function generateUUID(): string {
    return crypto.randomUUID();
}

export const createLume = async (
    name: string,
    description: string,
    latitude: number,
    longitude: number,
    address: string,
    type: LumeType = LumeType.CITY
) => {
    const request = create(CreateLumeRequestSchema, {
        lumoId: generateUUID(), // Generate a valid UUID
        type: type,
        name: name,
        latitude: latitude,
        longitude: longitude,
        address: address,
        description: description,
        images: [],
        categoryTags: [],
        bookingLink: ""
    });

    return await lumeClient.createLume(request);
};