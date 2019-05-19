export const Util = {
    formatSize(sz) {
        const name = ["Byte", "KiB", "MiB", "GiB", "TiB", "PiB"];
        let i = 0;
        for (; i < name.length; i++) {
            if (sz >= 512) {
                sz /= 1024;
            } else {
                break;
            }
        }
        if (i >= 2) {
            sz = sz.toFixed(2);
        }
        return `${sz} ${name[i]}`;
    }
}