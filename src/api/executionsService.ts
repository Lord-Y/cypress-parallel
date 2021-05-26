import axios, { AxiosInstance } from 'axios'

const axiosClient: AxiosInstance = axios.create({
  headers: {
    'Content-type': 'application/json',
  },
})

export interface Execution {
  execution_id: number
  project_id: number
  branch: string
  execution_status: string
  uniq_id: string
  spec: string
  result: string
  date: Date
  execution_error_output: string
  pod_name: string
  pod_cleaned: boolean
  total: number
}

class ExecutionsService {
  async get(id: number): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/cypress-parallel-api/executions/${id}`,
    )
  }

  async list(page = 1): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/cypress-parallel-api/executions/list?page=${page}`,
    )
  }

  async uniqid(id: string): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/cypress-parallel-api/executions/list/by/uniqid/${id}`,
    )
  }

  async search(q: string, page = 1): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/cypress-parallel-api/executions/search?q=${q}&page=${page}`,
    )
  }
}

export default new ExecutionsService()
