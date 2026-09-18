#!/bin/bash

set -e

DEFINE_SERVER_ADDR="localhost"
DEFINE_SERVER_PORT=12346
DEFINE_HELLO_PACKET="2000010208864011061834612f01300000000000000006025e23b7685163c3fa"
DEFINE_CFG_PACKET="090605e68b5e23b7680200060600000c00000000000000000000000001003c6d326d2e6265656c696e652e72750000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002002a6265656c696e65000000000000000000000000000000000000000000000000000000000000000000000003002a6265656c696e65000000000000000000000000000000000000000000000000000000000000000000000004003e3139332e3139332e3136352e3136350000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000050004d05700000600040100000007000200000800040000000009004400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00093000000000000000000b00010071003800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006f0010000000000000000000000000000000000d001041444d353030544553540000000000000e0001000f0001ff100001ff1100010612000101130004a00fa00f1400049919581b150001001600021e001700022c01180001ff190001041a00010a1b0002e8031c0003140a051d00030606061e002a0000000000000000000000000000000000000000000000000000000000000000000000000000000000001f000101200001ff210001022200010023000100240044000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000025003cfc61cfef5e310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002600050000000000270005000000000028000a0000000000000000000029000a000000000000000000002a000500000000002b000500000000005b00300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002f0001043000010131000100320001ff33000affffffffffffffffffff3400010035000100360001003700010038000106390001ff3a00020f003b0001053c00033c3c3c3d0010000102030001020300010203000102033e0006000000000000400002000041000100420001005c0001ff4400400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000045000100460001006e000105470001ff48000100490001004a0001004b0001004c0005ffffffffff4d0014004b0000004b0000004b0000004b0000004b000090001400000000000000000000000000000000000000004f0005020202020250000500ffffffff5100050303030303520005000000000053000400000000540001017f000101550001ff7d00010156000101570001ff58001e0000000000000000000000000000000000000000000000000000000000005900013c5a0001055d000100600004fffdffff610004030000006200102c01000058020000100e0000b0040000630001006400010565000100670001ff690001ff680001017e0001016a0001006b00021e006d0001006c000200007c00020303700001608000010181000103820001018300802f61646d3530302f3f696d65693d3836343031313036313833343631320000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007900017884000200007a0001007b0001017200010673000114740001198b00010a8d0002ffff9300060000000000008f00030000ff95001e38392e3131312e3134312e313730000000000000000000000000000000009600023a30970001055163c3fa"
DEFINE_SYNC_PACKET="0d000406025e23b7685163c3fa"
DEFINE_KEEPALIVE_PACKET="aa"

DEFINE_EMULATOR_LOG="EMULATOR"
DEFINE_SERVER_LOG="TCPSERV"

NUM_EMULATORS=${1:-1}

log() {
    local type=$1
    local msg=$2

    echo "$(date '+%b %d %H:%M:%S') admrc-emulator[$$]: $type: $msg"
}

send_pack() {
    local packet_hex=$1
    local desc=$2

    log "$DEFINE_EMULATOR_LOG" "Send packet | $desc"

    printf "%s" "$packet_hex" | xxd -r -p >&3
}

handle_command() {
    local type_hex=$1
    local cmd_type_hex=$2
    local cmd_data_hex=$3

    local cmd_data=""
    if [[ -n "$cmd_data_hex" ]]; then
        cmd_data=$(echo "$cmd_data_hex" | xxd -r -p 2>/dev/null || true)
    fi

    case "$cmd_type_hex" in
        01) # строковая команда
            log "$DEFINE_EMULATOR_LOG" "Received string command -> '$cmd_data'"
            cmd_data_upper=$(echo "$cmd_data" | tr '[:lower:]' '[:upper:]')

            case "$cmd_data_upper" in
                "BLESENSOR")
                    response="3e0003424c4553454e534f523a202830293a204643363143464546354533312c202831293a202c202832293a202c202833293a202c202834293a202c2000"
                    send_pack "$response" "Response to BLESENSOR"
                    ;;
                "WHO")
                    response="1c0003312c41444d353030452c3231352c61646d3530302c362c3100"
                    send_pack "$response" "Response to WHO"
                    ;;
                "STATUS")
                    response="78000349443d3120536f66743d30783332204750533d3235302054696d653d31383a30383a32382031352e30392e32352056616c3d30204c61743d35372e3939383330204c6f6e3d35362e313935333020563d30563d3020536174436e743d352b3620537461743d3078303034632c20496e5f616c61726d3d3000"
                    send_pack "$response" "Response to STATUS"
                    ;;
                *)
                    log "EMULATOR" "Unknown string command: '$cmd_data'"
                    response="130003556e6b6e6f776e20636f6d6d616e6400"
                    send_pack "$response" "Response to unknown string command"
                    ;;
            esac
            ;;
        02) # запрос синхронизации
            send_pack "$DEFINE_SYNC_PACKET" "Sync response"
            ;;
        03) # запрос конфигурации
            log "$DEFINE_EMULATOR_LOG" "Configuration request received"
            send_pack "$DEFINE_CFG_PACKET" "Config response"
            ;;
        *)
            log "$DEFINE_EMULATOR_LOG" "Unknown CMD_TYPE: $cmd_type_hex"
            ;;
    esac
}

