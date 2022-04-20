# go-unpad
Unpad padded data. This tool can be used to extract car files directly from unsealed sectors.

1. Get the offset and length of a deal from a sector
`lotus-miner sectors refs`
Example output
```
Block XXXXXXX:
        2770+34359738368 34091302912 bytes
```
The offset is `34359738368` and the length is `34091302912`

2. Use the tool to generate car file from the sector
```
go-unpad --offset 34359738368 --length 34091302912 -i /storage/unsealed/s-t0xxxxx-2770 -o deal.car
```
