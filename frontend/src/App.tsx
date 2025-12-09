import { useEffect, useState } from 'react'
import Upan from './components/Upan'
import { getConfig } from './api/client'

interface Config {
  file_max_size: number
  file_expire_time: number
}

function App() {
  const [config, setConfig] = useState<Config | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getConfig()
      .then(setConfig)
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg">加载中...</div>
      </div>
    )
  }

  if (!config) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg text-red-500">加载配置失败</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="container mx-auto px-4 max-w-2xl">
        <h1 className="text-3xl font-bold text-center mb-8">网络闪存</h1>
        <Upan fileMaxSize={config.file_max_size} fileExpireTime={config.file_expire_time} />
      </div>
    </div>
  )
}

export default App

