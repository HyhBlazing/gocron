import httpClient from '../utils/httpClient'

export default {
  // 任务列表
  list (query, callback, onError) {
    httpClient.batchGet([
      {
        uri: '/task',
        params: query
      },
      {
        uri: '/host/all'
      }
    ], callback, onError)
  },

  detail (id, callback) {
    httpClient.batchGet([
      {
        uri: `/task/${id}`
      },
      {
        uri: '/host/all'
      }
    ], callback)
  },

  update (data, callback) {
    httpClient.post('/task/store', data, callback)
  },

  remove (id, callback) {
    httpClient.post(`/task/remove/${id}`, {}, callback)
  },

  enable (id, callback) {
    httpClient.post(`/task/enable/${id}`, {}, callback)
  },

  disable (id, callback) {
    httpClient.post(`/task/disable/${id}`, {}, callback)
  },

  run (id, callback) {
    httpClient.get(`/task/run/${id}`, {}, callback)
  },
  tags (query, callback, onError) {
    httpClient.get('/task/tags', query || {}, callback, onError)
  },
  nextRuns (spec, callback, onError) {
    httpClient.get('/task/next-runs', { spec }, callback, onError)
  }
}
