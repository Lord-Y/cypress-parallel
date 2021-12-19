import axios, { AxiosInstance } from 'axios'

const axiosClient: AxiosInstance = axios.create({
  headers: {
    'Content-type': 'application/json',
  },
})

export interface Environments {
  project_id: number
  project_name: string
  environment_id: number
  key: string
  value: string
  date: Date
  total: number
}

export interface Environment {
  project_id: number
  project_name: string
  environment_id: number
  key: string
  value: string
  date: Date
}

class EnvironmentsService {
  async create(data: any): Promise<any> {
    return await axiosClient.post('/api/v1/environments', data)
  }

  async get(id: number): Promise<any> {
    return await axiosClient.get<any>(`/api/v1/environments/${id}`)
  }

  async list(page = 1): Promise<any> {
    return await axiosClient.get<any>(`/api/v1/environments/list?page=${page}`)
  }

  async update(data: any): Promise<any> {
    return await axiosClient.put<any>('/api/v1/environments', data)
  }

  async delete(id: number): Promise<any> {
    return await axiosClient.delete<any>(`/api/v1/environments/${id}`)
  }

  async search(q: string, page = 1): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/environments/search?q=${q}&page=${page}`,
    )
  }

  async projectID(id: number): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/environments/list/by/projectid/${id}`,
    )
  }
}

export default new EnvironmentsService()
