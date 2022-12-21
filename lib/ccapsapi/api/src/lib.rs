extern crate libc;

use std::ffi::{CStr, CString};
use std::str::FromStr;

pub mod utils;
pub mod api;
pub mod constants;


#[no_mangle]
pub extern "C" fn dynamic_decode_storage(pallet_name: *const libc::c_char, storage_entry: *const libc::c_char, raw: *const libc::c_char, metadata: *const libc::c_char) -> *const libc::c_char {
    let result = api::inner_dynamic_decode_storage(pallet_name, storage_entry, raw, metadata);
    return match result {
        Ok(res) => {
            utils::success(res)
        }
        Err(err) => {
            utils::fail(err.to_string())
        }
    };
}


#[no_mangle]
pub extern "C" fn decode_extrinsic(raw: *const libc::c_char) -> *const libc::c_char {
    let result = api::inner_decode_extrinsic(raw);
    return match result {
        Ok(res) => {
            utils::success(res)
        }
        Err(err) => {
            utils::fail(err.to_string())
        }
    };
}


#[no_mangle]
pub extern "C" fn pallet_info(data: *const libc::c_char, pallet_name: *const libc::c_char, call_name: *const libc::c_char) -> *const libc::c_char {
    let result = api::inner_pallet_info(data, pallet_name, call_name);
    return match result {
        Ok(res) => {
            utils::success(res)
        }
        Err(err) => {
            utils::fail(err.to_string())
        }
    };
}


#[no_mangle]
pub extern "C" fn signed_extrinsic(hash: *const libc::c_char, seed: *const libc::c_char, to: *const libc::c_char, amount: *const libc::c_char, nonce: *const libc::c_char, spec_version: *const libc::c_char, transaction_version: *const libc::c_char) -> *const libc::c_char {
    let result = api::inner_signed_extrinsic(hash, seed, to, amount, nonce, spec_version, transaction_version);
    return match result {
        Ok(res) => {
            utils::success(res)
        }
        Err(err) => {
            utils::fail(err.to_string())
        }
    };
}


#[no_mangle]
pub unsafe extern "C" fn free(value: *mut libc::c_char) {
    CString::from_raw(value);
}
