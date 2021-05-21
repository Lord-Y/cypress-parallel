export interface Loading {
  loading: boolean
  class: string
}

export interface Pagination {
  url: string
  actualPage: number
  total: number
  mode?: string
}
