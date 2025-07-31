import { http } from "./http";
import type { Job } from "./types";


export const getJobs = async (): Promise<Job[]> => {
    const res = await http.get<Job[]>('/jobs');
    return res.data;
}

export const getJob = async (jobId: string): Promise<Job> => {
    const res = await http.get<Job>(`/jobs/${jobId}`);
    return res.data;
}

export const removeJob = async (jobId: string): Promise<void> => {
    await http.delete(`/jobs/${jobId}`);
}