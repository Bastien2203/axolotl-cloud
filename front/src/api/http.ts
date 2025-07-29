import axios from "axios"

const env = import.meta.env.VITE_APP_ENV;
const API_HOST = env == "production" ? document.location.origin : "http://localhost:8888";

export const http = axios.create({
  baseURL: API_HOST + "/api",
  timeout: 5000
})
