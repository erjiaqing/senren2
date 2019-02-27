import axios from 'axios'

export const RPC = {
    doRPC: async function (api, req) {
        // if cannot get pci_sid, try to get senren sid for problem ci group
        try {
            req.sid = window.localStorage.getItem('pci_sid');
            let data = await axios.post(`/rpc/pci/${api}`, req);
            return data.data;
        } catch (e) {
            return null;
        }
    },
    refreshSession: async function () {
        try {
            if (window.localStorage.getItem('pci_sid') && window.localStorage.getItem('pci_sid') !== "") {
                return;
            }
            let req = {};
            req.sid = window.localStorage.getItem('sid');
            let data = (await axios.post(`/rpc/pci/loginBySenrenSid`, req)).data;
            if (data.sid && data.sid !== "") {
                localStorage.setItem('pci_sid', data.sid);
            }
        } catch (e) {
            return;
        }
    },
    setSession(sid) {
        localStorage.setItem('pci_sid', sid);
    }
}