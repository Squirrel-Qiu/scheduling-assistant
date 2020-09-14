const app = getApp()
Page({
  data: {
    rotaId: '',
    frees: [],
    times: [
      {value: '0'},{value: '1'},{value: '2'},{value: '3'},{value: '4'},
      {value: '5'},{value: '6'},{value: '7'},{value: '8'},{value: '9'},
      {value: '10'},{value: '11'},{value: '12'},{value: '13'},{value: '14'},
      {value: '15'},{value: '16'},{value: '17'},{value: '18'},{value: '19'},
      {value: '20'},{value: '21'},{value: '22'},{value: '23'},{value: '24'},
      {value: '25'},{value: '26'},{value: '27'},{value: '28'},{value: '29'},
      {value: '30'},{value: '31'},{value: '32'},{value: '33'},{value: '34'}
    ]
  },
  onLoad: function (options) {
    wx.showShareMenu({
      // withShareTicket: true,
      menus: ['shareAppMessage']
    })
    let that = this;
    that.data.rotaId = options.rotaId
    var rotaId = options.rotaId
    wx.request({
      url: app.globalData.host + '/rota/' + rotaId,
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
        } else if(res.data.status == 0) {
          if(res.data.frees != null) {
            // int数组 --> 字符数组
            // console.log('getFree的结果为：' + res.data.frees)
            var free = res.data.frees //.map(function(item) {return item + ''});
            var times = that.data.times
            for (let i = 0, lenI = times.length; i < lenI; ++i) {
              times[i].checked = false

              for (let j = 0, lenJ = free.length; j < lenJ; ++j) {
                // Tips：times[i].value和e.detail.value都为字符串类型
                if (parseInt(times[i].value) === free[j]) {
                  times[i].checked = true
                  break
                }
              }
            }
            that.setData({
              times
            })
          }
          
        }
      },
      fail: function(res) {
        wx.showModal({
          title: '错误!',
          content: '获取失败',
          showCancel: false
        })
      }
    })
  },
  checkboxChange: function(e) {
    // console.log('checkbox发生change事件，携带value值为：', e.detail.value)
    var times = this.data.times
    var f = e.detail.value
    var arr = []
    
    // console.log('数组的值:' + new Array(e.detail.value))

    for (let i = 0, lenI = times.length; i < lenI; ++i) {
      times[i].checked = false

      for (let j = 0, lenJ = f.length; j < lenJ; ++j) {
        if (times[i].value === f[j]) {
          arr.push(parseInt(f[j]))
          times[i].checked = true
          break
        }
      }
    }
    
    this.setData({
      frees: arr,
      times
    });
    // console.log('绑定到frees:' + this.data.frees)
  },
  handleClick1: function() {
    var that = this;
    // console.log('要提交的frees为:' + that.data.frees)
    wx.request({
      url: app.globalData.host + '/chooseFree/' + that.data.rotaId,
      data: {
        frees: that.data.frees
      },
      header: {
        'content-type': 'application/json',
        'cookie': wx.getStorageSync('cookie')
      },
      method: 'POST',
      dataType: 'json',
      success: function(res) {
        if(res.data.status == 5) {
          wx.showModal({
            title: '错误!',
            content: res.data.msg,
            showCancel: false
          })
        } else if(res.data.status == 3) {
          wx.showModal({
            title: '错误!',
            content: res.data.msg,
            showCancel: false
          })
        } else if(res.data.status == 0) {
          wx.showModal({
            title: '恭喜!',
            content: "提交成功",
            showCancel: false
          })
          wx.switchTab({
            url: '../join/join'
          })
        } else if(res.data.status == 1 || res.data.status == 2) {
          wx.showModal({
            title: '警告!',
            content: "请求参数错误",
            showCancel: false
          })
        }
      },
      fail: function(res) {
        wx.showModal({
          title: '警告!',
          content: "提交失败,服务器出错",
          showCancel: false
        })
      }
    })
  },
  onShareAppMessage: function() {
    return {
      title: '排班助手',
      imageUrl: '../../images/capoo2.jpg'
    }
  }
})
