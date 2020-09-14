//index.js
//获取应用实例
var app = getApp()

Page({
  data: {
    current: 'homepage',
    userInfo: {},
    hasUserInfo: false,
    canIUse: wx.canIUse('button.open-type.getUserInfo'),
    title: '',
    shift: 1,
    limitChoose: 1,
    counter: 1
  },

  handleChange ({detail}) {
    this.setData({
      current: detail.key
    });
  },
  bindTitleInput: function (e) {
    this.setData({
      title: e.detail.value,
    });
  },
  handleChange1: function(e) {
    this.setData({
      shift: e.detail.value
    })
  },
  handleChange2: function(e) {
    this.setData({
      limitChoose: e.detail.value
    })
  },
  handleChange3: function(e) {
    this.setData({
      counter: e.detail.value
    })
  },
  onLoad: function () {
    wx.showShareMenu({
      withShareTicket: true,
      menus: ['shareAppMessage']
    })
    if (app.globalData.userInfo) {
      this.setData({
        userInfo: app.globalData.userInfo,
        hasUserInfo: true
      })
    } else if (this.data.canIUse){
      // 由于 getUserInfo 是网络请求，可能会在 Page.onLoad 之后才返回
      // 所以此处加入 callback 以防止这种情况
      app.userInfoReadyCallback = res => {
        this.setData({
          userInfo: res.userInfo,
          hasUserInfo: true
        })
      }
    } else {
      // 在没有 open-type=getUserInfo 版本的兼容处理
      wx.getUserInfo({
        success: res => {
          app.globalData.userInfo = res.userInfo
          this.setData({
            userInfo: res.userInfo,
            hasUserInfo: true
          })
        }
      })
    }
  },
  getUserInfo: function(e) {
    app.globalData.userInfo = e.detail.userInfo
    this.setData({
      userInfo: e.detail.userInfo,
      hasUserInfo: true
    })
    app.storeuserInfo()
  },
  handleClick: function() {
    var that = this;
    if (that.data.title == '') {
      wx.showModal({
        title: '警告!',
        content: '标题不能为空',
        showCancel: false
      })
    } else {
      wx.showLoading({
        title: '创建中...',
      })
      wx.request({
        url: app.globalData.host + '/newRota',
        data: {
          title: that.data.title,
          shift: that.data.shift,
          limit_choose: that.data.limitChoose, 
          counter: that.data.counter
        },
        header: {
          'content-type': 'application/json',
          'cookie': wx.getStorageSync('cookie')
        },
        method: 'POST',
        dataType: 'json',
        success: function (res) {
          wx.hideLoading();
          console.log(res)
          if(res.data.status == 1) {
            wx.showModal({
              title: '警告!',
              content: res.data.msg,
              showCancel: false
            })
          } else if(res.data.status == 2) {
            wx.showModal({
              title: '警告!',
              content: res.data.msg,
              showCancel: false
            })
          } else if(res.data.status == 3) {
            wx.showModal({
              title: '警告!',
              content: "创建失败,服务器出错",
              showCancel: false
            })
          } else if(res.data.status == 0) {
            wx.showModal({
              title: '恭喜!',
              content: '创建成功',
              showCancel: false,
            })
            console.log('创建成功跳转去:' + res.data.rota_id)
            wx.redirectTo({
              url: '../rota/rota?rotaId=' + res.data.rota_id
            })
          }
        },
        fail: function(res) {
          wx.showModal({
            title: '哎呀',
            content: '网络好像有点不对劲?',
            showCancel: false
          })
        }
      });
    }
  },
  onShareAppMessage: function() {
    return {
      title: '排班助手',
      imageUrl: '../../images/capoo2.jpg'
    }
  }
})
