<?php

namespace App\Model;

use Illuminate\Database\Eloquent\Model;

class Files extends Model
{
    protected $table = 'files';
    /**
     * 可以被批量赋值的属性。
     *
     * @var array
     */
    protected $fillable = [
        'code',
        'path',
        'filename',
        'mime',
    ];
}
