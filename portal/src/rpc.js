import axios from 'axios'

export const RPC = {
    doRPC: async function(api, req) {
        try {
            req.sid = window.localStorage.getItem('sid');
            let data = await axios.post(`/rpc/class/${api}`, req);
            return data.data;
        } catch (e) {
            return null;
        }
    },
    setSession(sid) {
        localStorage.setItem('sid', sid);
    }
}