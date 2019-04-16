<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;
use Log;
use Exception;
use App\Model\Files as FilesModel;

class MainController extends Controller
{
    public function index(Request $request) {
        $fileMaxSize = env('FILE_MAX_SIZE', 30);
        $fileExpireTime = env('FILE_EXPIRE_TIME', 6);
        return response()->view(
            'main/index',
            [
                'sitename' => '网络闪存',
                'fileMaxSize' => $fileMaxSize,
                'fileExpireTime' => $fileExpireTime,
            ],
            200
        );
    }

    public function upload(Request $request, $t = null) {
        $fileMaxSize = env('FILE_MAX_SIZE', 30);
        $fileExpireTime = env('FILE_EXPIRE_TIME', 6);
        try {
            if (!$request->file('o')) {
                throw new Exception('请选择要上传的文件!');
            }
            $data = [];
            $data['size'] = $request->file('o')->getClientSize();
            if ($data['size'] > $fileMaxSize * 1024 * 1024) {
                throw new Exception('文件大小限制: '.$fileMaxSize.'MB');
            }
            $data['filename'] = $request->file('o')->getClientOriginalName();
            $data['mime'] = $request->file('o')->getClientMimeType();
            $data['path'] = Storage::putFile('files', $request->file('o'));
            $data['code'] = bin2hex(random_bytes(3));
            while (FilesModel::where('code', $data['code'])->count() > 0) {
                $data['code'] = bin2hex(random_bytes(3));
            }
            $file = FilesModel::create($data);
            if ($t != 'api') {
                return redirect('/')->with('status1', '上传成功! 你的提取码是: <span>'.$file->code.'</span>');
            } else {
                return response()->json([
                    'status' => true,
                    'code' => $file->code,
                    'expired_at' => strtotime($file->created_at) + $fileExpireTime * 3600,
                ]);
            }
        } catch(Exception $e) {
            Log::error('upload_error', [$e->getMessage(), $request->input()]);
            if ($t != 'api') {
                return redirect('/')->with('status0', $e->getMessage());
            } else {
                return response()->json([
                    'status' => false,
                    'error' => $e->getMessage(),
                ]);
            }
        }
    }

    public function getFile(Request $request, $code) {
        $file = FilesModel::where('code', $code)->first();
        if ($file) {
            return Storage::download($file->path, $file->filename, []);
        } else {
            return redirect('/')->with('status0', '提取码不存在');
        }
    }
}
