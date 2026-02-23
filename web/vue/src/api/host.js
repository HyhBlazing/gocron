import httpClient from '../utils/httpClient'

export default {
  list (query, callback, onError) {
    httpClient.get('/host', query, callback, onError)
  },

  all (query, callback) {
    httpClient.get('/host/all', {}, callback)
  },

  detail (id, callback) {
    httpClient.get(`/host/${id}`, {}, callback)
  },

  update (data, callback) {
    httpClient.post('/host/store', data, callback)
  },

  remove (id, callback) {
    httpClient.post(`/host/remove/${id}`, {}, callback)
  },

  ping (id, callback) {
    httpClient.get(`/host/ping/${id}`, {}, callback)
  }
}
