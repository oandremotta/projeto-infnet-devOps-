import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '1m', target: 50 },
        { duration: '2m', target: 150 },
        { duration: '2m', target: 200 },
        { duration: '2m', target: 0 },
    ],
};

export default function () {

    const url = 'http://localhost:8000/healthz';


    const res = http.get(url);


    check(res, {
        'status is 200': (r) => r.status === 200,
    });


    check(res, {
        'response time is under 500ms': (r) => r.timings.duration < 500,
    });


    check(res, {
        'response body contains "ok"': (r) => r.body.includes('ok'),
    });


    check(res, {
        'content-type is application/json': (r) => r.headers['Content-Type'] === 'application/json',
    });


    sleep(1);
}
