//app.js
const util = require('./utils/util.js')
const host = '...'
App({
  onLaunch: function () {
    // 展示本地存储能力
    // var logs = wx.getStorageSync('logs') || []
    // logs.unshift(Date.now())
    // wx.setStorageSync('logs', logs)

    // 检查登录态是否过期
    wx.checkSession({
      success() {
        //session_key 未过期，并且在本生命周期一直有效
        const value = wx.getStorageSync('cookie');
        util.login()
      },
      fail () {
        // session_key 已经失效，需要重新执行登录流程
        console.log('session_key 已经失效')
        util.login() //重新登录
      }
    })
    
    // 获取用户信息
    wx.getSetting({
      success: res => {
        if (res.authSetting['scope.userInfo']) {
          // 已经授权，可以直接调用 getUserInfo 获取头像昵称，不会弹框
          wx.getUserInfo({
            success: res => {
              this.globalData.userInfo = res.userInfo
    
              // 由于 getUserInfo 是网络请求，可能会在 Page.onLoad 之后才返回
              // 所以此处加入 callback 以防止这种情况
              if (this.userInfoReadyCallback) {
                this.userInfoReadyCallback(res)
              }
              // console.log('获取用户信息后,昵称为:' + this.globalData.userInfo.nickName)
              this.storeuserInfo()
            }
          })
        }
      }
    })
  },
  //将登陆的用户详情信息存储到数据库
  storeuserInfo: function () {
    let userInfo = this.globalData.userInfo
    wx.request({
      url: host + '/savePerson',
      data: {
        nick_name: userInfo.nickName,
      },
      header: {
        'content-type': 'application/json',
        'cookie': wx.getStorageSync('cookie') 
      },
      method: 'POST',
      dataType: 'json',
      success: function(res) {
        if(res.data.status != 0) {
          wx.showModal({
            title: '警告!',
            content: res.data.msg,
            showCancel: false
          })
        }
      },
      fail: function(res) {
        wx.showModal({
          title: '警告!',
          content: '昵称存储出错!',
          showCancel: false
        })
      }
    });
  },
  globalData: {
    host,
    userInfo: null,
  }
})