use anyhow::Result;
use codec::Decode;
use frame_metadata::{RuntimeMetadataPrefixed, StorageEntryType};
use sp_core::Bytes;

use metadata_type::Metadata;

mod metadata_type;
mod hash_cache;
mod subxt_metadata;


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




pub type DecodedValue = scale_value::Value<scale_value::scale::TypeId>;


pub trait DecodeWithMetadata {
    fn decode_storage_with_metadata(
        bytes: &mut &[u8],
        pallet_name: &str,
        storage_entry: &str,
        metadata: &Metadata,
    ) -> Result<DecodedValue>;


    fn decode_with_metadata(
        bytes: &mut &[u8],
        type_id: u32,
        metadata: &Metadata,
    ) -> Result<DecodedValue>;
}


impl DecodeWithMetadata for DecodedValue {
    fn decode_storage_with_metadata(
        bytes: &mut &[u8],
        pallet_name: &str,
        storage_entry: &str,
        metadata: &Metadata,
    ) -> Result<DecodedValue> {
        let ty = &metadata.pallet(pallet_name)?.storage(storage_entry)?.ty;
        let id = match ty {
            StorageEntryType::Plain(ty) => ty.id(),
            StorageEntryType::Map { value, .. } => value.id(),
        };
        Self::decode_with_metadata(bytes, id, metadata)
    }

    fn decode_with_metadata(
        bytes: &mut &[u8],
        type_id: u32,
        metadata: &Metadata,
    ) -> Result<DecodedValue> {
        let res = scale_value::scale::decode_as_type(bytes, type_id, metadata.types())?;
        Ok(res)
    }
}


#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {}
}
