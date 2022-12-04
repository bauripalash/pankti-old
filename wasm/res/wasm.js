const WASM_URL = '/wasm.wasm';

var wasm;


function runSource(){
	const res =  document.getElementById("result")

	const src = document.getElementById("src_input").value;
  //window.set_output(res)
	res.value+= evs(src);
}

(function() {
 let oldconsole = console.log;
  var logger = document.getElementById("result");
  console.log = function(msg){
    logger.value += msg + `
`;
  }



  document.getElementById("runbtn").onclick = function(){
	  console.log(evs(document.getElementById("src_input").value));
  }

  document.getElementById("clearbtn").onclick = function(){
    document.getElementById("result").value = "";
  }



   const go = new Go();
  WebAssembly.instantiateStreaming(fetch("/wasm.wasm"), go.importObject).then((result) => {
    
    go.run(result.instance);
});
}());
