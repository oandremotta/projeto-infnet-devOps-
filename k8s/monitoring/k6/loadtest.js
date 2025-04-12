import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
    vus: 100,
    duration: '2m',
};

export default function () {
    http.get('https://test.k6.io');
    sleep(1);
}