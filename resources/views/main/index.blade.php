@extends('layouts.app')

@section('title', $sitename.' -- By ETY001')

@section('body')
    <upan
        file-max-size="{{$fileMaxSize}}"
        file-expire-time="{{$fileExpireTime}}"
    ></upan>
@endsection

@section('customjs')
@endsection
