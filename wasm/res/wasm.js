'use strict';

const WASM_URL = '/wasm.wasm';

var wasm;

function runSource(){
	var src = document.getElementById("src_input").value;
	document.getElementById("result").value = runner(src);
}

function init() {
document.getElementById("runbtn").onclick = runSource;
  const go = new Go();
  if ('instantiateStreaming' in WebAssembly) {
    WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then(function (obj) {
      wasm = obj.instance;
      go.run(wasm);
      runSource();
    })
  } else {
    fetch(WASM_URL).then(resp =>
      resp.arrayBuffer()
    ).then(bytes =>
      WebAssembly.instantiate(bytes, go.importObject).then(function (obj) {
        wasm = obj.instance;
        go.run(wasm);
        runSource();
      })
    )
  }
}

init();
