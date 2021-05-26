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
  total: number
}

export interface Environment {
  project_id: number
  project_name: string
  environment_id: number
  key: string
  value: string
}

class EnvironmentsService {
  async create(data: any): Promise<any> {
    return await axiosClient.post(
      '/api/v1/cypress-parallel-api/environments',
      data,
    )
  }

  async get(id: number): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/cypress-parallel-api/environments/${id}`,
    )
  }

  async list(page = 1): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/cypress-parallel-api/environments/list?page=${page}`,
    )
  }

  async update(data: any): Promise<any> {
    return await axiosClient.put(
      '/api/v1/cypress-parallel-api/environments',
      data,
    )
  }

  async delete(id: number): Promise<any> {
    return await axiosClient.delete<Environments[]>(
      `/api/v1/cypress-parallel-api/environments/${id}`,
    )
  }

  async search(q: string, page = 1): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/cypress-parallel-api/environments/search?q=${q}&page=${page}`,
    )
  }
}

export default new EnvironmentsService()
