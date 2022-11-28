use anyhow::Result;
use codec::{Decode, Encode};
use frame_metadata::RuntimeMetadataPrefixed;
use scale_decode::visitor::Str;
use sp_core::Bytes;

use crate::dynamic::DecodedValue;
use crate::metadata::{DecodeStaticType, DecodeWithMetadata, Metadata};

mod error;
mod dynamic;
mod metadata;
mod subxt_metadata;
mod storage;


pub fn metadata(data: &str) -> Result<Metadata> {
    let metadata_data = hex::decode(data)?;
    let mut bytes: Bytes = Bytes(metadata_data);
    let meta: RuntimeMetadataPrefixed = Decode::decode(&mut &bytes[..])?;
    let metadata: Metadata = meta.try_into()?;
    Ok(metadata)
}


pub fn dynamic_decode_storage(pallet_name: &str, storage_entry: &str, raw: &str, data: &str) -> Result<String> {
    let metadata_data = hex::decode(data)?;
    let mut bytes: Bytes = Bytes(metadata_data);
    let meta: RuntimeMetadataPrefixed = Decode::decode(&mut &bytes[..])?;
    let metadata: Metadata = meta.try_into()?;

    let raw_data = hex::decode(raw)?;
    let val: DecodedValue = DecodedValue::decode_storage_with_metadata(
        &mut &*raw_data,
        pallet_name,
        storage_entry,
        &metadata,
    )?;

    let json = serde_json::to_string(&val)?;
    let result = json.replace("\\", "");
    Ok(result)
}


pub fn dynamic_decode_storage_by_type<T: Decode>(pallet_name: &str, storage_entry: &str, raw: &str, data: &str) -> Result<T> {
    let metadata_data = hex::decode(data)?;
    let mut bytes: Bytes = Bytes(metadata_data);
    let meta: RuntimeMetadataPrefixed = Decode::decode(&mut &bytes[..])?;
    let metadata: Metadata = meta.try_into()?;

    let raw_data = hex::decode(raw)?;

    let val: T = DecodeStaticType::<T>::decode_storage_with_metadata(
        &mut &*raw_data,
        pallet_name,
        storage_entry,
        &metadata,
    )?;
    Ok(val)
}


#[cfg(test)]
mod tests {
    use std::fs;

    use codec::HasCompact;
    pub use sp_core::crypto::AccountId32 as AccountId;

    use super::*;

    #[test]
    fn test_dynamic_decode_storage() {
        let metadata = fs::read("/Users/jianli.rao/workspace/go/subclient/metadata.txt").unwrap();

        let data = &String::from_utf8(metadata).unwrap();

        let raw = "24000000000000006853e20700000000020000000100000000002812948b0000000002000000020000000508d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d3742590a0000000000000000000000000000020000000502d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48003cdb748cf20300000000000000000000000200000005076d6f646c70792f747273727900000000000000000000000000000000000000002c68470800000000000000000000000000000200000013062c6847080000000000000000000000000000020000000507be5ddb1579b72e84524fc29e78609e3caf42e85aa118ebfe0b0ad404b5bdd25f0bda11020000000000000000000000000000020000002000d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d3742590a0000000000000000000000000000000000000000000000000000000000000200000000009874f00700000000000000";
        let raw = "0000000001aa1a2bb482010000";
        let raw = "be5ddb1579b72e84524fc29e78609e3caf42e85aa118ebfe0b0ad404b5bdd25f070010a5d4e8070010a5d4e80000";

        let pallet_name = "Staking";
        let storage_name = "Ledger";

        let result = dynamic_decode_storage(pallet_name, storage_name, raw, data).unwrap();
        println!("{}", result);
    }


    #[test]
    fn test_dynamic_decode_storage_by_type() {
        let metadata = fs::read("/Users/jianli.rao/workspace/go/subclient/metadata.txt").unwrap();

        let data = &String::from_utf8(metadata).unwrap();

        let raw = "24000000000000006853e20700000000020000000100000000002812948b0000000002000000020000000508d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d3742590a0000000000000000000000000000020000000502d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48003cdb748cf20300000000000000000000000200000005076d6f646c70792f747273727900000000000000000000000000000000000000002c68470800000000000000000000000000000200000013062c6847080000000000000000000000000000020000000507be5ddb1579b72e84524fc29e78609e3caf42e85aa118ebfe0b0ad404b5bdd25f0bda11020000000000000000000000000000020000002000d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d3742590a0000000000000000000000000000000000000000000000000000000000000200000000009874f00700000000000000";
        let raw = "0000000001aa1a2bb482010000";
        let raw = "be5ddb1579b72e84524fc29e78609e3caf42e85aa118ebfe0b0ad404b5bdd25f070010a5d4e8070010a5d4e80000";

        let pallet_name = "Staking";
        let storage_name = "Ledger";
        let result = dynamic_decode_storage_by_type::<StakingLedger>(pallet_name, storage_name, raw, data).unwrap();
        println!("{}", result.stash);
        println!("{}", result.active);
        println!("{}", result.total);
    }

    pub type BalanceOf = u128;


    #[derive(PartialEq, Encode, Decode)]
    pub struct StakingLedger {
        /// The stash account whose balance is actually locked and at stake.
        pub stash: AccountId,
        /// The total amount of the stash's balance that we are currently accounting for.
        /// It's just `active` plus all the `unlocking` balances.
        #[codec(compact)]
        pub total: BalanceOf,
        /// The total amount of the stash's balance that will be at stake in any forthcoming
        /// rounds.
        #[codec(compact)]
        pub active: BalanceOf,

    }


    #[test]
    fn test_demo() {
        let res = r#""[{\"phase\":{\"name\":\"ApplyExtrinsic\",\"values\":[0]},\"event\":{\"name\":\"System\",\"values\":[{\"name\":\"ExtrinsicSuccess\",\"values\":{\"dispatch_info\":{\"weight\":132273000,\"class\":{\"name\":\"Mandatory\",\"values\":[]},\"pays_fee\":{\"name\":\"Yes\",\"values\":[]}}}}]},\"topics\":[]},{\"phase\":{\"name\":\"ApplyExtrinsic\",\"values\":[1]},\"event\":{\"name\":\"System\",\"values\":[{\"name\":\"ExtrinsicSuccess\",\"values\":{\"dispatch_info\":{\"weight\":2341737000,\"class\":{\"name\":\"Mandatory\",\"values\":[]},\"pays_fee\":{\"name\":\"Yes\",\"values\":[]}}}}]},\"topics\":[]}]""#;

        let res = res.replace("\\", "");
        println!("{}", res);
    }
}
