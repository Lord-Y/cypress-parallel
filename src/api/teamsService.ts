import axios, { AxiosInstance } from 'axios'

export interface Teams {
  team_id: number
  team_name: string
  total: number
  date: Date
}

const axiosClient: AxiosInstance = axios.create({
  headers: {
    'Content-type': 'application/json',
  },
})

class TeamsService {
  async get(page = 1): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/cypress-parallel-api/teams?page=${page}`,
    )
  }

  async create(data: any): Promise<any> {
    return await axiosClient.post('/api/v1/cypress-parallel-api/teams', data)
  }

  async update(id: number): Promise<any> {
    return await axiosClient.put(`/api/v1/cypress-parallel-api/teams/${id}`)
  }

  async delete(id: number): Promise<any> {
    return await axiosClient.delete<Teams[]>(
      `/api/v1/cypress-parallel-api/teams/${id}`,
    )
  }

  async search(q: string, page = 1): Promise<any> {
    return await axiosClient.get<any>(
      `/api/v1/cypress-parallel-api/teams/search?q=${q}&page=${page}`,
    )
  }
}

export default new TeamsService()
