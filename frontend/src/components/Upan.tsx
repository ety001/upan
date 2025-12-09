import { useState } from 'react'
import { uploadFile, downloadFile } from '../api/client'

interface UpanProps {
  fileMaxSize: number
  fileExpireTime: number
}

export default function Upan({ fileMaxSize, fileExpireTime }: UpanProps) {
  const [code, setCode] = useState('')
  const [uploadStatus, setUploadStatus] = useState<{ type: 'success' | 'error'; message: string } | null>(null)
  const [uploading, setUploading] = useState(false)

  const handleFileUpload = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setUploadStatus(null)
    setUploading(true)

    const form = e.currentTarget
    const fileInput = form.querySelector<HTMLInputElement>('input[type="file"]')
    
    if (!fileInput || !fileInput.files || fileInput.files.length === 0) {
      setUploadStatus({ type: 'error', message: '请选择要上传的文件!' })
      setUploading(false)
      return
    }

    const file = fileInput.files[0]

    // Check file size
    if (file.size > fileMaxSize * 1024 * 1024) {
      setUploadStatus({ 
        type: 'error', 
        message: `文件大小限制: ${fileMaxSize}MB` 
      })
      setUploading(false)
      return
    }

    try {
      const result = await uploadFile(file)
      if (result.status) {
        setUploadStatus({ 
          type: 'success', 
          message: `上传成功! 你的提取码是: ${result.code}` 
        })
        form.reset()
      } else {
        setUploadStatus({ 
          type: 'error', 
          message: result.error || '上传失败' 
        })
      }
    } catch (error) {
      setUploadStatus({ 
        type: 'error', 
        message: '上传失败，请重试' 
      })
    } finally {
      setUploading(false)
    }
  }

  const handleDownload = () => {
    if (code.trim()) {
      downloadFile(code.trim())
    }
  }

  return (
    <div className="space-y-5">
      {/* Download Section */}
      <div className="bg-white rounded-lg shadow-md">
        <div className="px-6 py-4 border-b border-gray-200">
          <h2 className="text-xl font-semibold">获取文件</h2>
        </div>
        <div className="px-6 py-4">
          <form>
            <div className="mb-4">
              <label htmlFor="code" className="block text-sm font-medium text-gray-700 mb-2">
                下载码
              </label>
              <input
                type="text"
                id="code"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="请输入下载码"
              />
            </div>
            <button
              type="button"
              onClick={handleDownload}
              className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors"
            >
              下载
            </button>
          </form>
        </div>
      </div>

      {/* Upload Section */}
      <div className="bg-white rounded-lg shadow-md">
        <div className="px-6 py-4 border-b border-gray-200">
          <h2 className="text-xl font-semibold">上传文件</h2>
        </div>
        <div className="px-6 py-4">
          <form onSubmit={handleFileUpload} encType="multipart/form-data">
            <div className="mb-4">
              <input
                type="file"
                name="o"
                className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
                disabled={uploading}
              />
              <p className="mt-2 text-sm text-gray-500">
                文件大小限制: {fileMaxSize}MB, 保留时间: {fileExpireTime}小时
              </p>
            </div>
            <button
              type="submit"
              disabled={uploading}
              className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed"
            >
              {uploading ? '上传中...' : '上传'}
            </button>
          </form>
        </div>
      </div>

      {/* Status Messages */}
      {uploadStatus && (
        <div
          className={`p-4 rounded-md ${
            uploadStatus.type === 'success'
              ? 'bg-green-50 text-green-800 border border-green-200'
              : 'bg-red-50 text-red-800 border border-red-200'
          }`}
        >
          <div className="flex items-center">
            {uploadStatus.type === 'success' ? (
              <svg
                className="w-5 h-5 mr-2"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                  clipRule="evenodd"
                />
              </svg>
            ) : (
              <svg
                className="w-5 h-5 mr-2"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  clipRule="evenodd"
                />
              </svg>
            )}
            <span dangerouslySetInnerHTML={{ __html: uploadStatus.message }} />
          </div>
        </div>
      )}
    </div>
  )
}

