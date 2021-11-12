

import TestUser from "../modules/user";

test('AK101: Create API Key', async function () {
    
   let u = new TestUser()
   await u.signupRandomUser()
   
   await u.login()
   
   let key = await u.createApiKey()
   
   expect(key.length).toBeGreaterThan(8)
   
   const keys = await u.getAPIKeys()
   console.log("Get api keys response ", keys)
   expect(keys.length).toBe(1)
   expect(keys[0].key).toBe(key)
   
   let key2 = await u.createApiKey()
   const keys2 = await u.getAPIKeys()
   console.log("Get api keys response(2) ", keys2)
   expect(keys2.length).toBe(2)
   expect(keys2[1].key).toBe(key2)
   
});

test('AK102: Delete API Key', async function () {
    
    let u = new TestUser()
    await u.signupRandomUser()
    
    await u.login()
    
    let key = await u.createApiKey()
    expect(key.length).toBeGreaterThan(8)
    
    const keys = await u.getAPIKeys()
    expect(keys.length).toBe(1)
    
    const resp = await u.deleteAPIKey(key)
    expect(resp).toBeDefined()
    console.log("api key delete response ", resp)
    
    const keys2 = await u.getAPIKeys()
    console.log(" get keys (2) ", keys2)
    expect(keys2.length).toBe(0)
})