//logs.js
var app = getApp()
Page({
  data: {
    // joins: [
    //   {rota_id: "291255583271555074", title: "一月份值班表"},
    //   {rota_id: "291255583271555075", title: "人事部值班表"},
    //   {rota_id: "291255583271555076", title: "上半年值班表"},
    // ]
    joins: {} // 数组
  },

  onLoad: function () {
    var that = this;
    wx.request({
      url: app.globalData.host + '/join',
      header: {
        'content-type': 'application/json',
        'cookie': wx.getStorageSync('cookie')
      },
      method: 'GET',
      dataType: 'json',
      success: function(res) {
        if(res.data.status != 0) {
          wx.showModal({
            title: '错误!',
            content: '获取失败',
            showCancel: false
          })
        }
        that.setData({
          joins: res.data.joins
        })
      }
    })
  },
  handleClick1: function(event) {
    var rotaId = event.currentTarget.dataset.rotaid;
    wx.navigateTo({
      url: '../rota/rota?rotaId=' + rotaId
    })
  },
  
})
