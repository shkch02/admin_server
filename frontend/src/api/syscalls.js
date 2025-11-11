import client from './client'

export const getCallableSyscalls = async () => {
  const response = await client.get('/syscalls/callable')
  return response.data
}

