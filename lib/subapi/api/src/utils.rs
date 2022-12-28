use std::ffi::{CStr, CString};
use anyhow::{Result};
use serde::{Deserialize, Serialize};


pub fn get_value(value: *const libc::c_char) -> Result<String> {
    let str_url = unsafe { CStr::from_ptr(value) }.to_str()?;
    return Ok(str_url.to_string());
}


pub fn twox_128(data: &[u8]) -> Result<String> {

    let vec = sp_core::twox_128(data).to_vec();
    let res = hex::encode(&vec);
    Ok(res)
}

pub fn blake2_128_concat(x: &[u8]) -> Result<String> {
    let vec = sp_core::blake2_128(x).iter().chain(x.iter()).cloned().collect::<Vec<_>>();
    let res = hex::encode(&vec);
    Ok(res)
}


pub fn twox_64_concat(x: &[u8]) -> Result<String> {
    let vec = sp_core::twox_64(x)
        .iter()
        .chain(x)
        .cloned()
        .collect::<Vec<_>>();
    let res = hex::encode(&vec);
    Ok(res)
}





#[derive(Serialize, Deserialize)]
struct Message {
    data: String,
    status: u8,
    msg: String,
}

pub fn success(data: String) -> *const ::libc::c_char {
    let msg = Message {
        data: data.clone(),
        status: 1,
        msg: "success".to_string(),
    };

    let result = serde_json::to_string(&msg);
    match result {
        Ok(res) => {
            CString::new(res).unwrap().into_raw()
        }
        Err(e) => {
            CString::new(e.to_string()).unwrap().into_raw()
        }
    }
}

pub fn fail(msg: String) -> *const ::libc::c_char {
    let msg = Message {
        data: "".to_string(),
        status: 0,
        msg,
    };
    let result = serde_json::to_string(&msg);
    match result {
        Ok(res) => {
            CString::new(res).unwrap().into_raw()
        }
        Err(e) => {
            CString::new(e.to_string()).unwrap().into_raw()
        }
    }
}


#[cfg(test)]
mod test {
    #[test]
    fn test_twox_128() {
        let x = "Staking";
        let vec = sp_core::twox_128(x.as_bytes()).to_vec();
        let res = hex::encode(&vec);
        println!("{:?}", res);
    }

    #[test]
    fn test_blake2_128() {
        let dst = "Staking";
        let x = dst.as_bytes();
        let vec = sp_core::blake2_128(x).iter().chain(x.iter()).cloned().collect::<Vec<_>>();
        let res = hex::encode(&vec);
        println!("{:?}", res);
    }

    #[test]
    pub fn twox_64_concat() {
        let dst = "Staking";
        let x = dst.as_bytes();
        let vec = sp_core::twox_64(x)
            .iter()
            .chain(x)
            .cloned()
            .collect::<Vec<_>>();
        let res = hex::encode(&vec);
        println!("{:?}", res);
    }
}