import { http } from "./http";
import type { Container, ContainerStatus } from "./types";

export const getContainers = async (projectId: string): Promise<Container[]> => {
    const res = await http.get<Container[]>(`/projects/${projectId}/containers`);
    return res.data;
}

export const getOneContainer = async (projectId: string, containerId: string): Promise<Container> => {
    const res = await http.get<Container>(`/projects/${projectId}/containers/${containerId}`);
    return res.data;
}

export const createContainer = async (projectId: string, containerData: Omit<Container, 'id'>): Promise<Container> => {
    const res = await http.post<Container>(`/projects/${projectId}/containers`, containerData);
    return res.data;
}

export const importComposeFile = async (projectId: string, composeFile: string): Promise<Container[]> => {
    const res = await http.post<Container[]>(`/projects/${projectId}/containers/import`, { compose_file: composeFile });
    return res.data;
}

export const updateContainer = async (projectId: string, containerId: string, containerData: Partial<Container>): Promise<Container> => {
    const res = await http.put<Container>(`/projects/${projectId}/containers/${containerId}`, containerData);
    return res.data;
}

export const deleteContainer = async (projectId: string, containerId: string): Promise<void> => {
    await http.delete(`/projects/${projectId}/containers/${containerId}`);
}

export const startContainer = async (projectId: string, containerId: string): Promise<void> => {
    await http.post(`/projects/${projectId}/containers/${containerId}/start`);
}

export const stopContainer = async (projectId: string, containerId: string): Promise<void> => {
    await http.post(`/projects/${projectId}/containers/${containerId}/stop`);
}


export const getContainerLogs = async (projectId: string, containerId: string, n: number): Promise<string> => {
    const res = await http.get<string>(`/projects/${projectId}/containers/${containerId}/logs?tail=${n}`);
    return res.data;
}

export const getContainerStatus = async (projectId: string, containerId: string): Promise<ContainerStatus> => {
    const res = await http.get(`/projects/${projectId}/containers/${containerId}/status`);
    return res.data.status as ContainerStatus;
}
