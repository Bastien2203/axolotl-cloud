
export type Project = {
  id: string
  name: string
  icon_url: string
  created_at: string
  updated_at: string
}

export type Container = {
  id: string
  name: string
  docker_image: string
  ports: Record<string, string>
  env: Record<string, string>
  volumes: Record<string, string>
  networks: string[]
}

export type ContainerStatus = "created" | "running" | "paused" | "restarting" | "removing" | "exited" | "dead"

export const statusColors: Record<ContainerStatus, string> = {
  created: "bg-gray-200 text-gray-800",
  running: "bg-green-200 text-green-800",
  paused: "bg-yellow-200 text-yellow-800",
  restarting: "bg-blue-200 text-blue-800",
  removing: "bg-red-200 text-red-800",
  exited: "bg-gray-400 text-gray-800",
  dead: "bg-black text-white",
}