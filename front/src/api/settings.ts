import { http } from "./http"

import type { Setting } from "./types"

export const getSettings = async (): Promise<Setting[]> => {
  const response = await http.get("/settings")
  return response.data
}

export const getSetting = async (key: string): Promise<Setting> => {
    const response = await http.get(`/settings/${key}`)
    return response.data
}

export const updateSetting = async (setting: Setting): Promise<Setting> => {
  const response = await http.post(`/settings`, setting)
  return response.data
}


export const deleteSetting = async (key: string): Promise<void> => {
  await http.delete(`/settings/${key}`)
}