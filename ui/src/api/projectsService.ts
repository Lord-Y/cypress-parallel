import axios, { AxiosInstance } from 'axios'

const axiosClient: AxiosInstance = axios.create({
  headers: {
    'Content-type': 'application/json',
  },
})

export interface Projects {
  team_name: string
  project_id: number
  project_name: string
  branch: string
  cypress_docker_version: string
  browser: string
  config_file: string
  total: number
  date: Date
}

export interface Project {
  team_name: string
  project_id: number
  project_name: string
  date: Date
  repository: string
  username: string
  password: string
  branch: string
  cypress_docker_version: string
  max_pods: number
  scheduling: string
  scheduling_enabled: boolean
  specs: string
  team_id: number
  timeout: number
  browser: string
  config_file: string
  total: number
}

export interface ProjectOnly {
  project_id: number
  project_name: string
  total: number
}

class ProjectsService {
  async create(data: any): Promise<any> {
    return await axiosClient.post('/api/v1/projects', data)
  }

  async get(id: number): Promise<any> {
    return await axiosClient.get<any>(`/api/v1/projects/${id}`)
  }

  async list(page = 1): Promise<any> {
    return await axiosClient.get<any>(`/api/v1/projects/list?page=${page}`)
  }

  async all(): Promise<any> {
    return await axiosClient.get<any>(`/api/v1/projects/all`)
  }

  async update(data: any): Promise<any> {
    return await axiosClient.put('/api/v1/projects', data)
  }

  async delete(id: number): Promise<any> {
    return await axiosClient.delete<any>(`/api/v1/projects/${id}`)
  }

  async search(q: string, page = 1): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/projects/search?q=${q}&page=${page}`,
    )
  }

  async hook(data: any): Promise<any> {
    return await axiosClient.post('/api/v1/hooks/launch/plain', data)
  }
}

export default new ProjectsService()
