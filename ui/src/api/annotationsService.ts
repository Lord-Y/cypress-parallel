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
  date: Date
  total: number
}

export interface Annotation {
  project_id: number
  project_name: string
  annotation_id: number
  key: string
  value: string
  date: Date
}

class AnnotationsService {
  async create(data: any): Promise<any> {
    return await axiosClient.post('/api/v1/annotations', data)
  }

  async get(id: number): Promise<any> {
    return await axiosClient.get<any>(`/api/v1/annotations/${id}`)
  }

  async list(page = 1): Promise<any> {
    return await axiosClient.get<any>(`/api/v1/annotations/list?page=${page}`)
  }

  async update(data: any): Promise<any> {
    return await axiosClient.put('/api/v1/annotations', data)
  }

  async delete(id: number): Promise<any> {
    return await axiosClient.delete<any>(`/api/v1/annotations/${id}`)
  }

  async search(q: string, page = 1): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/annotations/search?q=${q}&page=${page}`,
    )
  }

  async projectID(id: number): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/annotations/list/by/projectid/${id}`,
    )
  }
}

export default new AnnotationsService()
