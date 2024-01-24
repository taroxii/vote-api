
## Description

This is an example of implementation of Clean Architecture in Go (Golang) projects.

Rule of Clean Architecture by Uncle Bob

- Independent of Frameworks. The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
- Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
- Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
- Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
- Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.

More at https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html

This project has 4 Domain layer :

- Models Layer
- Repository Layer
- Usecase Layer
- Interface Layer

## migrate

__Database migrations written in Go. Use as [CLI](#cli-usage) or import as [library](#use-in-your-go-project).__

* Migrate reads migrations from [sources](#migration-sources)
   and applies them in correct order to a [database](#databases).
* Drivers are "dumb", migrate glues everything together and makes sure the logic is bulletproof.
   (Keeps the drivers lightweight, too.)
* Database drivers don't assume things or try to correct user input. When in doubt, fail.
### Basic usage

```bash
$ migrate -source db/migrations -database postgres://localhost:5432/database up 
```

### Load test result

  execution: local
     script: scripts/create-item.js
     output: json (create-item-test.json)

  scenarios: (100.00%) 1 scenario, 20 max VUs, 10m30s max duration (incl. graceful stop):
           * reserve: 100 iterations for each of 20 VUs (maxDuration: 10m0s, gracefulStop: 30s)


     ✓ is status 200
     ✗ is status 400
      ↳  0% — ✓ 0 / ✗ 2000
     ✗ is status 500
      ↳  0% — ✓ 0 / ✗ 2000

     checks.........................: 33.33% ✓ 2000       ✗ 4000
     data_received..................: 745 kB 169 kB/s
     data_sent......................: 1.1 MB 257 kB/s
     http_req_blocked...............: avg=26.3µs  min=1µs    med=3µs     max=2.96ms   p(90)=4µs     p(95)=4µs    
     http_req_connecting............: avg=15µs    min=0s     med=0s      max=1.81ms   p(90)=0s      p(95)=0s     
     http_req_duration..............: avg=41.86ms min=1.18ms med=38.13ms max=297.41ms p(90)=69.24ms p(95)=82.72ms
       { expected_response:true }...: avg=41.86ms min=1.18ms med=38.13ms max=297.41ms p(90)=69.24ms p(95)=82.72ms
     http_req_failed................: 0.00%  ✓ 0          ✗ 2000
     http_req_receiving.............: avg=46.6µs  min=19µs   med=40µs    max=499µs    p(90)=65µs    p(95)=73µs   
     http_req_sending...............: avg=22.97µs min=10µs   med=21µs    max=592µs    p(90)=28µs    p(95)=31µs   
     http_req_tls_handshaking.......: avg=0s      min=0s     med=0s      max=0s       p(90)=0s      p(95)=0s     
     http_req_waiting...............: avg=41.79ms min=1.12ms med=38.06ms max=297.33ms p(90)=69.18ms p(95)=82.66ms
     http_reqs......................: 2000   453.087645/s
     iteration_duration.............: avg=42.18ms min=1.42ms med=38.52ms max=297.63ms p(90)=69.52ms p(95)=82.94ms
     iterations.....................: 2000   453.087645/s
     vus............................: 19     min=19       max=20
     vus_max........................: 20     min=20       max=20