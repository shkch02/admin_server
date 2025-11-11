import client from './client'

export const getAlerts = async (limit = 50, since = null) => {
  const params = { limit }
  if (since) {
    params.since = since
  }
  const response = await client.get('/alerts', { params })
  return response.data
}

