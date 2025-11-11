import client from './client'

export const triggerTest = async (testType) => {
  const response = await client.post('/tests/trigger', { test_type: testType })
  return response.data
}

