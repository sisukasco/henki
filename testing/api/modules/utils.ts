import Axios from "axios";

export function sleep(ms:number) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

export async function post_direct(path:string, options:any){
    const res = await Axios(global.test_config.endpoint+path, options)
    return res.data
}

export async function get_direct(path:string, options:any){
    const res = await Axios(global.test_config.endpoint+path, { method:"get", ...options})
    return res.data
}