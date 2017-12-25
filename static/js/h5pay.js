var getJsonFromUrl = function() {
  var query  = location.search.substr( 1 );
  var result = {};
  query.split("&").forEach(function(part) {
    var item = part.split( "=" );
    result[item[0]] = decodeURIComponent( item[1] );
  });
  return result;
};

var selectedPkg = null;
var ipaddr      = '127.0.0.1';
var getParam = function( key ) {
  var params = getJsonFromUrl();
  if( params && params.hasOwnProperty( key ) ){
    return params[key];
  }
  return null;
};

var appInfos = {
  'guangxi': {
    'title': '充值中心',
    'name' : '摩语广西'
  },
  'hubei': {
    'title': '充值中心',
    'name' : '摩语荆楚'
  },
  'xiangyang': {
    'title': '充值中心',
    'name' : '幸运卡五星'
  }
}

var baseUrl = 'http://127.0.0.1:8080/pay';
var pkgs    = [];
var getH5Pkgs = function( app, cb ){
  var url = baseUrl + '/productInfo?area=' + app;
  $.get( url, function( d ) {
    if( d.status != 0 ) {
      cb( [] );
    } else {
      cb( d.data )
    }
  });
};

var fetchPlayerInfoById = function() {
  var url = baseUrl + '/playerInfo?uid=' + getParam( 'uid' ) + '&area=' + getParam( 'app' );
  $.get( url, function( d ) {
    if( d.status != 0 ){
      alert( d.message );
      return;
    };
    var playerInfo = JSON.parse( d.data )
    $( '#player-id'     ).text( playerInfo.id       );
    $( '#player-name'   ).text( playerInfo.name     );
    $( '#player-cards'  ).text( playerInfo.diamonds );
    $( '#player-avatar' ).attr( 'src', playerInfo.avatar );
  });
};


var select = function( id, name ) {
  if( !(selectedPkg = pkgs[id]) ) {
    return
  }
  $( '#target-name'   ).text( name               );
  $( '#target-fee'    ).text( selectedPkg.amount );
  $( '#target-amount' ).text( selectedPkg.cnt    );
  $( '#target-area'   ).show();

  window.scrollTo( 200, 200 );
}

/**
 * 生成H5支付链接
 **/
var createPayURL = function( uid, area, pkgId, title, cb ) {
  var json = {
    "playerId": uid,
    "area":     area,
    "cid":      pkgId,
    "title":    title,
    "ipaddr":   ipaddr
  }
  $.post( baseUrl + "/h5/create", { orderInfo: JSON.stringify( json ) }, function( result ) {
  	if( cb ) {
      cb( result );
    }
  }, 'json' );
}

$(function() {
  var area    = getParam( 'app' );
  var appInfo = appInfos[area];
  if( appInfo ) {
    $( 'title'      ).text( appInfo['title'] );
    $( '#app-title' ).text( appInfo['name']  );
  }

  // 抓取用户信息
  fetchPlayerInfoById();
  ipaddr = returnCitySN.cip
  
  //抓取套餐信息
  getH5Pkgs( getParam( 'app' ), function( d ) {
    if( !d ) {
      return;
    }
    pkgs = JSON.parse( d )
    var img = '/static/img/icon-diamond.jpg';
    // var img = '/static/img/icon-' + getParam( 'app' ) + '.png';
    for( var idx = 0; idx < pkgs.length; idx++ ) {
      var desc = pkgs[idx].amount + '元=' + pkgs[idx].cnt + '颗钻石';
      var dom  = '<a href="#target-area" class="cart-item item" onclick="select(' + idx +', \''+ (pkgs[idx].cnt + '颗钻石') +'\')">' +
                 '<img src="' + img +'">' +
                 '<h1>' + (pkgs[idx].cnt + '颗钻石') + '</h1>' +
                 '<h2>选择</h2>' +
                 '<h3>' + desc + '</h3>' +
                 '</a>';
      $( '#pkg-list' ).append( dom );
    }

    var pid = getParam( 'pid' )
    if( pid ) {
      select( pid, pkgs[pid].name )
    }
  });

  $( '#paybtn' ).click( function() {
    var uid  = getParam( 'uid' );
    var area = getParam( 'app' );
    if( uid <= 0 ) {
      alert( '参数错误!' );
      return;
    }
    if( !area ) {
      alert( '参数错误!' );
      return;
    }
    if( !appInfos[area] ) {
      alert( '参数错误!' );
      return;
    }
    createPayURL( uid, area, selectedPkg.id, appInfos[area]['title'], function( d ) {
      if( !d ) {
        return;
      }
      if( d.return_code != 0 ) {
        alert( '支付失败!' );
        return;
      } else {
        var cburl = encodeURIComponent( 'http://127.0.0.1/pay/result.html' );
        window.location.href = d.data + '&redirect_url=' + cburl;
      }
    });
  });
});
