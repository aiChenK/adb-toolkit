import axios from 'axios'
import { message } from 'ant-design-vue'

axios.defaults.timeout = 25000

export function GET(url, params, option = {}, config = {}) {
  if (params) {
    config.params = params
  }
  return fetch(url, 'get', option, config)
}

export function POST(url, params, option = {}, config = {}) {
  if (params) {
    config.data = params
  }
  return fetch(url, 'post', option, config)
}

export function DELETE(url, params, option = {}, config = {}) {
  if (params) {
    config.params = params
  }
  return fetch(url, 'delete', option, config)
}

export function PATCH(url, params, option = {}, config = {}) {
  if (params) {
    config.data = params
  }
  return fetch(url, 'patch', option, config)
}

// 统一请求封装
function fetch(url, method, option, config) {
  config.url = url
  config.method = method || 'get'
  return axios(config)
    .then(res => {
      if (typeof res === 'string') {
        return Promise.reject(new Error(`接口 ${url} 未返回正确格式`))
      }
      if (option.successMsg) {
        message.success(typeof option.successMsg === 'string' ? option.successMsg : '操作成功')
      }
      return Promise.resolve(res)
    })
    .catch(err => {
      const errorMsg =
        option.errorMsg ||
        (err.response && err.response.data && (err.response.data.errMessage || err.response.data.error)) ||
        err.message ||
        err.statusText ||
        '网络请求失败'

      if (option.errorMsg !== false) {
        message.error(errorMsg)
      }
      return Promise.reject(errorMsg)
    })
}