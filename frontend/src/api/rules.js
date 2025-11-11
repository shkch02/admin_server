import client from './client'

export const getRules = async () => {
  const response = await client.get('/rules')
  return response.data
}

export const updateRules = async (ruleSet) => {
  const response = await client.put('/rules', ruleSet)
  return response.data
}

