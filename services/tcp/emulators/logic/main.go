package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/pkg/protocol/adm"
)

type Emulator struct {
	addr              string
	state             *DeviceState
	cfg               []byte
	store             *StateStore
	syncInterval      time.Duration
	keepaliveInterval time.Duration
	reconnectInterval time.Duration
	conn              net.Conn
	assembler         adm.Assembler
}

func main() {
	var (
		host         = flag.String("host", "127.0.0.1", "NeoSync TCP server host")
		port         = flag.Int("port", 12346, "NeoSync TCP server port")
		imei         = flag.String("imei", "862843047104450", "device IMEI")
		password     = flag.String("password", "0", "device password")
		stateDir     = flag.String("state-dir", "./emulator-state", "directory for state.json and config.bin")
		firmware     = flag.Uint("firmware", 0x38, "firmware version (hex, e.g. 56=0x38)")
		cfgVersion   = flag.Uint("cfg-version", 1, "configuration version byte")
		syncSec      = flag.Int("sync-sec", 300, "sync interval seconds")
		keepaliveSec = flag.Int("keepalive-sec", 180, "keepalive interval seconds")
		reconnectSec = flag.Int("reconnect-sec", 0, "force reconnect every N seconds (0=keep session)")
		importHex    = flag.String("import-config-hex", "", "optional path to hex config to import on first run")
	)
	flag.Parse()

	defaults := DeviceState{
		IMEI:            *imei,
		Password:        *password,
		FirmwareVersion: uint16(*firmware),
		CfgVersion:      uint8(*cfgVersion),
	}

	store := &StateStore{Dir: *stateDir}
	if *importHex != "" {
		if err := importConfigHex(store, *importHex); err != nil {
			log.Fatalf("import config: %v", err)
		}
	}

	state, cfg, err := store.Load(defaults)
	if err != nil {
		log.Fatalf("load state: %v", err)
	}

	em := &Emulator{
		addr:              fmt.Sprintf("%s:%d", *host, *port),
		state:             state,
		cfg:               cfg,
		store:             store,
		syncInterval:      time.Duration(*syncSec) * time.Second,
		keepaliveInterval: time.Duration(*keepaliveSec) * time.Second,
		reconnectInterval: time.Duration(*reconnectSec) * time.Second,
	}

	log.Printf("emulator imei=%s fw=0x%02x cfg_hash=%08x last_mod=%d config_bytes=%d state_dir=%s",
		em.state.IMEI, em.state.FirmwareVersion, em.state.CfgHash, em.state.LastModTime, len(em.cfg), store.Dir)

	for {
		if err := em.runSession(); err != nil {
			log.Printf("session ended: %v", err)
		}
		time.Sleep(2 * time.Second)
	}
}

func importConfigHex(store *StateStore, path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	cfg, err := hex.DecodeString(compactHex(string(raw)))
	if err != nil {
		return err
	}
	packet, err := ensureCfgPacket(cfg)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(store.Dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(store.configPath(), packet, 0o644)
}

func (em *Emulator) runSession() error {
	log.Printf("connecting to %s", em.addr)
	conn, err := net.Dial("tcp", em.addr)
	if err != nil {
		return err
	}
	em.conn = conn
	defer func() {
		_ = conn.Close()
		em.conn = nil
	}()

	if err := em.sendHello(); err != nil {
		return err
	}

	stop := make(chan struct{})
	var stopOnce sync.Once
	stopSession := func() { stopOnce.Do(func() { close(stop) }) }
	defer stopSession()

	go em.periodicLoop(stop)
	go em.readLoop(stopSession)

	if em.reconnectInterval > 0 {
		timer := time.NewTimer(em.reconnectInterval)
		defer timer.Stop()
		select {
		case <-timer.C:
			log.Printf("reconnect timer fired (%s)", em.reconnectInterval)
			return nil
		case <-stop:
			return fmt.Errorf("connection closed")
		}
	}

	<-stop
	return fmt.Errorf("connection closed")
}

func (em *Emulator) sendHello() error {
	packet, err := buildHelloPacket(em.state)
	if err != nil {
		return err
	}
	log.Printf(">> hello (%d bytes) hash=%08x", len(packet), em.state.CfgHash)
	_, err = em.conn.Write(packet)
	return err
}

func (em *Emulator) sendSync(reason string) error {
	packet := buildSyncPacket(em.state)
	log.Printf(">> sync (%s) hash=%08x", reason, em.state.CfgHash)
	_, err := em.conn.Write(packet)
	return err
}

func (em *Emulator) sendKeepalive() error {
	log.Printf(">> keepalive 0xAA")
	_, err := em.conn.Write([]byte{constants.ADM_RC_KEEP_ALIVE_VALUE})
	return err
}

func (em *Emulator) sendConfig(reason string) error {
	if len(em.cfg) == 0 {
		log.Printf("config empty, skip GET_CFG response")
		return nil
	}
	log.Printf(">> config (%d bytes, %s) hash=%08x", len(em.cfg), reason, em.state.CfgHash)
	_, err := em.conn.Write(em.cfg)
	return err
}

func (em *Emulator) periodicLoop(stop <-chan struct{}) {
	syncTicker := time.NewTicker(em.syncInterval)
	keepTicker := time.NewTicker(em.keepaliveInterval)
	defer syncTicker.Stop()
	defer keepTicker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-syncTicker.C:
			if err := em.sendSync("periodic"); err != nil {
				return
			}
		case <-keepTicker.C:
			if err := em.sendKeepalive(); err != nil {
				return
			}
		}
	}
}

