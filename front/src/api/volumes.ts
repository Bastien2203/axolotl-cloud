import { http } from "./http";
import type { Volume } from "./types";

export const getVolumes = async (): Promise<Volume[]> => {
    const res = await http.get<Volume[]>("/volumes");
    return res.data;
}