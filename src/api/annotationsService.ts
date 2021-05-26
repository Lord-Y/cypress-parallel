import axios, { AxiosInstance } from 'axios'

const axiosClient: AxiosInstance = axios.create({
  headers: {
    'Content-type': 'application/json',
  },
})

export interface Annotations {
  project_id: number
  project_name: string
  annotation_id: number
  key: string
  value: string
  total: number
}

export interface Annotation {
  project_id: number
  project_name: string
  annotation_id: number
  key: string
  value: string
}

class AnnotationsService {
  async create(data: any): Promise<any> {
    return await axiosClient.post(
      '/api/v1/cypress-parallel-api/annotations',
      data,
    )
  }

  async get(id: number): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/cypress-parallel-api/annotations/${id}`,
    )
  }

  async list(page = 1): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/cypress-parallel-api/annotations/list?page=${page}`,
    )
  }

  async update(data: any): Promise<any> {
    return await axiosClient.put(
      '/api/v1/cypress-parallel-api/annotations',
      data,
    )
  }

  async delete(id: number): Promise<any> {
    return await axiosClient.delete<Annotations[]>(
      `/api/v1/cypress-parallel-api/annotations/${id}`,
    )
  }

  async search(q: string, page = 1): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/cypress-parallel-api/annotations/search?q=${q}&page=${page}`,
    )
  }
}

export default new AnnotationsService()
