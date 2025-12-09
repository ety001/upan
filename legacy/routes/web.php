<?php

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| contains the "web" middleware group. Now create something great!
|
*/

Route::get('/', 'MainController@index')->name('main_index');
Route::post('/upload/{t?}', 'MainController@upload')->name('main_upload');
Route::get('/file/{code}', 'MainController@getFile')->name('main_getfile');
