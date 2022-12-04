const WASM_URL = '/wasm.wasm';

var wasm;


function runSource(){
	const res =  document.getElementById("result")

	const src = document.getElementById("src_input").value;
  //window.set_output(res)
	res.value+= evs(src);
}

(function() {

  var editor = CodeMirror.fromTextArea(document.getElementById('src_input') , {
    lineNumbers: true
  });
  editor.save();

  document.getElementById("runbtn").onclick = function(){
	  document.getElementById("result").innerHTML+= evs(document.getElementById("src_input").value);
  }

  document.getElementById("clearbtn").onclick = function(){
    document.getElementById("result").value = "";
  }



  let oldconsole = console.log;
  var logger = document.getElementById("result");
  console.log = function(msg){
    logger.value += msg + `
`;
  }
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("/wasm.wasm"), go.importObject).then((result) => {
    
    go.run(result.instance);
});
}());