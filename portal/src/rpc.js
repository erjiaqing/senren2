import axios from 'axios'

export const RPC = {
    doRPC: async function(api, req) {
        try {
            let data = await axios.post(`/rpc/class/${api}`, req);
            return data.data;
        } catch (e) {
            return null;
        }
    }
}