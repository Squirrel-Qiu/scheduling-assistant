//logs.js
var app = getApp()
Page({
  data: {
    // rotas: [
    //   {rota_id: "291255583271555078", title: "一月份值班表"},
    //   {rota_id: "291255583271555079", title: "人事部值班表"},
    // ]
    rotas: {} // 数组
  },

  onLoad: function () {
    var that = this;
    wx.request({
      url: app.globalData.host + '/rotas',
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
          rotas: res.data.rotas
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
  handleClick2: function(event) {
    var rotaId = event.currentTarget.dataset.rotaid;
    // console.log('当前标题对应的rotaId:' + rotaId);
    wx.navigateTo({
      url: '../table/table?rotaId=' + rotaId
    })
  },
  handleClick3: function(event) {
    let that = this;
    var rotaId = event.currentTarget.dataset.rotaid;
    wx.showModal({
      title: '警告!',
      content: '是否确认删除',
      success (res) {
        if (res.confirm) {
          // console.log('用户点击确定')
          wx.request({
            url: app.globalData.host + '/delete/' + rotaId,
            header: {
              'content-type': 'application/json',
              'cookie': wx.getStorageSync('cookie')
            },
            method: 'DELETE',
            dataType: 'json',
            success: function(res) {
              if(res.data.status == 1) {
                wx.showModal({
                  title: '警告!',
                  content: '删除失败',
                  showCancel: false
                })
              } else if(res.data.status == 0) {
                wx.showModal({
                  title: '恭喜~',
                  content: '删除成功',
                  showCancel: false
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
          })
        } else if (res.cancel) {
          // console.log('用户点击取消')
        }
      }
    })
  }
})
