import httpClient from '../utils/httpClient'

export default {
  loginLogList (query, callback, onError) {
    httpClient.get('/system/login-log', query, callback, onError)
  },
  loginLogUsers (query, callback, onError) {
    httpClient.get('/system/login-log/users', query, callback, onError)
  },
  stats (query, callback, onError) {
    httpClient.get('/system/stats', query, callback, onError)
  }
}
