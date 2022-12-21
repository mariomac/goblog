import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 200,
  duration: '5m',
};


const urlGroups = [
    ['https://macias.info'],
    ['https://macias.info/entry/about.md'],
    ['https://macias.info/entry/202212020000_go_streams.md',
    'https://macias.info/static/assets/2022/streams/us.png',
    'https://macias.info/static/assets/2022/streams/bytes.png',
    'https://macias.info/static/assets/2022/streams/allocs.png'],
    ['https://macias.info/entry/202211040000_goblog_v1.md'],
    ['https://macias.info/entry/202201211818_https_support.md'],
    ['https://macias.info/entry/202109081800_k8s_informers.md'],
    ['https://macias.info/entry/202003151900_go_wasm_js.md',
    'https://macias.info/static/assets/2020/03/go_wasm/log_console.png',
    'https://macias.info/static/assets/2020/03/go_wasm/result.png'],
    ['https://macias.info/entry/201912201300_graal_aot.md',
    'https://macias.info/static/assets/2019/graal_aot/exec_time.png',
    'https://macias.info/static/assets/2019/graal_aot/max_rss.png']
];

export default function () {
    let idx = Math.floor(Math.random() * urlGroups.length);
    let urlGroup = urlGroups[idx];
    urlGroup.forEach(http.get);
}
