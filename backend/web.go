package main

import (
	"net/http"
	"strings"
)

func staticHandler() http.Handler {
	const index = `<!doctype html><html><head><meta charset="utf-8"><title>Observation runs</title></head><body><h1>Astronomy observation runs</h1><p id="health">Loading...</p><button id="refresh">Refresh runs</button><ul id="items"></ul><script src="/app.js"></script></body></html>`
	const app = `const list=document.querySelector('#items');async function load(){const xs=await (await fetch('/api/runs')).json();list.innerHTML=xs.map(x=>'<li>'+x.target+' / '+x.instrument+' - '+x.status+' <button data-id="'+x.id+'">Archive</button></li>').join('');document.querySelectorAll('[data-id]').forEach(b=>b.onclick=async()=>{await fetch('/api/runs/'+b.dataset.id+'/status',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({status:'archived'})});load()})}fetch('/healthz').then(r=>r.json()).then(x=>document.querySelector('#health').textContent=x.status+' / '+x.service);document.querySelector('#refresh').onclick=load;load();`
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			_, _ = w.Write([]byte(index))
			return
		}
		if r.URL.Path == "/app.js" {
			w.Header().Set("Content-Type", "application/javascript")
			_, _ = w.Write([]byte(app))
			return
		}
		if strings.HasPrefix(r.URL.Path, "/web/") {
			http.NotFound(w, r)
			return
		}
		http.NotFound(w, r)
	})
}
