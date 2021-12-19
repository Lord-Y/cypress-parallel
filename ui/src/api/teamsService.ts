import axios, { AxiosInstance } from 'axios'

const axiosClient: AxiosInstance = axios.create({
  headers: {
    'Content-type': 'application/json',
  },
})

export interface Teams {
  team_id: number
  team_name: string
  total: number
  date: Date
}

export interface Team {
  team_id: number
  team_name: string
  date: Date
}

class TeamsService {
  async create(data: any): Promise<any> {
    return await axiosClient.post('/api/v1/teams', data)
  }

  async get(id: number): Promise<any> {
    return await axiosClient.get<any>(`/api/v1/teams/${id}`)
  }

  async list(page = 1): Promise<any> {
    return await axiosClient.get<any>(`/api/v1/teams/list?page=${page}`)
  }

  async all(): Promise<any> {
    return await axiosClient.get<any>(`/api/v1/teams/all`)
  }

  async update(data: any): Promise<any> {
    return await axiosClient.put('/api/v1/teams', data)
  }

  async delete(id: number): Promise<any> {
    return await axiosClient.delete<any>(`/api/v1/teams/${id}`)
  }

  async search(q: string, page = 1): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/teams/search?q=${q}&page=${page}`,
    )
  }
}

export default new TeamsService()
