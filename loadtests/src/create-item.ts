import exec from 'k6/execution';
import { SharedArray } from 'k6/data';
import { Options } from 'k6/options';
import http from 'k6/http';
import { check } from 'k6';

// @ts-ignore
import file from 'k6/x/file';
function uuidv4(): string {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        let r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

const CONFIGS = {
    VUS: 20,
    API_PATH: 'http://localhost:8080/items',
    REQUEST_SAMPLE_SIZE: 100,
};

type CreateItemRequest = {
    correlation_id: string;
    name: string;
    description: string;
   
};

export const options: Options = {
    scenarios: {
        reserve: {
            executor: 'per-vu-iterations',
            vus: CONFIGS.VUS,            
            iterations: CONFIGS.REQUEST_SAMPLE_SIZE,
            maxDuration: '600s',
        }
    }
};


let data: Array<Omit<CreateItemRequest, 'correlation_id'>> = new SharedArray('data', () => {
    return JSON.parse(open('../data/feed-item.json'));
});
export default function main() {
    const request = data[exec.scenario.iterationInTest];
    const res = http.post(
        CONFIGS.API_PATH, 
        JSON.stringify({ ...request, correlation_id: uuidv4() } as CreateItemRequest), {
        headers: {
            'Content-Type': 'application/json',
            'Authorization':"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiaXNzIjoiQUNDT1VOVF9DRU5URVIiLCJzdWJqZWN0IjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwidXNlcm5hbWUiOiJnc3VtcHRvbjAifQ.3m-yhUU7Dm_VGQgRGvvDY0B3Rs1pre6fqTfN9ZA8lw0"
            
        }
    });
    if (res.status !== 200) {
        console.log(res.body);
    }
    check(res, {
        'is status 200': (r) => r.status === 200,
        'is status 400': (r) => r.status === 400,
        'is status 500': (r) => r.status === 500,
    });
}