func (em *Emulator) readLoop(stop func()) {
	reader := bufio.NewReader(em.conn)
	buf := make([]byte, 4096)

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("read error: %v", err)
			}
			stop()
			return
		}

		for _, packet := range em.assembler.Feed(buf[:n]) {
			em.handlePacket(packet)
		}
	}
}

func (em *Emulator) handlePacket(packet []byte) {
	if len(packet) < 3 {
		return
	}

	typ := packet[2]
	log.Printf("<< packet type=0x%02x len=%d hex=%s", typ, len(packet), hex.EncodeToString(packet))

	switch typ {
	case constants.ADM_RC_PACK_TYPE_COMMAND_V2:
		em.handleCommand(packet)
	case constants.ADM_RC_PACK_TYPE_CFG_V2:
		em.handleSetConfig(packet)
	default:
		log.Printf("unhandled packet type 0x%02x", typ)
	}
}

func (em *Emulator) handleCommand(packet []byte) {
	if len(packet) < 4 {
		return
	}

	cmdType := packet[3]
	payload := packet[4:]

	switch cmdType {
	case constants.ADM_RC_TYPE_STRING:
		cmd := strings.TrimRight(string(bytes.TrimRight(payload, "\x00")), "\x00")
		log.Printf("<< string command: %q", cmd)
		em.handleStringCommand(cmd)
	case constants.ADM_RC_TYPE_GET_SYNC:
		log.Printf("<< GET_SYNC")
		_ = em.sendSync("GET_SYNC")
	case constants.ADM_RC_TYPE_GET_CFG:
		log.Printf("<< GET_CFG")
		_ = em.sendConfig("GET_CFG")
	case constants.ADM_RC_TYPE_SET_CFG:
		log.Printf("<< SET_CFG inside command packet (%d bytes)", len(packet))
		_ = em.saveConfig(packet)
	default:
		log.Printf("<< unknown command type 0x%02x", cmdType)
	}
}

func (em *Emulator) handleSetConfig(packet []byte) {
	log.Printf("<< SET_CFG packet (%d bytes)", len(packet))
	_ = em.saveConfig(packet)
}

func (em *Emulator) saveConfig(packet []byte) error {
	cfg, err := ensureCfgPacket(packet)
	if err != nil {
		log.Printf("save config: %v", err)
		return err
	}

	if err := applyConfigPacket(em.state, cfg); err != nil {
		log.Printf("apply config: %v", err)
		return err
	}

	em.cfg = cfg
	if err := em.store.Save(em.state, cfg); err != nil {
		log.Printf("persist config: %v", err)
		return err
	}

	log.Printf("config saved: hash=%08x last_mod=%d bytes=%d", em.state.CfgHash, em.state.LastModTime, len(cfg))
	return nil
}

func (em *Emulator) handleStringCommand(cmd string) {
	upper := strings.ToUpper(strings.TrimSpace(cmd))
	var answer []byte

	switch upper {
	case "WHO":
		answer = buildStringAnswer("1,adm333,215,adm333,6,1")
	case "COM0":
		answer = buildStringAnswer("COM0")
	case "ADM20INFO":
		answer = buildStringAnswer("ADM20INFO: 0; 0; 0,0,0,0,0")
	case "BLESENSORINFO":
		answer = buildStringAnswer("BLESENSORINFO: [0]:0;")
	case "FUELINFO":
		answer = buildStringAnswer("FUELINFO: [0,BL]: 0;")
	case "STATUS":
		answer = buildStringAnswer("ID=1 Soft=0x38 GPS=0 Time=00:00:00 Val=1")
	default:
		log.Printf("no stub answer for %q", cmd)
		return
	}

	log.Printf(">> string answer for %q (%d bytes)", cmd, len(answer))
	_, _ = em.conn.Write(answer)
}
