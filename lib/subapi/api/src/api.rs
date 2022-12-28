extern crate libc;

use anyhow::Result;
use serde_json::json;

use ext::UncheckedExtrinsicV4;

use super::utils::get_value;

pub fn inner_dynamic_decode_storage(pallet_name: *const libc::c_char, storage_entry: *const libc::c_char, raw: *const libc::c_char, metadata: *const libc::c_char) -> Result<String> {
    let r_pallet_name = get_value(pallet_name)?;
    let r_storage_entry = get_value(storage_entry)?;
    let r_raw = get_value(raw)?;
    let r_metadata = get_value(metadata)?;
    sub_decode::dynamic_decode_storage(&r_pallet_name, &r_storage_entry, &r_raw, &r_metadata)
}


pub fn inner_signed_extrinsic(hash: *const libc::c_char, seed: *const libc::c_char, to: *const libc::c_char, amount: *const libc::c_char, nonce: *const libc::c_char, spec_version: *const libc::c_char, transaction_version: *const libc::c_char,
                              network_id: *const libc::c_char,module_index: *const libc::c_char,call_index: *const libc::c_char) -> Result<String> {
    let r_hash = get_value(hash)?;
    let r_seed = get_value(seed)?;
    let r_to = get_value(to)?;
    let r_amount = get_value(amount)?;
    let r_nonce = get_value(nonce)?;
    let r_spec_version = get_value(spec_version)?;
    let r_transaction_version = get_value(transaction_version)?;
    let r_network_id = get_value(network_id)?;
    let r_module_index = get_value(module_index)?;
    let r_call_index = get_value(call_index)?;
    ext::signed_extrinsic(r_hash, r_seed, r_to, r_amount, r_nonce, r_spec_version, r_transaction_version, r_network_id,r_module_index,r_call_index)
}


pub fn inner_decode_extrinsic(ext_raw: *const libc::c_char) -> Result<String> {
    let raw = get_value(ext_raw)?;
    ext::decode_extrinsic(raw)
}


pub fn inner_pallet_info(data: *const libc::c_char, pallet_name: *const libc::c_char, call_name: *const libc::c_char) -> Result<String> {
    let r_data = get_value(data)?;
    let r_pallet = get_value(pallet_name)?;
    let r_call_name = get_value(call_name)?;
    let metadata = sub_decode::metadata(&r_data)?;
    let pallet = metadata.pallet(&r_pallet)?;
    let pallet_index = pallet.index();
    let call_index = pallet.call_index(&r_call_name)?;
    let result = json!(
        {
            "pallet_index":pallet_index,
            "call_index":call_index,
        }
    );
    Ok(result.to_string())
}


#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test02() {}
}

