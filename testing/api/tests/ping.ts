
import Axios from "axios";


test('PG101: Ping test', async function () {

    const res = await Axios.get(global.test_config.endpoint+"/ping")
    console.log("ping result ",res.data)
    expect(res.data).toBe("pong")
})