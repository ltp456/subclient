#![cfg_attr(not(feature = "std"), no_std)]

use std::str::FromStr;

use anyhow::{anyhow, bail, Result};
use codec::{Compact, Decode, Encode};
use hex::FromHex;
use mainnet_runtime::{SignedExtra, UncheckedExtrinsic};
use serde::{Deserialize, Serialize};
use sp_core::{Pair, sr25519};
use sp_core::crypto::Ss58AddressFormat;
use sp_runtime::AnySignature;
use sp_runtime::generic::Era;

pub use extrinsic::*;
pub use extrinsic_params::*;

pub mod extrinsic;
pub mod compose;
pub mod extrinsic_params;


#[derive(Serialize, Deserialize)]
pub struct CallExt {
    pub module: String,
    pub call: String,
    pub pallet_index: i16,
    pub call_index: i16,
}

impl Default for CallExt {
    fn default() -> Self {
        CallExt {
            module: "UnKnown".to_string(),
            call: "UnKnown".to_string(),
            pallet_index: -1,
            call_index: -1,
        }
    }
}




// todo

pub fn decode_extrinsic(raw: String) -> Result<String> {
    let mut call_ext = CallExt::default();
    let mut vec = Vec::from_hex(raw)?;
    let ext = UncheckedExtrinsicV4::<([u8;2],), SignedExtra>::decode(&mut vec.as_slice());
    //let ext = UncheckedExtrinsicV4::decode(&mut vec.as_slice());
    match ext {
        Ok(e) => {
            let (call_index, ..): ([u8; 2],) = e.function;
            call_ext.pallet_index = call_index[0] as i16;
            call_ext.call_index = call_index[1] as i16
        }
        Err(e) => {}
    }
    let result = serde_json::to_string(&call_ext)?;
    Ok(result)
}


pub fn signed_extrinsic(hash: String, seed: String, to: String, amount: String, nonce: String, spec_version: String, transaction_version: String, network_id: u16) -> Result<String> {
    // println!("{:?},{:?},{:?},{:?},{:?},{:?},{:?},{:?}",hash,seed,to,amount,nonce,spec_version,transaction_version,network_id);

    let pair = sr25519::Pair::from_string(seed.as_str(), None).map_err(|e| anyhow!("gen pair error {:?}",e))?;
    sp_core::crypto::set_default_ss58_version(Ss58AddressFormat::custom(network_id));
    let to_addr = AccountId::from_str(to.as_str()).map_err(|e| anyhow!("gen to addr error {:?}",e))?;
    let address = GenericAddress::Id(to_addr);

    let new_amount = amount.parse::<u128>().map_err(|e| anyhow!("parse amount to err {:?}",e))?;
    let new_nonce = nonce.parse::<u32>().map_err(|e| anyhow!("parse nonce  err {:?}",e))?;

    let genesis_hash = sp_core::H256::from_str(hash.as_str()).map_err(|e| anyhow!("parse hash error {:?}",e))?;

    let module_index = 4 as u8;
    let call_index = 0 as u8;
    let call = ([module_index, call_index], address, Compact(new_amount));

    let new_spec_version = spec_version.parse::<u32>().map_err(|e| anyhow!("parse spec version error: {:?}",e))?;
    let new_transaction_version = transaction_version.parse::<u32>().map_err(|e| anyhow!("parse transaction version error: {:?}",e))?;
    let tx_params = PlainTipExtrinsicParamsBuilder::default();

    // other
    // let tx_params = PlainTipExtrinsicParamsBuilder::new()
    //     .era(Era::mortal(period, h.number.into()), head)
    //     .tip(0);
    let extrinsic_params = PlainTipExtrinsicParams::new(
        new_spec_version,
        new_transaction_version,
        new_nonce,
        genesis_hash,
        tx_params,
    );
    let xt = compose_extrinsic_offline!(
            pair,
            call,
            extrinsic_params
        );
    Ok(format!("{:?}", xt.hex_encode()))
}


#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_ext() {
        // let hash = "4830cee9d9d002fbbe8fa446e8e6b15465f0f282c5aee6d41e8b657523d8beb3".to_string();
        // let seed = "0xcb1b3421af91a4860a3e82310561853a82668f803c4174dbd6a58535012f7a18".to_string();
        // let to = "esozUamXY9R14rwM3G5cPaTX1haVSo47orFtKGJSn6pPrLvwD".to_string();
        // let amount = "1000000000000".to_string();
        // let nonce = "0".to_string();
        // let spec_version = "277".to_string();
        // let transaction_version = "2".to_string();

        let hash = "4830cee9d9d002fbbe8fa446e8e6b15465f0f282c5aee6d41e8b657523d8beb3".to_string();
        let seed = "0x908fb32a3389776ea34b053c054adcd3b5ed49dfca380729b3d37610fdb60dbf".to_string();
        let to = "esozUamXY9R14rwM3G5cPaTX1haVSo47orFtKGJSn6pPrLvwD".to_string();
        let amount = "1500000000000".to_string();
        let nonce = "0".to_string();
        let spec_version = "277".to_string();
        let transaction_version = "2".to_string();

        let result = signed_extrinsic(
            hash,
            seed,
            to,
            amount,
            nonce,
            spec_version,
            transaction_version,
            88,
        ).unwrap();
        println!("{:?}", result);
    }

    #[test]
    fn test_decode_extrinsic() {
        let raw = "c1018400fadda0af24a7e6d0bba8d3dd0615915b654d68e8735e213fb37577fc681b3d27010c788166b53ce53f0ec566c950d098b47fc7545ab8b9479dd10f48c84222cd722d5e4d2b3463f1995a336cd0684e2bd16c7979641c5c4a54913de9d8826bf98f4703040009130b0030ef7dba02".to_string();
        let res = decode_extrinsic(raw).unwrap();
        println!("{:?}", res);
    }
}