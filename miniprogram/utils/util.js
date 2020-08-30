var login = function() {
  // 通过这个命令获取 code, 这个 code 5 分钟有效, 只能用一次
  wx.login({
    success: res => {
      // console.log('js_code为:' + res.code);

      // 发送 res.code 到后台换取 cookie
      wx.request({
        url: '.../login',
        data: {
          appid: '...',
          secret: '...',
          js_code: res.code,
          grant_type: 'authorization_code'
        },
        header: {
          'content-type': 'application/json'
        },
        method: 'POST',
        success: function (res) {
          if(res.data.status != 0) {
            wx.showModal({
              title: '警告!',
              content: res.data.msg,
              showCancel: false
            })
          }
          // 请求成功后, 把后台返回的 cookie 缓存起来
          wx.setStorageSync('cookie', res.header['Set-Cookie'])
          // console.log('cookie为' + res.header['Set-Cookie'])
        },
        fail () {
          console.log('登录失败')
        }
      })
    }
  })
}

// 输出这个函数
module.exports = {
  login: login,
}
