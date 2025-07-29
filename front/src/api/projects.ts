import { http } from "./http"
import type { Project } from "./types"

export const getProjects = async (): Promise<Project[]> => {
  const res = await http.get<Project[]>("/projects")
  return res.data
}

export const createProject = async (data: Omit<Project, 'id' | 'created_at' | 'updated_at'>) : Promise<Project> => {
  const res = await http.post("/projects", data)
  return res.data
}

export const deleteProject = async (id: string) => {
  await http.delete(`/projects/${id}`)
}

export const updateProject = async (id: string, data: Omit<Project, 'id' | 'created_at' | 'updated_at'>): Promise<Project> => {
  const res = await http.put(`/projects/${id}`, data)
  return res.data
}

export const getProject = async (id: string): Promise<Project> => {
  const res = await http.get<Project>(`/projects/${id}`)
  return res.data
}

export const runProject = async (id: string) => {
    await http.post(`/projects/${id}/run`)
}

export const stopProject = async (id: string) => {
    await http.post(`/projects/${id}/stop`)
}

export const getProjectLogs = async (id: string, n: number): Promise<string> => {
  const res = await http.get<string>(`/projects/${id}/logs?n=${n}`)
  return res.data
}



