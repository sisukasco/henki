import TestUser from "../modules/user";
import faker from "faker"
import qs from 'qs';
import getEmails from "../modules/mail-check";
import {sleep, post_direct, get_direct} from "../modules/utils";


test('SU101: Raw Signup should work fine', async function () {
    const signup = {
        email: faker.internet.email(),
        password: faker.random.alphaNumeric(20)
    };
    
    const resp = await post_direct("/signup", {method:"post", data:signup })
    console.log("Signup response ", resp)
    expect(resp).toBeDefined()
    
    await sleep(2000);
    
    const emails = await getEmails(signup.email)
    console.log(" signup response email: ", emails.items[0].Content)
    expect(emails.items.length).toBe(1)
    
    
    const token = await post_direct("/token", {
        method:"post", 
        data:qs.stringify({
            grant_type: "password",
            username: signup.email,
            password: signup.password
        }),
        headers: {
            'content-type': 'application/x-www-form-urlencoded;charset=utf-8'
            } 
        })
    expect(token).toBeDefined()
    expect(token.access_token).toBeDefined()
    console.log("received token ", token)
})
    
test('SU102: Signup should work fine', async function () {
    /*let auth= new Chimes( "http://local.dockform.com:3121", 
                "local.dockform.com") 
    let res = await auth.signup("somename2@ajshdkjash.com", "alonglongpwd",{})
    expect(res.status).toBe("ok")
    */
   let u = new TestUser()
   if(await u.signupRandomUser())
   {
       console.log("signed up random user ", u.email)
   }else{
    fail('failed signing up random user');
   }
   
   if(!await u.login())
   {
        fail('failed logging in as random user');
   }
   console.log("Logged in as user ", u.email)
   
});

test('SU103: Confirm Email After Signup', async function () {
   let u = new TestUser()
   if(!await u.signupRandomUser())
   {
        fail('failed signing up random user')
   }
   
   await sleep(2000);
    
   const emails = await getEmails(u.email)
   if(!emails || emails.length < 1)
   {
       fail("Didn't get the confirmation email ")
       return
   }
   const content = emails.items[0].Content.Body;
   console.log(" signup response email: ", content)
   //emails.items[0].Content
   let rx = /confirm\?code=3D([\w]+)/g;
   
   let matches = rx.exec(content);
   if(!matches || matches.length < 2)
   {
       fail("Email does not have confirmation code ")
       return
   }
   console.log("confirmation code ", matches[1]);
   const code = matches[1]
   
   const resp = await post_direct("/confirm", {method:"post", data:{code} })
   expect(resp).toBeDefined()
   console.log(" confirmation response ", resp )
   
   await u.login()
   
   if(!u.isEmailConfirmed())
   {
       fail("Email confirmation does not update the email confirmation status")
   }
})

test('SU104: Renew Refresh Token Should work fine ', async function () {
    const signup = {
        email: faker.internet.email(),
        password: faker.random.alphaNumeric(20)
    };
    
    const resp = await post_direct("/signup", {method:"post", data:signup })
    //console.log("Signup response ", resp)
    expect(resp).toBeDefined()
    
    await sleep(2000);
    
    const token = await post_direct("/token", {
        method:"post", 
        data:qs.stringify({
            grant_type: "password",
            username: signup.email,
            password: signup.password
        }),
        headers: {
            'content-type': 'application/x-www-form-urlencoded;charset=utf-8'
            } 
        })
    expect(token).toBeDefined()
    expect(token.access_token).toBeDefined()
    expect(token.refresh_token).toBeDefined()
    //console.log("received token ", token)
    
    //console.log("renewing the refresh token... ")
    
    const token2 = await post_direct("/token", {
        method:"post", 
        data:qs.stringify({
            grant_type: "refresh_token",
            refresh_token: token.refresh_token
        }),
        headers: {
            'content-type': 'application/x-www-form-urlencoded;charset=utf-8'
            } 
        })
    //console.log("refresh token received ", token2)
    expect(token2).toBeDefined()
    expect(token2.access_token).toBeDefined()
    expect(token2.access_token.length).toBeGreaterThan(12)
    expect(token2.refresh_token).toBeDefined()
    expect(token2.refresh_token.length).toBeGreaterThan(8)
    
    const user_info = await get_direct("/user",{headers:{Authorization: `Bearer ${token2.access_token}` }})
    console.log("Got user info ", user_info)
    expect(user_info).toBeDefined()
    expect(user_info.id).toBeDefined()
    expect(user_info.id.length).toBeGreaterThan(8)
})