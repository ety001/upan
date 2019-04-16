<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;
use Log;
use Exception;
use Illuminate\Support\Str;
use App\Model\Files as FilesModel;

class MainController extends Controller
{
    public function index(Request $request) {

    }

    public function upload(Request $request, $t = null) {
        $fileMaxSize = env('FILE_MAX_SIZE', 30);
        $fileExpireTime = env('FILE_EXPIRE_TIME', 6);
        try {
            if (!$request->file('o')) {
                throw new Exception('no file input');
            }
            $data = [];
            $data['size'] = $request->file('o')->getClientSize();
            if ($data['size'] > $fileMaxSize * 1024 * 1024) {
                throw new Exception('File size limit: '.$fileMaxSize.'MB');
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
                return redirect('/')->with('status1', 'Upload Success!');
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
}
