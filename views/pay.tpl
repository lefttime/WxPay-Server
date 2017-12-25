<!DOCTYPE html>
<html lang="zh"><head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
<meta name="viewport" content="user-scalable=no, initial-scale=1.0, maximum-scale=1.0 minimal-ui">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="black">

<title>充值中心</title>
    
<link rel="stylesheet" type="text/css" href="/static/css/style.css">
<link rel="stylesheet" type="text/css" href="/static/css/skin.css">
<link rel="stylesheet" type="text/css" href="/static/css/framework.css">

<script type="text/javascript" src="/static/js/jquery.js"></script>
<script type="text/javascript" src="/static/js/plugins.js"></script>
<script type="text/javascript" src="/static/js/custom.js"></script>
<script type="text/javascript" src="http://pv.sohu.com/cityjson?ie=utf-8"></script>  


    <style>
        @media all and (orientation: portrait) {
            .item{}
        }
        @media all and (orientation: landscape) {
            .item{
            width: 33%;float: left;
            }
            .item h2 {display:none;}
        }
    </style>
</head>

<body>
<div>
    
<div class="page-preloader page-preloader-dark">
    <div class="spinner"></div>
</div>
            
<div id="page-content">
    <div id="page-content-scroll" style="right: 0px;">
        <div class="heading-strip">
            <h3><b id="app-title">摩语网络</b>-选择支付</h3>
            <div class="overlay dark-overlay"></div>
        </div>
        <div class="content" style="font-size: 1.2em;">
            <div class="cart-costs" id="target-area" style="display: none;">
                <h4>充值给</h4>
                <div class="activity-status">
                    <img id="player-avatar" src="/static/img/0" style="width: 45px;height: 45px;float: left;">
                    <strong id="player-name"></strong>
                    <em>ID:<b id="player-id"></b></em>
                </div>
                <h5><strong>商品</strong><em id="target-name">-</em></h5>
                <h5><strong>数量</strong><em id="target-amount">-</em></h5>
                <h5><strong>应付</strong><em class="color-red-light">
                    <b id="target-fee">-</b><b>元</b></em></h5>
                <div class="clear"></div>
                <a href="javascript:;" id="paybtn" class="button button-green button-full half-top">去付款</a>
            </div>
        </div>
        <div class="decoration decoration-margins"></div>
        <div class="content">
            <div class="store-cart-1" id="pkg-list">
                <div class="clear"></div>
                <div class="cart-costs">
                </div>
                <div class="clear"></div>
            </div>
        </div>
        <div class="decoration decoration-margins"></div>
    </div>  
</div>
</div>
<script type="text/javascript" src="/static/js/h5pay.js"></script>
<div class="share-bottom-tap-close"></div></body></html>
