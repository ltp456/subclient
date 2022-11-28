// Copyright 2019-2022 Parity Technologies (UK) Ltd.
// This file is dual-licensed as Apache-2.0 or GPL-3.0.
// see LICENSE for license details.

use super::storage_address::{
    StorageAddress,
    Yes,
};
use crate::{

    error::Error,
    metadata::{
        DecodeWithMetadata,
        Metadata,
    },

};

use frame_metadata::StorageEntryType;
use scale_info::form::PortableForm;
use sp_core::storage::{
    StorageData,
    StorageKey,
};



/// Validate a storage entry against the metadata.
pub fn validate_storage(
    pallet_name: &str,
    storage_name: &str,
    hash: [u8; 32],
    metadata: &Metadata,
) -> Result<(), Error> {
    let expected_hash = match metadata.storage_hash(pallet_name, storage_name) {
        Ok(hash) => hash,
        Err(e) => return Err(e.into()),
    };
    match expected_hash == hash {
        true => Ok(()),
        false => Err(crate::error::MetadataError::IncompatibleMetadata.into()),
    }
}

/// look up a return type ID for some storage entry.
pub fn lookup_storage_return_type(
    metadata: &Metadata,
    pallet: &str,
    entry: &str,
) -> Result<u32, Error> {
    let storage_entry_type = &metadata.pallet(pallet)?.storage(entry)?.ty;

    Ok(return_type_from_storage_entry_type(storage_entry_type))
}

/// Fetch the return type out of a [`StorageEntryType`].
pub fn return_type_from_storage_entry_type(entry: &StorageEntryType<PortableForm>) -> u32 {
    match entry {
        StorageEntryType::Plain(ty) => ty.id(),
        StorageEntryType::Map { value, .. } => value.id(),
    }
}
