const app = getApp()
Page({
  data: {
    rotaId: '',
    colorArrays: [ "#85B8CF", "#90C652", "#D8AA5A", "#FC9F9D", "#0A9A84", "#61BC69", "#12AEF3", "#E29AAD"],
    interval: []
  },
  onLoad: function (options) {
    let that = this;
    that.data.rotaId = options.rotaId
    var rotaId = options.rotaId
    wx.request({
      url: app.globalData.host + '/generate/' + rotaId,
      header: {
        'content-type': 'application/json'
      },
      method: 'GET',
      dataType: 'json',
      success: function(res) {
        if(res.data.status == 1) {
          wx.showModal({
            title: '错误!',
            content: '生成失败',
            showCancel: false
          })
        } else if(res.data.status == 2) {
          wx.showModal({
            title: '哎呀',
            content: '值班表还没人填写',
            showCancel: false
          })
        } else if(res.data.status == 0) {
          wx.showModal({
            title: '恭喜',
            content: '生成成功',
            showCancel: false
          })
          that.setData({
            interval: res.data.interval
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
  },
  handleClick1: function() {
    var that = this;
    wx.downloadFile({
      url: app.globalData.host + '/download/' + that.data.rotaId,
      success: function (res) {
        // if(res.data.status == 1) {
        //   wx.showModal({
        //     title: '错误!',
        //     content: res.data.msg,
        //     showCancel: false
        //   })
        // } else {
        //  const tempFilePath = res.tempFilePath
        wx.saveFile({
          tempFilePath: res.tempFilePath,
          success (res) {
            const savedFilePath = res.savedFilePath
            wx.openDocument({
              filePath: savedFilePath,
              fileType: 'csv',
              success: function (res) {
                console.log('成功打开文档')
              },
            });
          },
          fail (res) {
            wx.showModal({
              title: '哎呀',
              content: "保存失败,服务器出错",
              showCancel: false
            })
          }
        })
          
        // }
      },
      fail: function (res) {
        wx.showModal({
          title: '哎呀',
          content: "导出失败,服务器出错",
          showCancel: false
        })
      }
    })
  },
})
