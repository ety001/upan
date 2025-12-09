<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use App\Model\Files;
use Illuminate\Support\Facades\Storage;
use Exception;
use Carbon\Carbon;
use Log;

class FileCleaner extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'file:clean';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Clean the expired files.';

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        try {
            $fileExpireTime = env('FILE_EXPIRE_TIME', 6); // hour
            $expiredTimePoint = new Carbon(date('Y-m-d H:i:s', time() - $fileExpireTime * 3600));
            $files = Files::where('created_at', '<', $expiredTimePoint)->get();
            if ($files->isNotEmpty()) {
                $waitingToDelFiles = [];
                $fileIds = [];
                foreach($files as $f) {
                    array_push($waitingToDelFiles, $f->path);
                    array_push($fileIds, $f->id);
                }
                Storage::delete($waitingToDelFiles);
                Files::destroy($fileIds);
            }
        } catch (Exception $e) {
            Log::error('del_file_err', [$e->getMessage()]);
            $this->error($e->getMessage());
        }
    }
}
