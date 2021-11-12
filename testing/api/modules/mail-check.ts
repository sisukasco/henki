import Axios from "axios";


declare global {
    namespace NodeJS {
      interface Global {
        test_config: any
      }
    }
}


export default async function getEmails(to:string){
    const res = await Axios(global.test_config.emails.endpoint+"/api/v2/search", {
        method: "get",
        params:{
            kind:"to",
            query: to
        }
    })
    return res.data
}