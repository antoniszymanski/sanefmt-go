#!/usr/bin/env nu

const root = path self .
let release = http get https://api.github.com/repos/sane-fmt/sane-fmt/releases/latest

def download [name: string] {
  let url = $release.assets | where name == $name | get 0.browser_download_url
  http get $url | save -f $"($root)/($name)"
}

download sane-fmt-x86_64-unknown-linux-musl
download sane-fmt-wasm32-wasi.wasm
wasm-opt -O4 $"($root)/sane-fmt-wasm32-wasi.wasm" -o $"($root)/sane-fmt-wasm32-wasi.wasm"

$'($release.tag_name)
https://github.com/sane-fmt/sane-fmt/releases/latest
' | save -f $"($root)/VERSION"