read_packet() {
    size_bytes=$(dd bs=1 count=2 <&3 2>/dev/null | xxd -p -c2 | tr -d '\n')
    [[ -z "$size_bytes" ]] && return 1

    size=$((0x$(echo "$size_bytes" | sed 's/\(..\)\(..\)/\2\1/')))
    data_len=$((size - 2))

    rest=""
    if (( data_len > 0 )); then
        rest=$(dd bs=1 count=$data_len <&3 2>/dev/null | xxd -p -c256 | tr -d '\n')
    fi

    echo "$size_bytes$rest"
}

imei_to_bcd() {
    local imei=$1
    local bcd=""
    for ((i=0; i<${#imei}; i+=2)); do
        pair=${imei:i:2}
        if [ ${#pair} -eq 1 ]; then
            pair="${pair}f"
        fi
        bcd+="$pair"
    done
    echo "$bcd"
}

generate_random_imei() {
    local count=$1
    local imeis=()
    for i in $(seq 1 "$count"); do
        imei=$(printf "%014d0" $((RANDOM*1000000 + RANDOM)))
        imeis+=("$imei")
    done
    echo "${imeis[@]}"
}

run_emulator() {
    local imei=$1
    local bcd_imei=$(imei_to_bcd $imei)
    local hello_packet="2000010208${bcd_imei}01300000000000000006025e23b7685163c3fa"

    exec 3<>/dev/tcp/$DEFINE_SERVER_ADDR/$DEFINE_SERVER_PORT
    log "$DEFINE_EMULATOR_LOG" "Connecting to $DEFINE_SERVER_ADDR:$DEFINE_SERVER_PORT..."

    echo "$hello_packet" | xxd -r -p >&3
    sleep 1

    (
        local start_time=$(date +%s)
        local last_keepalive=$start_time
        local last_sync=$start_time
        while true; do
            now=$(date +%s)
            (( now - last_keepalive >= 180 )) && { echo "$DEFINE_KEEPALIVE_PACKET" | xxd -r -p >&3; last_keepalive=$now; log "$DEFINE_EMULATOR_LOG" "Keepalive sent for $imei"; }
            (( now - last_sync >= 300 )) && { echo "$DEFINE_SYNC_PACKET" | xxd -r -p >&3; last_sync=$now; log "$DEFINE_EMULATOR_LOG" "Sync sent for $imei"; }
            sleep 1
        done
    ) &

    while true; do
        if raw=$(read_packet); then
            log "$DEFINE_EMULATOR_LOG" "Received -> '$raw'"

            type_hex=${raw:4:2}
            cmd_type_hex=${raw:6:2}
            payload_hex=${raw:8}

            handle_command "$type_hex" "$cmd_type_hex" "$payload_hex"
        fi
    done
}

cat << 'EOF'

  ___           __   __   _____           __
 / _ \    /\   |  \ /  | |  ___)          \ \
| |_| |  /  \  |   v   | | |_   _   _ _   _\ \   __  _____ ___   ___
|  _  | / /\ \ | |\_/| | |  _) | | | | | | |> \ /  \/ (   ) _ \ / _ \
| | |/ /__\ \| |   | | | |___| |_| | |_| / ^ ( ()  < | ( (_) ) |_) )
|_| |_/________\_|   |_| |_____) ._,_|\___/_/ \_\__/\_\ \_)___/|  __/
                               | |                             | |
                               |_|                             |_|

EOF

IMEIS=($(generate_random_imei "$NUM_EMULATORS"))
for imei in "${IMEIS[@]}"; do
    run_emulator "$imei" &
    sleep 0.5
done

wait
