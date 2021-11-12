import faker from "faker"
import Chimes,{AuthedConnection} from "@sisukas/chimes";
declare global {
    namespace NodeJS {
      interface Global {
        test_config: any
      }
    }
}

export default class TestUser
{
    public email:string = ""
    private password:string = ""
    private auth:Chimes|null=null
    constructor()
    {
        
    }
    updatePassword(pwd:string)
    {
        this.password = pwd
    }
    async signupRandomUser()
    {
        console.log("endpoint ", global.test_config.endpoint)
        
        this.email = faker.internet.email()
        this.password= faker.random.alphaNumeric(16)
        let auth= new Chimes( global.test_config.endpoint, 
            global.test_config.endpoint.aud) 
        try{
            let res = await auth.signup(this.email,this.password)    
            if(res.status != "ok")
            {
                return false
            }
        }catch(e)
        {
            console.log("Error signing up user" , e)
            return false;    
        }
        return true
    }
    //signup()
    async login()
    {
        this.auth= new Chimes( global.test_config.endpoint, 
            global.test_config.endpoint.aud) 
            
        try{
            let res = await this.auth.login(this.email, this.password )
            console.log("login result ", res) 
            if(res.ok != true){
                
                return false
            }           
        }catch(e)
        {
            console.log("Caught error while logging in ",e )
            return false;
        }

        return true;
    }
    async createApiKey()
    {
        if(!this.auth || !this.auth.user)
        {
            throw new Error("The user is not loggd in!")
        }
        let res = await this.auth.user.request("/api-key", 
        {method:'post', data: {}})
        
        console.log("create api key respose ", res)
        
        return res.key
    }
    getAuthConnection():AuthedConnection
    {
        if(!this.auth){
            throw new Error("User not authenticated! ")
        }
        const ret = this.auth.getAuthConnection()
        if(!ret){
            throw new Error("User not authenticated! ")
        }
        return ret
    }
    async getForms()
    {
        if(!this.auth || !this.auth.user){ return []}
        let res = await this.auth.user.request("/v1.0/forms", 
        {method:'get', data: {}})
        console.log("user's forms ", res)
        return res
    }
    
    async doRequest(path:string, options:any){
        if(!this.auth || !this.auth.user)
        {
            throw new Error("Not logged in !");
        }
        
        return this.auth.user.request(path, options)
    }
    
    
    isEmailConfirmed()
    {
        return this.auth?.user?.info.email_confirmed
    }
    getAPIKeys()
    {
        return this.doRequest("/api-key", {method:"get"})
    }
    deleteAPIKey(key: string)
    {
        return this.doRequest("/api-key", {method:"delete", data:{key}})
    }
}