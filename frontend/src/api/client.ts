const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'

export interface Config {
  file_max_size: number
  file_expire_time: number
}

export interface UploadResponse {
  status: boolean
  code?: string
  expired_at?: number
  error?: string
}

export async function getConfig(): Promise<Config> {
  const response = await fetch(`${API_BASE_URL}/config`)
  if (!response.ok) {
    throw new Error('Failed to fetch config')
  }
  return response.json()
}

export async function uploadFile(file: File): Promise<UploadResponse> {
  const formData = new FormData()
  formData.append('o', file)

  const response = await fetch(`${API_BASE_URL}/upload`, {
    method: 'POST',
    body: formData,
  })

  return response.json()
}

export function downloadFile(code: string): void {
  window.location.href = `${API_BASE_URL}/file/${code}`
}

