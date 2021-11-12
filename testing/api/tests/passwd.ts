import TestUser from "../modules/user";
import getEmails from "../modules/mail-check";
import {sleep, post_direct} from "../modules/utils";
import faker from "faker"

async function signupRandomUser()
{
    const u = new TestUser()
    if(!await u.signupRandomUser())
    {
        fail('failed signing up random user')
    }
    return u
}

async function requestPasswordReset(email:string)
{
    const resp = await post_direct("/reset/init", 
        {method:"post", data:{email: email} })
    console.log("reset passwd response ", resp)
    expect(resp).toBeDefined()
    expect(resp.status).toBe("ok")
}

async function getPasswordResetTokenFromEmail(email:string)
{
    const emails = await getEmails(email)
    if(!emails || emails.length < 2)
    {
        fail("Didn't get the password reset email ")
        return
    }
    expect(emails.items[0].Content.Body).toBeDefined()
    //console.log("Email 0\n", emails.items[0].Content.Body)
   // console.log("Email 1\n", emails.items[1].Content.Body)
    const content = emails.items[0].Content.Body;
    
    let rx = /reset\?token=3D([\w]+)/g;
   
    let matches = rx.exec(content);
    if(!matches || matches.length < 2)
    {
        fail("Email does not have confirmation code ")
        return
    }
    //console.log("reset password code ", matches[1]);
    return matches[1]
}

async function updatePassword(token:string, new_pwd:string)
{
    const resp2 = await post_direct("/reset/update", 
    {
    method:"post", 
    data:{
        token: token,
        password: new_pwd
    } })
    expect(resp2).toBeDefined()

    expect(resp2.status).toBe("ok") 
}

test('PW101: Reset Password', async function () {
    
    const u = await signupRandomUser()
    await sleep(2000);
    
    await requestPasswordReset(u.email)
    
    await sleep(2000);
    
    const token = await getPasswordResetTokenFromEmail(u.email)
    
    if(!token){ fail("Failed getting token from email ") }
    
    const new_pwd = faker.random.alphaNumeric(22)
    
    await updatePassword(token, new_pwd)
    
    const login_resp = await u.login()
    //console.log("login using old password ", login_resp)
    if(login_resp != false)
    {
        fail("Could login using old password after reset password")
    }
    u.updatePassword(new_pwd)
    
    const login2 = await u.login()
    if(!login2)
    {
        fail("failed logging in after reset password ")
    }
    //console.log("login response after reset password ", login2)
})

test('PW102: Shouldnt send Reset Password if the user has not signed up', async function () {
  
    const email = faker.internet.email()
    let throws_error = false
    try{
        await requestPasswordReset(email)
    }catch(e)
    {
        throws_error = true
        console.log("expected error got from reset request ", e.message)
    }
    //reset password without non existing email should throw error
    expect(throws_error).toBe(true)
    
})

test('PW103: should be able to login using old password of reset link is not clicked', async function () {
    const u = await signupRandomUser()
    await sleep(2000);
    
    await requestPasswordReset(u.email)
    
    await sleep(2000);  
    
    const login2 = await u.login()
    expect(login2).toBe(true)
    
    //Note: this is just do any action that can be done only after logging in 
    let respKey = await u.createApiKey()
    console.log("api key created ", respKey)
    if(!respKey || respKey.length <= 6)
    {
        fail("Expected to login and do some action ")
    }
})

test('PW104: reset password with bad token', async function () {
    const u = await signupRandomUser()
    
    await sleep(2000);
    
    await requestPasswordReset(u.email)
    
    await sleep(2000);

    let token = await getPasswordResetTokenFromEmail(u.email)
    
    let throws_error = false
    try{
        await updatePassword(token+"x", faker.random.alphaNumeric(8))
    }catch(e)
    {
        throws_error = true
        console.log("Reset password with bad code throws error ", e.message)
    }
    expect(throws_error).toBe(true)
   
    expect(await u.login()).toBe(true)
})

test('PW105: reset password with empty token', async function () {
    const u = await signupRandomUser()
    
    await sleep(2000);
    
    await requestPasswordReset(u.email)
    
    await sleep(2000);

    let throws_error = false
    try{
        await updatePassword("", faker.random.alphaNumeric(8))
    }catch(e)
    {
        throws_error = true
        console.log("Reset password with bad code throws error ", e.message)
    }
    expect(throws_error).toBe(true)
   
    expect(await u.login()).toBe(true)
})


test('PW106: reset password with empty password', async function () {
    console.log("tc PW106 reset password ")
    const u = await signupRandomUser()
    
    await sleep(2000);
    
    await requestPasswordReset(u.email)
    
    await sleep(2000);
    
    let token = await getPasswordResetTokenFromEmail(u.email)
    if(!token){ fail("password reset token is empty ") }
    
    let throws_error = false
    try{
        const resp = await updatePassword(token, "")
        console.log(" reset passwd res ", resp)
    }catch(e)
    {
        throws_error = true
        console.log("Reset password with bad code throws error ", e.message)
    }
    expect(throws_error).toBe(true)
    expect(await u.login()).toBe(true)
})


