"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.options = void 0;
const execution_1 = __importDefault(require("k6/execution"));
const data_1 = require("k6/data");
const http_1 = __importDefault(require("k6/http"));
const k6_1 = require("k6");
function uuidv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        let r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}
const CONFIGS = {
    VUS: 20,
    API_PATH: 'http://localhost:8080/items',
    REQUEST_SAMPLE_SIZE: 100,
};
exports.options = {
    scenarios: {
        reserve: {
            executor: 'per-vu-iterations',
            vus: CONFIGS.VUS,
            iterations: CONFIGS.REQUEST_SAMPLE_SIZE,
            maxDuration: '600s',
        }
    }
};
let data = new data_1.SharedArray('data', () => {
    return JSON.parse(open('../data/feed-item.json'));
});
function main() {
    const request = data[execution_1.default.scenario.iterationInTest];
    const res = http_1.default.post(CONFIGS.API_PATH, JSON.stringify(Object.assign(Object.assign({}, request), { correlation_id: uuidv4() })), {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiaXNzIjoiQUNDT1VOVF9DRU5URVIiLCJzdWJqZWN0IjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwidXNlcm5hbWUiOiJnc3VtcHRvbjAifQ.3m-yhUU7Dm_VGQgRGvvDY0B3Rs1pre6fqTfN9ZA8lw0"
        }
    });
    if (res.status !== 200) {
        console.log(res.body);
    }
    (0, k6_1.check)(res, {
        'is status 200': (r) => r.status === 200,
        'is status 400': (r) => r.status === 400,
        'is status 500': (r) => r.status === 500,
    });
}
exports.default = main;
