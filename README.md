### Playing with Graphql and Go

#### Endpoints:
* Get one `curl -g 'http://localhost:3000/redirects?query={redirect(from:"/oldUrl"){from,to}}'`
* Get all `curl -g 'http://localhost:3000/redirects?query={redirects(offset:0,limit:30){from,to}}'`
* Create `curl -g 'http://localhost:3000/redirects?query=mutation+_{createRedirect(from:"/oldUrl",to:"/newUrl"){from}}'`
* Delete `curl -g 'http://localhost:3000/redirects?query=mutation+_{deleteRedirect(from:"/oldUrl"){from}}'`