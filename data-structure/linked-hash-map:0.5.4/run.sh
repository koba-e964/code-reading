set -e

URI=https://raw.githubusercontent.com/contain-rs/linked-hash-map/0531e100ef052fd49b2f465abf96cd88aea84692/src/lib.rs
curl ${URI} -o linked_hash_map.rs
patch linked_hash_map.rs diff.patch -o linked_hash_map-mod.rs
sha1sum --check sum.txt 
rustc linked_hash_map-mod.rs
