<!doctype html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, minimum-scale=1.0, maximum-scale=1.0, user-scalable=no"/>
        <title>@yield('title')</title>
        <link rel="stylesheet" href="/css/app.css">
        <meta name="csrf-token" content="{{ csrf_token() }}">
        @yield('customcss')
        <!-- Global site tag (gtag.js) - Google Analytics -->
        <script async src="https://www.googletagmanager.com/gtag/js?id=UA-132752007-2"></script>
        <script>
            window.dataLayer = window.dataLayer || [];
            function gtag(){dataLayer.push(arguments);}
            gtag('js', new Date());

            gtag('config', 'UA-132752007-2');
        </script>
        <script async src="//pagead2.googlesyndication.com/pagead/js/adsbygoogle.js"></script>
        <script>
            (adsbygoogle = window.adsbygoogle || []).push({
                google_ad_client: "ca-pub-7536831447654223",
                enable_page_level_ads: true
            });
        </script>
    </head>
    <body style="background-color: #eee;">
        <div id="app" class="container">
            <div class="row">
                <div class="col col-md-6 offset-md-3">
                    <div class="text-center">
                        <h1 style="margin-top: 10px;font-size: 2em;">{{ $sitename }}</h1>
                    </div>
                </div>
            </div>
            <div class="row">
                <div class="col col-md-6 offset-md-3">
                    @if (session('status1'))
                        <div class="alert alert-success">
                            {!! session('status1') !!}
                        </div>
                    @endif
                    @if (session('status0'))
                        <div class="alert alert-danger">
                            {!! session('status0') !!}
                        </div>
                    @endif
                </div>
            </div>
            <div class="row">
                <div class="col col-md-6 offset-md-3">
                    @yield('body')
                </div>
            </div>
            <div class="row">
                <div class="col col-md-6 offset-md-3">
                    <hr />
                </div>
            </div>
            <div class="row">
                <div class="col col-md-6 offset-md-3 text-center">
                @ <a href="https://blog.domyself.me">ety001</a>
                </div>
            </div>
        </div>
        <script src="/js/app.js"></script>
        @yield('customjs')
    </body>
</html>
