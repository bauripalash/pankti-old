use std::ffi::{CStr , CString};
use std::os::raw::c_char;
use std::env;
use std::fs;

extern "C"{
    fn DoParse(src : GoString) -> *const c_char;
}

#[repr(C)]
struct GoString{
    a: *const c_char,
    b: i64
}

fn run_file(src : String) {
    //let src = "show(123);";
    let c_str = CString::new(src).expect("failed in CString::new");
    let ptr = c_str.as_ptr();
    let go_string = GoString{
        a : ptr,
        b: c_str.as_bytes().len() as i64
    };

    let result = unsafe {
        DoParse(go_string)
    };

    let c_res = unsafe {
        CStr::from_ptr(result)
    };

    let res = c_res.to_str().expect("Failed to parse result as native string");
    match res.is_empty() {
        true => println!("Failed program"),
        false => println!("{}" , res)
    }
}

fn main(){
    let args : Vec<String> = env::args().collect();
    let file_path = &args[1];
    let source = fs::read_to_string(file_path)
        .expect("Failed to read the file");

    run_file(source);

}
