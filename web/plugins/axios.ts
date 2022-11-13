import axios from 'axios'

const client = axios.create({ baseURL: process.env.BASE_URL_BACKEND })

export function setAuthToken(token: string | null) {
  client.defaults.headers.Authorization = 'Bearer ' + token
}

export default client
